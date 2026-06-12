#!/usr/bin/env bash
set -euo pipefail

GITEE_INSTALL_URL="${GITEE_INSTALL_URL:-https://gitee.com/shengsuanyun/loomloom/raw/main/install.sh}"
export LOOMLOOM_RELEASE_SOURCE="${LOOMLOOM_RELEASE_SOURCE:-gitee}"

SOURCE_PATH="${BASH_SOURCE[0]:-}"
if [[ -n "$SOURCE_PATH" && "$SOURCE_PATH" != "-" && -f "$SOURCE_PATH" ]]; then
  SCRIPT_DIR="$(cd "$(dirname "$SOURCE_PATH")" && pwd)"
  if [[ -f "$SCRIPT_DIR/install.sh" ]]; then
    exec "$SCRIPT_DIR/install.sh" --source gitee "$@"
  fi
fi

TMP_INSTALLER="$(mktemp)"
trap 'rm -f "$TMP_INSTALLER"' EXIT

curl -fsSL -o "$TMP_INSTALLER" "$GITEE_INSTALL_URL"
exec bash "$TMP_INSTALLER" --source gitee "$@"
