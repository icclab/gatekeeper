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

# After installing gatekeeper, simply run this script to test if all API endpoints 
# are working properly or not. Make sure that Gatekeeper has been srated before
# runnig this test script. Change APIPATH to appropriate value.

APIPATH="localhost:8000"
echo "---- testing the discovery api ----"
curl -X GET "$APIPATH/" --header "Content-Type:application/json"
echo "---- end of the discovery api ----"
echo
echo "---- testing token generation (default admin) ----"
curl -X POST "$APIPATH/token/" --header "Content-Type:application/json" --header "X-Auth-Password:Eq7K8h9gpg" --header "X-Auth-Uid:1"
echo "---- end of admin token generation ----"
echo
echo "---- extracting the admin token value ----"
TOKEN=`curl -X POST "$APIPATH/token/" --header "Content-Type:application/json" --header "X-Auth-Password:Eq7K8h9gpg" --header "X-Auth-Uid:1" | python -mjson.tool | echo "{$(grep id)}" | python -c 'import sys, json; print json.load(sys.stdin)[sys.argv[1]]' id`
echo "---- admin token extracted from response ---"
echo $TOKEN
echo "---- end of token extraction process ----"
echo
echo "---- testing the user-list api ----"
curl -X GET "$APIPATH/admin/user/" --header "Content-Type:application/json" --header "X-Auth-Token:$TOKEN"
echo "---- end of user-list api ----"
echo
TESTUSER="user-$(date +'%T')"
echo "---- testing user creation ----"
curl -X POST "$APIPATH/admin/user/" --header "Content-Type:application/json" --header "X-Auth-Token:$TOKEN" -d '{"username":"'"$TESTUSER"'","password":"somepass","isadmin":"n", "accesslist":"ALL"}'
echo "---- end of user creation test ----"
echo
echo "--- testing user details retrieval api ----"
curl -X GET "$APIPATH/admin/user/1" --header "Content-Type:application/json" --header "X-Auth-Token:$TOKEN"
echo "---- end of user details retrieval test ----"
echo
echo "---- testing user update ----"
curl -X PUT "$APIPATH/admin/user/1" --header "Content-Type:application/json" --header "X-Auth-Token:$TOKEN" -d '{"accesslist":"servicex,servicey,ALL","isadmin":"y"}'
echo "---- end of user update test ----"
echo
echo "---- start authentication test ----"
curl -X GET "$APIPATH/auth/1" --header "Content-Type:application/json" --header "X-Auth-Password:Eq7K8h9gpg"
echo "---- end of simple authentication test ----"
echo
echo "---- start token validation user perspective ----"
curl -X GET "$APIPATH/token/validate/$TOKEN" --header "Content-Type:application/json" --header "X-Auth-Uid:1"
echo "---- end of token validate test ----"
echo
echo "---- testing service list api ----"
curl -X GET "$APIPATH/admin/service/" --header "Content-Type:application/json" --header "X-Auth-Token:$TOKEN"
echo "---- end of service list test ----"