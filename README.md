# Judge0 Uploader

###### _CLI that helps you submit your local code into Judge0 easily_

```bash
###############################
# REAL
# Create challenge-how-edabit-works
npx @roeeyn/challenge-generator --edabit-id ARr5tA458o2tC9FTN -s

# Submit challenge
go run main.go submit sample_challenge --verbose

# fetch status
go run main.go status 4ec633f1-4dca-4870-803c-7a01bf4aef23 --verbose

# Submit and wait for the response
go run main.go submit challenge-how-edabit-works | xargs go run main.go status --wait --verbose
###############################
```

Whenever we're working on new coding challenges, we want to try them as fast as possible with [Judge0](https://github.com/judge0/judge0), the platform we use to test our participants code.
As there are too much submissions, we wanted a CLI that helped us to try the local code with the least friction possible.

This CLI is expected to work with the file structure of the challenges generated with the [challenge-generator](https://github.com/roeeyn/challenge-generator) project. If you would like to support any other structure, please create a new issue.

## Configuration

You need to setup the token via config file or flag
You can set your token in different ways:

### ENV variable

```bash
export JUDGE0_AUTH_TOKEN='YOUR_TOKEN'
```

### YAML Config File

The default config file location is `~` or `$HOME`, and the expected name is `.judge0-uploader.yaml`, so the expected file complete path should be `~/.judge0-uploader.yaml`.

You can create it with the following command.

```bash
echo 'judge0_auth_token: "YOUR_TOKEN"' > ~/.judge0-uploader.yaml
```

You can also specify a different path for your file configuration by writing:

```bash
... --config YOUR_CONFIG_FILE_PATH
```

## Usage

### 0. Alias

Create a new alias to avoid writing the complete command.

> **Warning** If you do not create the alias you will have to use `judge0-uploader` instead of `j0` in each command.

```bash
# OPTIONAL - to avoid writing judge0-uploader every time.
alias j0=judge0-uploader

```

### 1. Upload New Submission

To upload and execute a new submission use the following command. Remember that the first version of this CLI is designed for the file structure of the [challenge-generator](https://github.com/roeeyn/challenge-generator) project. This means we're expecting to find in the challenge path:

- `run`
- `index.js`\*
- `test.js`\*
- `testframework.js`\*

\* This files can be of any extension (e.g. `js`, `java`, `py`) but the extension must match between all the files.

```
# To upload a new submission:
j0 upload path-to-your-challenge/

123456 # Your submission ID
```

### 2. Look For the Result

As the multipart submission cannot be waited directly in the API we have to be polling the Judge0 API for the execution result. To know the result of your current challenge execute and wait until the execution finishes use:

> **Note** A submission is considered in execution if the status is `In Queue` or `Processing`.

```bash
# To get your current submission status
j0 status YOUR-SUBMISSION-ID
```

If for some reason you don't want to wait until the execution finishes, pass the flag `--no-wait`.

```bash
# To get you current submission status as it is at that point.
j0 status YOUR-SUBMISSION-ID --no-wait
```
