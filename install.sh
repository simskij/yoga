#!/bin/sh
set -e

REPO="simskij/yoga"
BIN_DIR="${HOME}/.local/bin"
YO_PATH="~/.yofiles"
CLONE_REPO=""

usage() {
    cat <<EOF
Usage: install.sh [options]

Options:
  --repo <url>        Dotfiles repo to clone after install (triggers one-shot mode)
  --path <path>    Override yo storage path (default: ~/.yofiles)
  --bin-dir <path>    Override binary install directory (default: ~/.local/bin)
  -h, --help          Show this help

One-shot example:
  curl -sL https://raw.githubusercontent.com/simskij/yoga/main/install.sh | sh -s -- \\
    --repo git@github.com:you/dotfiles.git
EOF
}

while [ "$#" -gt 0 ]; do
    case "$1" in
        --repo)      CLONE_REPO="$2"; shift 2 ;;
        --path)   YO_PATH="$2";    shift 2 ;;
        --bin-dir)   BIN_DIR="$2";    shift 2 ;;
        -h|--help)   usage; exit 0 ;;
        *) echo "unknown option: $1" >&2; usage; exit 1 ;;
    esac
done

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)          ARCH=amd64 ;;
    aarch64|arm64)   ARCH=arm64 ;;
    armv7l)          ARCH=arm ;;
    *)
        echo "error: unsupported architecture: $ARCH" >&2
        exit 1
        ;;
esac

case "$OS" in
    linux|darwin) ;;
    *)
        echo "error: unsupported OS: $OS" >&2
        exit 1
        ;;
esac

ASSET="yo_${OS}_${ARCH}"

# Resolve latest release tag
echo "· Fetching latest release..."
TAG="$(curl -sI "https://github.com/${REPO}/releases/latest" \
    | grep -i '^location:' \
    | sed 's|.*/tag/||' \
    | tr -d '\r\n')"

if [ -z "$TAG" ]; then
    echo "error: could not determine latest release tag" >&2
    exit 1
fi
echo "· Latest release: $TAG"

# Download binary
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"
mkdir -p "$BIN_DIR"
TMP="$(mktemp)"
echo "· Downloading $ASSET..."
curl -sL "$DOWNLOAD_URL" -o "$TMP"
chmod +x "$TMP"
mv "$TMP" "${BIN_DIR}/yo"
echo "✓ Installed yo to ${BIN_DIR}/yo"

# Ensure bin dir is on PATH for remainder of script
export PATH="${BIN_DIR}:${PATH}"

if [ -z "$CLONE_REPO" ]; then
    # Interactive mode: just print next steps
    cat <<EOF

Run the following to get started:

  yo init
  yo dots init
  yo dots apply
EOF
else
    # One-shot mode: init + clone + apply
    echo "· Running yo init --path ${YO_PATH}..."
    yo init --path "$YO_PATH"
    echo "· Cloning dotfiles from ${CLONE_REPO}..."
    yo dots clone "$CLONE_REPO"
    echo "✓ Done"
fi
