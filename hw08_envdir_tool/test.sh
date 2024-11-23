#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-envdir

export HELLO="SHOULD_REPLACE"
export FOO="SHOULD_REPLACE"
export UNSET="SHOULD_REMOVE"
export ADDED="from original env"
export EMPTY="SHOULD_BE_EMPTY"

result=$(./go-envdir "$(pwd)/testdata/env" "/bin/bash" "$(pwd)/testdata/echo.sh" arg1=1 arg2=2)
expected='HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2'

[ "${result}" = "${expected}" ] || (echo -e "invalid output: ${result}" && exit 1)
echo "PASS"

# testing return code
set +e # to trap exit code of the program and go on

./go-envdir "$(pwd)/testdata/empty" "/bin/bash" "$(pwd)/testdata/retcode10.sh" arg1=1 arg2=2
result=$?

set -e # to exit if the tests fail

[ "${result}" = "10" ] || (echo -e "invalid return code: ${result}" && exit 1)

rm -f go-envdir
echo "PASS"
