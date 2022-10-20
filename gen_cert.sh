#!/bin/env bash

IP=$1

"${IP:?Must specify an IP address}" 2>/dev/null

openssl genrsa -out self-signed-ca.key 2048
openssl req -new -x509 -days 365 -key self-signed-ca.key -subj "/C=US/ST=NC/O=David Greeson/CN=greeson.xyz" -out self-signed-ca.crt

openssl req -newkey rsa:2048 -nodes -keyout self-signed.key -subj "/C=US/ST=NC/O=David Greeson/CN=${IP}" -out self-signed.csr
openssl x509 -req -extfile <(printf "\n[SAN]\nsubjectAltName=DNS:greeson.xyz,IP.1:${IP}") -days 1825 -in self-signed.csr -CA self-signed-ca.crt -CAkey self-signed-ca.key -CAcreateserial -out self-signed.crt
