#!/bin/bash

 # https://smallstep.com/hello-mtls/doc/server/kafka

 openssl pkcs12 -export -in server.crt -inkey server.key -name myserver.internal.net > server.p12
 keytool -importkeystore -srckeystore server.p12 -destkeystore kafka.server.keystore.jks -srcstoretype pkcs12
 keytool -keystore kafka.server.truststore.jks -alias CARoot -import -file ca.crt


 # https://smallstep.com/hello-mtls/doc/client/kafka-cli


 openssl pkcs12 -export -in client.a.crt -inkey client.a.key -name clienta > client.a.p13
 keytool -importkeystore -srckeystore client.a.p12 -destkeystore kafka.client.a.keystore.jks -srcstoretype pkcs12 -alias clienta
 keytool -keystore kafka.client.a.truststore.jks -alias CARoot -import -file ca.crt

