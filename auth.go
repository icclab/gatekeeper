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
	"bytes"
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

func UserAuthHandler(out http.ResponseWriter, in *http.Request) {
	id := mux.Vars(in)["id"]
	var passWord string
	out.Header().Set("Content-Type", "application/json")

	if len(in.Header["X-Auth-Password"]) == 0 {
		MyFileWarning.Println("Authentication Module - Can't Proceed: Password Missing! UserId =", id)
		out.WriteHeader(http.StatusBadRequest) //400 status code
		var jsonbody = staticMsgs[5]
		fmt.Fprintln(out, jsonbody)
	} else {
		passWord = in.Header["X-Auth-Password"][0]
		MyFileInfo.Println("A valid password [password hidden] received for userid:", id)
		data := []byte(passWord)
		hash := sha1.Sum(data)
		//sha1hash := string(hash[:])
		sha1hash := hex.EncodeToString(hash[:])
		MyFileInfo.Println("SHA-1 Hash Generated for the incoming password:", sha1hash)
		//get the stored SHA-1 Hash for the incoming user
		storedHash := LocatePasswordHash("file:foo.db?cache=shared&mode=rwc", "user", id)
		MyFileInfo.Println("SHA-1 hash retrieved for the incoming user:", storedHash)
		if strings.HasPrefix(storedHash, sha1hash) && strings.HasSuffix(storedHash, sha1hash) {
			out.WriteHeader(http.StatusAccepted) //202 status code
			var jsonbody = staticMsgs[7]
			//get the list of all valid tokens associated with this user
			tokenlist, validitylist := GetTokenList(dbArg, "token", id)
			var buffer1 bytes.Buffer
			var buffer2 bytes.Buffer
			runCount := 0
			for i := 0; i < len(validitylist); i++ {
				//check if the token is valid or not
				unixTimeStr := validitylist[i]
				x, _ := strconv.ParseInt(unixTimeStr, 10, 64)
				storedTime := time.Unix(x, 0)
				if time.Now().Before(storedTime) {
					if runCount == 0 {
						buffer1.WriteString("\"")
						buffer1.WriteString(tokenlist[i])
						buffer1.WriteString("\"")
						buffer2.WriteString("\"")
						buffer2.WriteString(storedTime.String())
						buffer2.WriteString("\"")
					} else {
						buffer1.WriteString(",\"")
						buffer1.WriteString(tokenlist[i])
						buffer1.WriteString("\"")
						buffer2.WriteString(",\"")
						buffer2.WriteString(storedTime.String())
						buffer2.WriteString("\"")
					}
					runCount += 1
				} else {
					// these tokens are expired so ignored
				}
			}
			jsonbody = strings.Replace(jsonbody, "uuid-xxx", buffer1.String(), 1)
			jsonbody = strings.Replace(jsonbody, "time-yyy", buffer2.String(), 1)
			fmt.Fprintln(out, jsonbody)
			MyFileInfo.Println("Password matches. User", id, "successfully authenticated.")
		} else {
			out.WriteHeader(http.StatusUnauthorized) //401 status code
			var jsonbody = staticMsgs[6]
			fmt.Fprintln(out, jsonbody)
			MyFileInfo.Println("Password does not match. User", id, "not authenticated.")
		}
	}

	MyFileInfo.Println("Received request on URI:/auth/{user-id} GET")
}

func LocatePasswordHash(filePath string, tableName string, userId string) string {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT password FROM tablename WHERE uid=searchterm;"
	queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
	queryStmt = strings.Replace(queryStmt, "searchterm", userId, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in user-password-locate method.")
		checkErr(err, 1, db)
	}
	defer rows.Close()
	if rows.Next() {
		var passWord string
		err = rows.Scan(&passWord)
		checkErr(err, 1, db)
		return passWord
	}

	return ""
}
