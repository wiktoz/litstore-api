#!/bin/sh
mkdir -p /app/keys
openssl ecparam -name secp384r1 -genkey -noout -out /app/keys/ec.pem
chmod 600 /app/keys/ec.pem
openssl ec -in /app/keys/ec.pem -pubout -out /app/keys/ec-pub.pem