/*
 * Copyright (c) 2015. Zuercher Hochschule fuer Angewandte Wissenschaften
 *  All Rights Reserved.
 *
 *     Licensed under the Apache License, Version 2.0 (the "License"); you may
 *     not use this file except in compliance with the License. You may obtain
 *     a copy of the License at
 *
 *          http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 *     WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 *     License for the specific language governing permissions and limitations
 *     under the License.
 */

/*
 *     Author: Piyush Harsh,
 *     URL: piyush-harsh.info
 */

package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func TokenValidateHandler(out http.ResponseWriter, in *http.Request) {
	out.Header().Set("Content-Type", "application/json")
	id := mux.Vars(in)["id"]

	if (len(in.Header["X-Auth-Service-Key"]) == 0 && len(in.Header["X-Auth-Uid"]) == 0) || len(id) == 0 {
		MyFileWarning.Println("Token Validation Module - Can't Proceed: Service-Key or UserID Missing!")
		out.WriteHeader(http.StatusBadRequest) //400 status code
		var jsonbody = staticMsgs[5]
		fmt.Fprintln(out, jsonbody)
	} else {
		if len(in.Header["X-Auth-Service-Key"]) > 0 {
			//not handling the service-key header for now
			serviceKey := in.Header["X-Auth-Service-Key"][0]
			MyFileInfo.Println("Received token authorization request for service key[", serviceKey, "] against token [", id, "]")
			//get user id from this token
			validity, uid := LocateTokenValidity(dbArg, "token", id)
			x, _ := strconv.ParseInt(validity, 10, 64)
			storedTime := time.Unix(x, 0)
			MyFileInfo.Println("Result of search for token[", id, "] was: Unix-validity", storedTime.String(), "user-id:", uid)
			if time.Now().Before(storedTime) {
				//token is valid, now check if the token is authorized
				//locate service short code now
				shortCode := LocateServiceCode(dbArg, "service", serviceKey)
				if len(shortCode) > 0 {
					MyFileInfo.Println("For service-key [", serviceKey, "]: short-code found: [", shortCode, "]")
					//now check if uid is allowed to access this service
					hasAccess := CheckUserAccess(dbArg, "user", strconv.Itoa(uid), shortCode)
					if hasAccess {
						MyFileInfo.Println("Token authorization result for service key[", serviceKey, "] against token [", id, "] was: true.")
						out.WriteHeader(http.StatusOK) //200 status code
						var jsonbody = staticMsgs[9]   //token is valid
						//constructing the correct JSON response
						fmt.Fprintln(out, jsonbody)
					} else {
						MyFileInfo.Println("Token authorization result for service key[", serviceKey, "] against token [", id, "] was: false.")
						//token is not authorized for this service
						out.WriteHeader(http.StatusNotAcceptable) //406 status code
						var jsonbody = staticMsgs[10]             //validation failed msg
						//constructing the correct JSON response
						fmt.Fprintln(out, jsonbody)
					}

				} else {
					//no such service found
					out.WriteHeader(http.StatusNotAcceptable) //406 status code
					var jsonbody = staticMsgs[10]             //validation failed msg
					//constructing the correct JSON response
					MyFileInfo.Println("Validation result for token: ", id, "was - unauthorized. No such service found.")
					fmt.Fprintln(out, jsonbody)
				}
			} else {
				//token is invalid
				out.WriteHeader(http.StatusNotAcceptable) //406 status code
				var jsonbody = staticMsgs[10]             //validation failed msg
				//constructing the correct JSON response
				MyFileInfo.Println("Validation result for token: ", id, "was - invalid. Token was expired.")
				fmt.Fprintln(out, jsonbody)
			}
		} else {
			userId := in.Header["X-Auth-Uid"][0]
			//locate validity of token from db
			MyFileInfo.Println("TokenValidate: Incoming User-ID in header:", userId)
			validity, uid := LocateTokenValidity(dbArg, "token", id)
			x, _ := strconv.ParseInt(validity, 10, 64)
			storedTime := time.Unix(x, 0)
			MyFileInfo.Println("Result of search for token[", id, "] was: Unix-validity", storedTime.String(), "user-id:", uid)
			//this matches the uid with the uid associated with the token in question
			if time.Now().Before(storedTime) && strings.HasPrefix(userId, strconv.Itoa(uid)) && strings.HasSuffix(userId, strconv.Itoa(uid)) {
				//token is valid
				out.WriteHeader(http.StatusOK) //200 status code
				var jsonbody = staticMsgs[9]   //token is valid
				//constructing the correct JSON response
				MyFileInfo.Println("Validation result for token: ", id, "was - valid.")
				fmt.Fprintln(out, jsonbody)
			} else {
				//token is invalid
				out.WriteHeader(http.StatusNotAcceptable) //406 status code
				var jsonbody = staticMsgs[10]             //user user creation msg, replace with actual content for xxx and yyy
				//constructing the correct JSON response
				MyFileInfo.Println("Validation result for token: ", id, "was - invalid. Either expired or user-id mismatch.")
				fmt.Fprintln(out, jsonbody)
			}
		}
	}
	MyFileInfo.Println("Received request on URI:/token/validate/{id} GET")
}

