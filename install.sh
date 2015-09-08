#!/bin/bash

# Copyright (c) 2015. Zuercher Hochschule fuer Angewandte Wissenschaften
# All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.
#
# Author: Piyush Harsh,
# URL: piyush-harsh.info

sudo apt-get update
sudo apt-get upgrade
sudo apt-get install -y gcc
sudo apt-get install -y git
sudo apt-get install -y uuid-runtime

cd $HOME

echo "downloading and installing go runtime from google servers, please wait ..."
wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz

sudo -k

echo "configuring your environment for go projects ..."

cat >> $HOME/.profile << EOF

export GOPATH=$HOME/go
EOF

cd $HOME
source .profile

cat >> $HOME/.profile << EOF
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
EOF

cd $HOME

source .profile

echo "done."

echo "testing new environment variables ..."
echo "GOPATH: $GOPATH"
echo "Path: $PATH"

echo "Downloading auth-utils code now, please wait ..."
mkdir $HOME/go
mkdir -p $HOME/go/src/github.com/piyush82
cd $HOME/go/src/github.com/piyush82
git clone https://github.com/piyush82/auth-utils.git
echo "done."

cd auth-utils
echo "getting all code dependencies for auth-utils now, be patient ~ 1-2 minutes"
go get
echo "done."

echo "compiling and installing the package"
go install
echo "done."

cd

echo "starting the auth-service next, you can start using it at port :8000"
echo "use Ctrl+c to stop it. The executable is located at: $GOPATH/bin/auth-utils"

auth-utils