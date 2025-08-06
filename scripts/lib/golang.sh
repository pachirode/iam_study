readonly IAM_SUPPORTED_SERVER_PLATFORMS=(
  linux/amd64
  linux/arm64
)

iam::golang::server_targets() {
  local targets=(
    iam-apiserver
  )
  echo "${targets[@]}"
}

IFS=" " read -ra IAM_SERVER_TARGETS <<< "$(iam::golang::server_targets)"
readonly IAM_SERVER_TARGETS
readonly IAM_SERVER_BINARIES=("${IAM_SERVER_TARGETS[@]##*/}")

iam::golang::server_image_targets() {
  # NOTE: this contains cmd targets for iam::build::get_docker_wrapped_binaries
  local targets=(
    cmd/iam-apiserver
  )
  echo "${targets[@]}"
}

IFS=" " read -ra IAM_SERVER_IMAGE_TARGETS <<< "$(iam::golang::server_image_targets)"
readonly IAM_SERVER_IMAGE_TARGETS
readonly IAM_SERVER_IMAGE_BINARIES=("${IAM_SERVER_IMAGE_TARGETS[@]##*/}")

# Returns a sorted newline-separated list containing only duplicated items.
iam::golang::dups() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort | uniq -d
}

# Returns a sorted newline-separated list with duplicated items removed.
iam::golang::dedup() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort -u
}

declare -a IAM_SERVER_PLATFORMS
iam::golang::setup_platforms() {
  if [[ -n "${IAM_BUILD_PLATFORMS:-}" ]]; then
    # IAM_BUILD_PLATFORMS needs to be read into an array before the next
    # step, or quoting treats it all as one element.
    local -a platforms
    IFS=" " read -ra platforms <<< "${IAM_BUILD_PLATFORMS}"

    # Deduplicate to ensure the intersection trick with iam::golang::dups
    # is not defeated by duplicates in user input.
    iam::util::read-array platforms < <(iam::golang::dedup "${platforms[@]}")

    # Use iam::golang::dups to restrict the builds to the platforms in
    # IAM_SUPPORTED_*_PLATFORMS. Items should only appear at most once in each
    # set, so if they appear twice after the merge they are in the intersection.
    iam::util::read-array IAM_SERVER_PLATFORMS < <(iam::golang::dups \
        "${platforms[@]}" \
        "${IAM_SUPPORTED_SERVER_PLATFORMS[@]}" \
      )
    readonly IAM_SERVER_PLATFORMS

    iam::util::read-array IAM_CLIENT_PLATFORMS < <(iam::golang::dups \
        "${platforms[@]}" \
        "${IAM_SUPPORTED_CLIENT_PLATFORMS[@]}" \
      )
    readonly IAM_CLIENT_PLATFORMS

  elif [[ "${IAM_FASTBUILD:-}" == "true" ]]; then
    IAM_SERVER_PLATFORMS=(linux/amd64)
    readonly IAM_SERVER_PLATFORMS
    IAM_CLIENT_PLATFORMS=(linux/amd64)
    readonly IAM_CLIENT_PLATFORMS
  else
    IAM_SERVER_PLATFORMS=("${IAM_SUPPORTED_SERVER_PLATFORMS[@]}")
    readonly IAM_SERVER_PLATFORMS

    IAM_CLIENT_PLATFORMS=("${IAM_SUPPORTED_CLIENT_PLATFORMS[@]}")
    readonly IAM_CLIENT_PLATFORMS
  fi
}
iam::golang::setup_platforms

iam::golang::setup_env() {
  iam::golang::verify_go_version

  # Unset GOBIN in case it already exists in the current session.
  unset GOBIN

  # This seems to matter to some tools
  export GO15VENDOREXPERIMENT=1

  # Open go module feature
  export GO111MODULE=on

  # This is for sanity.  Without it, user umasks leak through into release
  # artifacts.
  umask 0022
}

iam::golang::verify_go_version() {
  if [[ -z "$(command -v go)" ]]; then
    iam::log::usage_from_stdin <<EOF
Can't find 'go' in PATH, please fix and retry.
See http://golang.org/doc/install for installation instructions.
EOF
    return 2
  fi

  local go_version
  IFS=" " read -ra go_version <<< "$(go version)"
  local minimum_go_version
  minimum_go_version=go1.24.5
  if [[ "${minimum_go_version}" != $(echo -e "${minimum_go_version}\n${go_version[2]}" | sort -s -t. -k 1,1 -k 2,2n -k 3,3n | head -n1) && "${go_version[2]}" != "devel" ]]; then
    iam::log::usage_from_stdin <<EOF
Detected go version: ${go_version[*]}.
IAM requires ${minimum_go_version} or greater.
Please install ${minimum_go_version} or later.
EOF
    return 2
  fi
}

iam::golang::host_platform() {
  echo "$(go env GOHOSTOS)/$(go env GOHOSTARCH)"
}

