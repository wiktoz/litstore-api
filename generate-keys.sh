#!/bin/sh

openssl ecparam -genkey -name secp384r1 -out /app/keys/ec.key
chmod 600 /app/keys/ec.key
openssl ec -in /app/keys/ec.key -pubout -out /app/keys/ec-pub.key