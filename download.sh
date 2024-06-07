#!/bin/sh
set -e

RESET="\\033[0m"
RED="\\033[31;1m"
GREEN="\\033[32;1m"
YELLOW="\\033[33;1m"
BLUE="\\033[34;1m"
WHITE="\\033[37;1m"

say_green()
{
    [ -z "${SILENT}" ] && printf "%b%s%b\\n" "${GREEN}" "$1" "${RESET}"
    return 0
}

say_red()
{
    printf "%b%s%b\\n" "${RED}" "$1" "${RESET}"
}

say_yellow()
{
    [ -z "${SILENT}" ] && printf "%b%s%b\\n" "${YELLOW}" "$1" "${RESET}"
    return 0
}

say_blue()
{
    [ -z "${SILENT}" ] && printf "%b%s%b\\n" "${BLUE}" "$1" "${RESET}"
    return 0
}

say_white()
{
    [ -z "${SILENT}" ] && printf "%b%s%b\\n" "${WHITE}" "$1" "${RESET}"
    return 0
}

print_unsupported_platform()
{
    >&2 say_red "error: We're sorry, but it looks like Omnistrate CLI is not supported on your platform"
    >&2 say_red "       We support 64-bit versions of Linux, macOS, and Windows."
}

at_exit()
{
    if [ "$?" -ne 0 ]; then
        >&2 say_red "We're sorry, but it looks like something might have gone wrong during installation."
        >&2 say_red "If you need help, please join us on our support channels."
    fi
}

trap at_exit EXIT

OS=""
ARCH=""
case $(uname) in
    "Linux") OS="linux";;
    "Darwin") OS="darwin";;
    "MINGW"*) OS="windows";;
    "MSYS"*) OS="windows";;
    *)
        print_unsupported_platform
        exit 1
        ;;
esac

case $(uname -m) in
    "x86_64") ARCH="amd64";;
    "arm64") ARCH="arm64";;
    "aarch64") ARCH="arm64";;
    *)
        print_unsupported_platform
        exit 1
        ;;
esac

VERSION="0.8"
BASE_URL="https://github.com/omnistrate/cli/releases/download/v${VERSION}/omnistrate-ctl-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BASE_URL="${BASE_URL}.exe"
fi

INSTALL_ROOT=""
if [ -z "${INSTALL_ROOT}" ]; then
    INSTALL_ROOT="${HOME}/.omnistrate"
fi
CLI_PATH="${INSTALL_ROOT}/bin/omnistrate-ctl"

if [ -d "${CLI_PATH}" ]; then
    say_red "error: ${CLI_PATH} already exists and is a directory, refusing to proceed."
    exit 1
elif [ ! -f "${CLI_PATH}" ]; then
    say_blue "=== Installing Omnistrate CLI v${VERSION} ==="
else
    say_blue "=== Upgrading Omnistrate CLI to v${VERSION} ==="
fi

mkdir -p "${INSTALL_ROOT}/bin"

say_blue "=== Downloading Omnistrate CLI v${VERSION} for ${OS}-${ARCH} ==="
curl -L -o "${CLI_PATH}" ${BASE_URL}

if [ "$OS" = "windows" ]; then
    mv "${CLI_PATH}" "${INSTALL_ROOT}/bin/omnistrate-ctl.exe"
    CLI_PATH="${INSTALL_ROOT}/bin/omnistrate-ctl.exe"
fi

chmod +x "${CLI_PATH}"
say_green "Omnistrate CLI downloaded to ${CLI_PATH}"

# Add to PATH if not already added
PROFILE_FILE=""
SHELL_NAME=$(basename "${SHELL}")
case "${SHELL_NAME}" in
    "bash")
        if [ "$(uname)" != "Darwin" ]; then
            if [ -e "${HOME}/.bashrc" ]; then
                PROFILE_FILE="${HOME}/.bashrc"
            elif [ -e "${HOME}/.bash_profile" ]; then
                PROFILE_FILE="${HOME}/.bash_profile"
            fi
        else
            if [ -e "${HOME}/.bash_profile" ]; then
                PROFILE_FILE="${HOME}/.bash_profile"
            elif [ -e "${HOME}/.bashrc" ]; then
                PROFILE_FILE="${HOME}/.bashrc"
            fi
        fi
        ;;
    "zsh")
        if [ -e "${ZDOTDIR:-$HOME}/.zshrc" ]; then
            PROFILE_FILE="${ZDOTDIR:-$HOME}/.zshrc"
        fi
        ;;
esac

if [ -n "${PROFILE_FILE}" ]; then
    LINE_TO_ADD="export PATH=\$PATH:${INSTALL_ROOT}/bin"
    if ! grep -q "# add Omnistrate to the PATH" "${PROFILE_FILE}"; then
        say_white "+ Adding ${INSTALL_ROOT}/bin to \$PATH in ${PROFILE_FILE}"
        printf "\\n# add Omnistrate to the PATH\\n%s\\n" "${LINE_TO_ADD}" >> "${PROFILE_FILE}"
    fi
    EXTRA_INSTALL_STEP="+ Please restart your shell or add ${INSTALL_ROOT}/bin to your \$PATH"
else
    EXTRA_INSTALL_STEP="+ Please add ${INSTALL_ROOT}/bin to your \$PATH"
fi

say_blue "=== Omnistrate CLI is now installed! ==="
if [ -n "${EXTRA_INSTALL_STEP}" ]; then
    say_white "${EXTRA_INSTALL_STEP}"
fi
say_green "+ Get started with Omnistrate: https://docs.omnistrate.com/getting-started/ctl-reference/#getting-started"
