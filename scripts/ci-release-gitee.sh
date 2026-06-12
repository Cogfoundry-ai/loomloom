#!/usr/bin/env bash
set +x
set -euo pipefail

RELEASE_TAG_PATTERN='^(v.*|gitee-release-test\..*)$'

log_context() {
  echo "CI context:"
  echo "GITEE_TAG=${GITEE_TAG:-}"
  echo "GITEE_REF_NAME=${GITEE_REF_NAME:-}"
  echo "GITEE_COMMIT=${GITEE_COMMIT:-}"
  echo "PWD=$(pwd)"
  git status --short || true
  git tag --points-at "${GITEE_COMMIT:-HEAD}" || true
}

require_env() {
  if [[ -z "${LOOMLOOM_GITEE_TOKEN:-}" ]]; then
    echo "LOOMLOOM_GITEE_TOKEN is required" >&2
    echo "Check the Gitee pipeline secret variable config, not only the YAML declaration." >&2
    exit 1
  fi
  echo "LOOMLOOM_GITEE_TOKEN is present"
}

setup_go() {
  export GOPROXY="${GOPROXY_OVERRIDE:-https://goproxy.cn,direct}"
  export GOSUMDB="${GOSUMDB_OVERRIDE:-sum.golang.org}"

  local go_version
  go_version="$(awk '/^go / { print $2; exit }' cli/go.mod)"
  if [[ -z "$go_version" ]]; then
    echo "failed to read Go version from cli/go.mod" >&2
    exit 1
  fi

  local current_version current_goroot
  current_version="$(go env GOVERSION 2>/dev/null || true)"
  current_goroot="$(go env GOROOT 2>/dev/null || true)"
  if [[ "$current_version" == "go${go_version}" && "$current_goroot" != *"/standard/golang/1.13" ]]; then
    go version
    go env GOVERSION GOROOT GOPROXY GOSUMDB
    return
  fi

  local go_platform
  go_platform="$(go_platform)"
  local go_archive="go${go_version}.${go_platform}.tar.gz"
  download_go "$go_version" "$go_archive"

  rm -rf /tmp/go
  tar -C /tmp -xzf "/tmp/${go_archive}"
  export GOROOT="/tmp/go"
  export PATH="/tmp/go/bin:$PATH"

  go version
  go env GOVERSION GOROOT GOPROXY GOSUMDB
}

download_go() {
  local go_version="$1"
  local go_archive="$2"
  local download_bases="${GO_DOWNLOAD_BASE_URLS:-https://mirrors.aliyun.com/golang https://golang.google.cn/dl https://go.dev/dl}"
  local base url attempt

  echo "installing Go ${go_version}"
  for base in $download_bases; do
    url="${base%/}/${go_archive}"
    echo "trying Go download: $url"
    for attempt in 1 2 3; do
      if curl -fL --connect-timeout 15 --retry 2 --retry-delay 2 -o "/tmp/${go_archive}" "$url"; then
        return
      fi
      echo "download Go failed: url=$url attempt=$attempt" >&2
      sleep 2
    done
  done

  echo "failed to download Go ${go_version} from all configured mirrors" >&2
  echo "GO_DOWNLOAD_BASE_URLS=${download_bases}" >&2
  exit 1
}

go_platform() {
  local os arch
  case "$(uname -s)" in
    Linux) os="linux" ;;
    Darwin) os="darwin" ;;
    *)
      echo "unsupported OS for Go bootstrap: $(uname -s)" >&2
      exit 1
      ;;
  esac

  case "$(uname -m)" in
    x86_64|amd64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *)
      echo "unsupported architecture for Go bootstrap: $(uname -m)" >&2
      exit 1
      ;;
  esac

  printf '%s-%s\n' "$os" "$arch"
}

resolve_tag() {
  git fetch --tags --force >&2

  local tag_name="${GITEE_TAG:-${GITEE_REF_NAME:-}}"
  if [[ -z "$tag_name" ]]; then
    tag_name="$(select_release_tag "${GITEE_COMMIT:-HEAD}")"
  fi

  if [[ -z "$tag_name" ]]; then
    echo "failed to resolve release tag" >&2
    echo "tags on HEAD:" >&2
    git tag --points-at "${GITEE_COMMIT:-HEAD}" >&2 || true
    exit 1
  fi

  if ! grep -Eq "$RELEASE_TAG_PATTERN" <<<"$tag_name"; then
    echo "resolved tag '$tag_name' is not a release tag" >&2
    exit 1
  fi

  printf '%s\n' "$tag_name" > .release-tag
  printf '%s\n' "$tag_name"
}

select_release_tag() {
  local ref="$1"
  local tags
  tags="$(git tag --points-at "$ref" | grep -E "$RELEASE_TAG_PATTERN" || true)"

  printf '%s\n' "$tags" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -n 1 && return
  printf '%s\n' "$tags" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-' | head -n 1 && return
  printf '%s\n' "$tags" | grep -E '^gitee-release-test\.' | head -n 1 || true
}

build_assets() {
  local tag_name="$1"
  scripts/package-release-assets.sh --version "$tag_name"

  echo "release assets:"
  find release -maxdepth 1 -type f -print | sort
}

publish_release() {
  local tag_name="$1"
  if [[ "${CI_RELEASE_GITEE_SKIP_PUBLISH:-}" == "1" ]]; then
    echo "CI_RELEASE_GITEE_SKIP_PUBLISH=1; skipping Gitee release publish"
    return
  fi
  scripts/publish-gitee-release.sh --tag "$tag_name"
}

main() {
  log_context
  require_env
  setup_go

  local tag_name
  tag_name="$(resolve_tag)"
  echo "release tag: $tag_name"

  build_assets "$tag_name"
  publish_release "$tag_name"
}

main "$@"
