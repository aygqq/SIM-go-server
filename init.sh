#!/bin/sh

echo "init.sh running"

socat PIPE PTY,link=/dev/ttyS10,raw,echo=1 &

/app/main
