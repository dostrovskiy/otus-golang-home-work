#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet

# 1. Happy path
(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -l localhost 4242 >/tmp/nc.out &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost 4242 >/tmp/telnet.out &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I
am
TELNET client'
fileEquals /tmp/nc.out "${expected_nc_out}"

expected_telnet_out='Hello
From
NC'
fileEquals /tmp/telnet.out "${expected_telnet_out}"

# 2. SIGINT
(echo -e "Hello\n" && cat) | nc -l localhost 4242 &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat) | ./go-telnet --timeout=5s localhost 4242 &
TL_PID=$!

sleep 1
kill -s SIGINT ${TL_PID}

sleep 5
! ps -p ${TL_PID} > /dev/null 2>&1 || (echo "2. Process telnet ${TL_PID} is still running after SIGINT" && kill ${TL_PID} && kill ${NC_PID} && exit 1)
kill ${NC_PID} 2>/dev/null || true

# 3.Server EOF
(echo -e "Hello\n" && sleep 1 && echo -e "Bye\n" && sleep 1) | nc -l localhost 4242 &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && sleep 2 && echo -e "Another message\n") | ./go-telnet --timeout=5s localhost 4242 >./telnet.out &
TL_PID=$!

sleep 5
! ps -p ${TL_PID} > /dev/null 2>&1 || (echo "3. Process telnet ${TL_PID} is still running after server EOF" && kill ${TL_PID} && kill ${NC_PID} && exit 1)
kill ${NC_PID} 2>/dev/null || true

# end of tests
rm -f /tmp/nc.out
rm -f /tmp/telnet.out
rm -f go-telnet
echo "PASS"
