#!/usr/bin/env bash

set -e -u -o pipefail # Fail on error

dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd "$dir"

gopath=${GOPATH:-}
logpath=${LOG_PATH:-}
: ${SCRIPT_PATH:?"Need to set SCRIPT_PATH to build script"}


if [ "$gopath" = "" ]; then
  echo "No GOPATH"
  exit 1
fi

client_dir="$gopath/src/github.com/keybase/client"

err_report() {
  "$client_dir/packaging/slack/send.sh" "Error building, see $logpath"
}

trap 'err_report $LINENO' ERR

$SCRIPT_PATH
