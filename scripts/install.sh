#!/usr/bin/env bash

{ # this ensures the entire script is downloaded #

BINARY_DEST="$HOME/bin"
BINARY_NAME="rhoas"
SRC_ORG="bf3fc6c"
SRC_REPO="cli"
OS_TYPE="linux"
OS_LONG_BIT="64"
API_BASE_URL="https://api.github.com"
API_RELEASE_URL="${API_BASE_URL}/repos/${SRC_ORG}/${SRC_REPO}/releases/latest"
DOWNLOAD_DIR="/tmp"

has_in_path() {
  type "$1" > /dev/null 2>&1
}

source_file() {
  # shellcheck disable=SC1090
  source "$1" 2> /dev/null
}

# get OS type
if [[ "$OSTYPE" == "darwin"* ]]; then
  OS_TYPE="macOS"
fi

if ! has_in_path "curl"; then
  echo "curl is required to download this binary"
  exit 1
fi

if [ ! -d "$BINARY_DEST" ]; then
  mkdir "$BINARY_DEST"
fi

LATEST_TAG=$(curl -s "${API_RELEASE_URL}" \
| grep "tag_name.*" \
| cut -d '"' -f 4)

ASSET_NAME="${BINARY_NAME}_${LATEST_TAG}_${OS_TYPE}_amd${OS_LONG_BIT}"
ASSET_NAME_COMPRESSED="${ASSET_NAME}.tar.gz"

DOWNLOAD_URL=$(curl -s "${API_RELEASE_URL}" \
| grep "browser.download_url.*${ASSET_NAME_COMPRESSED}" \
| cut -d '"' -f 4)

cd "$DOWNLOAD_DIR" || exit

# wget is faster, use it to download the release if available
if has_in_path "wget"; then
  echo "Downloading $BINARY_NAME v${LATEST_TAG}"
  wget -q "$DOWNLOAD_URL"
else
  echo "Downloading $BINARY_NAME v${LATEST_TAG}"
  curl -sL "$DOWNLOAD_URL" --output "${ASSET_NAME_COMPRESSED}"
fi

  echo "$BINARY_NAME v${LATEST_TAG} downloaded"


# unpack and place the binary in the users PATH
tar xvf "$ASSET_NAME_COMPRESSED"
rm -rf "${ASSET_NAME_COMPRESSED}"
cp "${ASSET_NAME}/bin/${BINARY_NAME}" "${BINARY_DEST}/${BINARY_NAME}"

echo "$BINARY_DEST"

echo "$BINARY_NAME v$LATEST_TAG"
echo
echo "Please ensure that $BINARY_DEST is on your PATH"
} # this ensures the entire script is downloaded #