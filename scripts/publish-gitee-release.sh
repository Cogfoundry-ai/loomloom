#!/usr/bin/env bash
set -euo pipefail

OWNER="${GITEE_OWNER:-shengsuanyun}"
REPO="${GITEE_REPO_NAME:-loomloom}"
TAG="${TAG:-}"
TARGET_COMMITISH="${TARGET_COMMITISH:-main}"
ASSET_DIR="${ASSET_DIR:-release}"
PRERELEASE="${PRERELEASE:-auto}"
API_BASE="${GITEE_API_BASE:-https://gitee.com/api/v5}"

usage() {
  cat <<'EOF'
Usage: scripts/publish-gitee-release.sh --tag <tag> [--asset-dir <path>]

Create or reuse a Gitee Release and upload files from the release asset directory.

Environment:
  LOOMLOOM_GITEE_TOKEN  Required. Personal access token with repo release write access.
  GITEE_TOKEN           Optional local alias for LOOMLOOM_GITEE_TOKEN.
  GITEE_OWNER           Repository owner/path (default: shengsuanyun)
  GITEE_REPO_NAME       Repository path (default: loomloom)
  TARGET_COMMITISH      Release target branch/commit (default: main)
  PRERELEASE            true, false, or auto (default: auto)
  GITEE_API_BASE        API base URL (default: https://gitee.com/api/v5)

Options:
  --tag <tag>        Release tag, for example v0.1.0-beta.1
  --asset-dir <path> Directory containing release assets (default: release)
  --help             Show this help text
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --tag)
      TAG="${2:-}"
      shift 2
      ;;
    --asset-dir)
      ASSET_DIR="${2:-release}"
      shift 2
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

if [[ -z "$TAG" ]]; then
  echo "--tag is required" >&2
  exit 1
fi
ACCESS_TOKEN="${LOOMLOOM_GITEE_TOKEN:-${GITEE_TOKEN:-}}"
if [[ -z "$ACCESS_TOKEN" ]]; then
  echo "LOOMLOOM_GITEE_TOKEN is required" >&2
  exit 1
fi
if [[ ! -d "$ASSET_DIR" ]]; then
  echo "asset directory not found: $ASSET_DIR" >&2
  exit 1
fi

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

require_cmd curl
require_cmd go

curl_api() {
  curl --retry 3 --retry-delay 2 --retry-connrefused -fsS "$@"
}

upload_asset() {
  local upload_url="$1"
  local asset="$2"
  local response_file status
  response_file="$(mktemp)"
  status="$(
    curl --retry 3 --retry-delay 2 --retry-connrefused -sS \
      -o "$response_file" \
      -w '%{http_code}' \
      -X POST "$upload_url" \
      -F "access_token=$ACCESS_TOKEN" \
      -F "file=@${asset}"
  )"

  if [[ "$status" != 2* ]]; then
    echo "failed to upload asset: $(basename "$asset") status=$status" >&2
    cat "$response_file" >&2
    rm -f "$response_file"
    exit 1
  fi
  rm -f "$response_file"
}

json_top_level_id() {
  local json_file parser_file
  json_file="$(mktemp)"
  parser_file="$(mktemp).go"
  cat > "$json_file"
  cat > "$parser_file" <<'GO'
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	body, err := os.ReadFile(os.Args[1])
	if err != nil {
		os.Exit(1)
	}
	var payload struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || payload.ID == 0 {
		os.Exit(1)
	}
	fmt.Println(payload.ID)
}
GO
  go run "$parser_file" "$json_file"
  rm -f "$json_file" "$parser_file"
}

if [[ "$PRERELEASE" == "auto" ]]; then
  if [[ "$TAG" =~ -(beta|rc|internal)\.[0-9]+$ ]]; then
    PRERELEASE="true"
  else
    PRERELEASE="false"
  fi
fi

release_body="LoomLoom ${TAG} Gitee distribution."

release_url="$API_BASE/repos/$OWNER/$REPO/releases/tags/$TAG"
create_url="$API_BASE/repos/$OWNER/$REPO/releases"

release_json="$(
  curl_api "$release_url?access_token=$ACCESS_TOKEN" || true
)"
release_id="$(printf '%s' "$release_json" | json_top_level_id || true)"

if [[ -z "$release_id" ]]; then
  echo "creating Gitee release: $OWNER/$REPO $TAG"
  release_json="$(
    curl_api -X POST "$create_url?access_token=$ACCESS_TOKEN" \
      --data-urlencode "tag_name=$TAG" \
      --data-urlencode "target_commitish=$TARGET_COMMITISH" \
      --data-urlencode "name=$TAG" \
      --data-urlencode "body=$release_body" \
      --data-urlencode "prerelease=$PRERELEASE"
  )"
  release_id="$(printf '%s' "$release_json" | json_top_level_id || true)"
fi

if [[ -z "$release_id" ]]; then
  echo "failed to resolve Gitee release id" >&2
  printf '%s\n' "$release_json" >&2
  exit 1
fi

echo "Gitee release id: $release_id"
upload_url="$API_BASE/repos/$OWNER/$REPO/releases/$release_id/attach_files"

while IFS= read -r asset; do
  echo "uploading: $(basename "$asset")"
  upload_asset "$upload_url" "$asset"
done < <(find "$ASSET_DIR" -maxdepth 1 -type f | sort)

echo "published Gitee release: https://gitee.com/$OWNER/$REPO/releases/tag/$TAG"
