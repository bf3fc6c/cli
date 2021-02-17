# bf3-uploader

This project copies the latest release from another repository into this one.

## Prerequisites

- Golang

## Configuring

### Personal Access Token

You will need to create a Personal Access Token with the `read:org` and `repo` scopes.

Once created, save the access token as a `BF3_TOKEN` environment variable.

### Other

- Set the name of the organization to copy the release from as a `CLONE_FROM_ORG` environment variable.
- Set the name of the repo to copy the release from as a `CLONE_FROM_REPO` environment variable.

## Installing

Run `make` once to install the program to your path.

## Running

Run `bf3-uploader` to run the program. Once completed, confirm that the release has been copied at https://github.com/bf3fc6c/cli/releases.


