# screeps-pusher

It is a small utility program and GitHub Action to push JavaScript code into Screeps.com.

## Inputs
| Input   | Description |
|---------|-------------|
| branch  | Required. Destination branch in screeps. Be aware that screeps does not create new branch automatically. |
| token   | Required. Screeps API token. Can be obtained from account settings: https://screeps.com/a/#!/account/auth-tokens. |
| dir     | Optional. Directory for .js files. |
| apiURL  | Optional. Screeps server API URL. |
| timeout | Optional. HTTP Client timeout. |

## Example

```yaml
steps:
- name: Checkout
  uses: actions/checkout@master

- name: Push the code
  id: myaction
  uses: kskitek/screeps-pusher@v0.1
  with:
    branch: experimental
    token: ${{ secrets.token }}
```

## Full example

```yaml
steps:
- name: Checkout
  uses: actions/checkout@master

- name: Push the code
  id: myaction
  uses: kskitek/screeps-pusher@v0.1
  with:
    branch: experimental
    token: ${{ secrets.token }}
    dir: src/
    apiURL: https://screeps.com/api/user/code
    timeout: 10s
```

## Full workflow example to push on tag/release

`.github/workflows/tag.yml`

```yaml
name: Push released code
on: tags
jobs:
  push:
    steps:
    - name: Checkout
      uses: actions/checkout@master

    - name: Push the code
      id: myaction
      uses: kskitek/screeps-pusher@v0.1
      with:
        branch: main
        token: ${{ secrets.token }}
```

## Full workflow example to push on PR

`.github/workflows/pr.yml`

```yaml
name: Push experimental code
on: pull_request
jobs:
  push:
    steps:
    - name: Checkout
      uses: actions/checkout@master

    - name: Push the code
      id: myaction
      uses: kskitek/screeps-pusher@v0.1
      with:
        branch: ${{ github.ref_name }} # be aware that screeps does not create new branch automatically
        token: ${{ secrets.token }}