func TokenGenHandler(out http.ResponseWriter, in *http.Request) {
	out.Header().Set("Content-Type", "application/json")

	if len(in.Header["X-Auth-Token"]) == 0 && (len(in.Header["X-Auth-Uid"]) == 0 || len(in.Header["X-Auth-Password"]) == 0) {
		MyFileWarning.Println("Token Gen Module - Can't Proceed: Token / UserID and Password Missing!")
		out.WriteHeader(http.StatusBadRequest) //400 status code
		var jsonbody = staticMsgs[5]
		fmt.Fprintln(out, jsonbody)
	} else {
		//Two paths - if X-Auth-Token is present, this has precedence
		if len(in.Header["X-Auth-Token"]) > 0 {

		} else {
			//username & password is present
			passWord := in.Header["X-Auth-Password"][0]
			userId := in.Header["X-Auth-Uid"][0]
			data := []byte(passWord)
			hash := sha1.Sum(data)
			//sha1hash := string(hash[:])
			sha1hash := hex.EncodeToString(hash[:])
			MyFileInfo.Println("TokenGen - SHA-1 Hash Generated for the incoming password:", sha1hash)
			//get the stored SHA-1 Hash for the incoming user
			storedHash := LocatePasswordHash("file:foo.db?cache=shared&mode=rwc", "user", userId)
			MyFileInfo.Println("SHA-1 hash retrieved for the incoming user:", storedHash)
			if strings.HasPrefix(storedHash, sha1hash) && strings.HasSuffix(storedHash, sha1hash) {
				uuid := genuuid()
				uuid = strings.TrimSpace(uuid)
				MyFileInfo.Println("Generated a new uuid for the token:", uuid)
				//generate a validity period
				now := time.Now()
				validity := now.Add(6 * 60 * 60 * 1000000000) //added 6 hours of validity from now
				//set the time format
				//timeValue := validity.Format(time.UnixDate)
				//fmt.Println(strconv.FormatInt(validity.Unix(), 10))
				status := InsertToken(dbArg, "token", uuid, userId, strconv.FormatInt(validity.Unix(), 10), "*")
				MyFileInfo.Println("Status of the attempt to store new token:", uuid, "into the table was:", status)
				out.WriteHeader(http.StatusOK) //200 status code
				var jsonbody = staticMsgs[8]   //user user creation msg, replace with actual content for xxx and yyy
				//constructing the correct JSON response
				jsonbody = strings.Replace(jsonbody, "uuid-xxx", uuid, 1)
				jsonbody = strings.Replace(jsonbody, "time-yyy", validity.String(), 1)
				fmt.Fprintln(out, jsonbody)
			} else {
				out.WriteHeader(http.StatusUnauthorized) //401 status code
				var jsonbody = staticMsgs[6]
				fmt.Fprintln(out, jsonbody)
				MyFileInfo.Println("TokenGen - Password does not match. User", userId, "not authenticated.")
			}
		}
	}
	MyFileInfo.Println("Received request on URI:/token/ POST")
}

func GetTokenList(filePath string, tableName string, uid string) ([]string, []string) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT uuid, validupto FROM tablename WHERE uid=userid;"
	queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
	queryStmt = strings.Replace(queryStmt, "userid", uid, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in user-list method.")
		checkErr(err, 1, db)
	}
	defer rows.Close()
	var tokenlist []string
	var validitylist []string
	for rows.Next() {
		var token string
		var validity string
		err = rows.Scan(&token, &validity)
		checkErr(err, 1, db)
		tokenlist = append(tokenlist, token)
		validitylist = append(validitylist, validity)
	}
	return tokenlist, validitylist
}

func LocateTokenValidity(filePath string, tableName string, uuId string) (string, int) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT validupto, uid FROM tablename WHERE uuid='searchterm';"
	queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
	queryStmt = strings.Replace(queryStmt, "searchterm", uuId, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in token-validity-locate method.")
		checkErr(err, 1, db)
	}
	defer rows.Close()
	if rows.Next() {
		var validity string
		var uid int
		err = rows.Scan(&validity, &uid)
		checkErr(err, 1, db)
		return validity, uid
	}

	return "", -1
}

func InsertToken(filePath string, tableName string, uuid string, uid string, validupto string, capability string) bool {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	insertStmt := "INSERT INTO tablename VALUES ('uuid', uid, 'validto', 'capability');"
	insertStmt = strings.Replace(insertStmt, "tablename", tableName, 1)
	insertStmt = strings.Replace(insertStmt, "uuid", uuid, 1)
	insertStmt = strings.Replace(insertStmt, "uid", uid, 1)
	insertStmt = strings.Replace(insertStmt, "validto", validupto, 1)
	insertStmt = strings.Replace(insertStmt, "capability", capability, 1)

	MyFileInfo.Println("SQLite3 Query:", insertStmt)

	res, err := db.Exec(insertStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in insert-token method,", res)
		checkErr(err, 1, db)
	}

	return true
}
