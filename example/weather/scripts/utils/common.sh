# Useful functions for re-use in different scripts

[[ "$DEBUG" != "" ]] && set -x

function is_mac {
  [[ "$OSTYPE" == "darwin"* ]]
}

function is_m1_mac {
  is_mac && [[ "$(uname -a)" == *"ARM64"* ]]
}

function check_required_cmd {
  cmd="$1"
  pkg="$2"
  # I feel like you could do this with substitution
  # but I didn't feel like fighting bash.
  [[ "$pkg" == "" ]] && pkg=$cmd

  if ! command -v $1 &> /dev/null; then
    echo "Unable to find '$cmd' in your PATH - cannot continue."
    if is_mac; then
      echo "Try 'brew install '$pkg'"
    else
      echo "Try 'apt-get install '$pkg'"
    fi
    exit 1
  fi
}
