#!/bin/sh

echo IP address is $1
echo Operator ID is $2

ssh admin@$1 'interface lte set operator='$2' lte1'

echo Done
