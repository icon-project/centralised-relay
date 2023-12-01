#!/bin/bash
set -e

# Source external scripts
. scripts/optimize-xcall-build.sh

# Define a usage function to display how to use the script
usage() {
  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  --clean : Clean contract directories (true/false, default: false)."
  echo "  --build-xcall : Build xCall contracts (true/false, default: false)."
  echo "  --xcall-branch <branch>: Specify the xCall branch to build (default: main)."
  echo "  --test <test_type>: Specify the type of test (e2e, e2e-demo, integration, default: e2e)."
  exit 1
}

# Define variables with default values
clean="false"
build_xcall="false"
test="e2e"
xcall_branch="main"

# Define functions

clean_contracts() {
  echo "Cleaning contract directories..."
  find artifacts/icon -type f -exec rm {} \;
  #  find artifacts/archway -type f -exec rm {} \; not required right now
  find artifacts/evm -type f -exec rm {} \;
}

e2e_test() {
  echo "Running e2e test..."
  go test -v ./test/e2e -timeout 0 -count=1
}

e2e_demo() {
  echo "Configuring e2e demo..."
  export PRESERVE_DOCKER=true && go test -v ./test/e2e-demo -testify.m TestSetup
}

integration_test() {
  echo "Running integration test..."
  go test -v ./test/integration -timeout 0 -count=1
}

# Parse command line arguments
while [ $# -gt 0 ]; do
  case "$1" in
  --clean)
    clean="true"
    ;;
  --build-xcall)
    build_xcall="true"
    ;;
  --xcall-branch)
    shift
    xcall_branch="$1"
    ;;
  --test)
    shift
    test="$1"
    ;;
  *)
    echo "Error: Unknown option '$1'."
    usage
    ;;
  esac
  shift
done
# Perform actions based on command line arguments

if [ "$clean" = "true" ]; then
  clean_contracts
fi

if [ "$build_xcall" = "true" ]; then
  echo "building xCall contracts..."
  build_xCall_contracts "$xcall_branch"
fi

# Run the selected test
echo "running $test......"
case "$test" in
"e2e")
  e2e_test
  ;;
"e2e-demo")
  e2e_demo
  ;;
"integration")
  integration_test
  ;;
*)
  echo "Error: Unknown test type '$test'."
  exit 1
  ;;
esac
