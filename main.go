package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var VERSION string

var jsRegex = regexp.MustCompile(`\.js$`)
var newLineRegex = regexp.MustCompile(`\n`)

type config struct {
	Branch  string        `required:"true"`
	Dir     string        `default:"./"`
	Token   string        `required:"true"`
	ApiURL  string        `default:"https://screeps.com/api/user/code"`
	Timeout time.Duration `default:"10s"`
	// TODO idea: add an option to automatically switch to the branch - POST /api/user/set-active-branch
}

func main() {
	var c config
	err := envconfig.Process("INPUT", &c)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	payload, err := preparePayload(c)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = send(payload, c)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func preparePayload(c config) (io.Reader, error) {
	fr := filesReader{modules: make(map[string]string)}

	err := filepath.Walk(c.Dir, fr.read)
	if err != nil {
		return nil, fmt.Errorf("unable to walk files: %w", err)
	}

	return fr.ToJson(c.Branch), nil
}

type filesReader struct {
	modules map[string]string
}

func (r *filesReader) read(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if !jsRegex.MatchString(path) {
		return nil
	}
	log.Println("Adding:", path)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file [%s]: %w", path, err)
	}

	fileContent := string(b)
	fileContent = newLineRegex.ReplaceAllString(fileContent, `\n`)
	fileContent = strings.ReplaceAll(fileContent, `"`, `\"`)

	module := jsRegex.ReplaceAllString(path, "")
	r.modules[module] = fileContent

	return nil
}

func (r *filesReader) ToJson(branch string) io.Reader {
	sb := strings.Builder{}
	sb.WriteString("{")
	sb.WriteString(`"branch":"`)
	sb.WriteString(branch)
	sb.WriteString(`","modules":{`)

	for module, fileContent := range r.modules {
		sb.WriteString(`"`)
		sb.WriteString(module)
		sb.WriteString(`":"`)
		sb.WriteString(fileContent)
		sb.WriteString(`", `)
	}
	sb.WriteString(`"empty":""`) // just a lazy way of not dealing with dangling comma
	sb.WriteString("}}")

	return bytes.NewBufferString(sb.String())
}

func send(payload io.Reader, c config) error {
	client := &http.Client{
		Timeout: c.Timeout,
	}

	req, err := prepareRequest(payload, c.ApiURL, c.Token)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send the code: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unable to send the code: Status: %d, Response: %s", resp.StatusCode, respBody)
	}

	var errBody errorBody
	_ = json.NewDecoder(resp.Body).Decode(&errBody)
	if errBody.Error != "" {
		return fmt.Errorf("API returned an error: Status: %d, Response: %s", resp.StatusCode, errBody.Error)
	}

	log.Println("Code pushed successfully")
	return nil
}

type errorBody struct {
	Error string `json:"error"`
}

func prepareRequest(payload io.Reader, apiURL, token string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, apiURL, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("X-Token", token)
	req.Header.Add("User-Agent", "screeps-pusher/"+VERSION)

	return req, nil
}
