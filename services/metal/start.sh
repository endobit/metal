#!/bin/sh

cat banner.txt

while true; do
	/bin/stackd --log-level=debug --cert-path=/certs --dbpath=/data
	echo
	date
	echo service crashed, restarting...
	echo
	sleep 10
done
