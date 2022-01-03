#!/bin/bash

pushd `dirname $0` > /dev/null
SCRIPTPATH=`pwd -P`
popd > /dev/null
SCRIPTFILE=`basename $0`

cd ${SCRIPTPATH}/certs
rm -f *.key *.crt *.srl *.csr *.log
