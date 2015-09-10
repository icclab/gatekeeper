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
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os/exec" //for calling system call uuidgen
	"strconv"
	"strings"
	"time"
)

func GetCount(filePath string, tableName string, columnName string, searchTerm string) int {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT count(*) FROM tablename WHERE columnname='searchterm';"
	queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
	queryStmt = strings.Replace(queryStmt, "columnname", columnName, 1)
	queryStmt = strings.Replace(queryStmt, "searchterm", searchTerm, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in count method.")
		checkErr(err, 1, db)
	}
	defer rows.Close()
	if rows.Next() {
		var userCount int
		err = rows.Scan(&userCount)
		checkErr(err, 1, db)
		return userCount
	}
	return 1
}

func CheckTokenAdmin(token string) bool {
	validity, uid := LocateTokenValidity(dbArg, "token", token)
	x, _ := strconv.ParseInt(validity, 10, 64)
	storedTime := time.Unix(x, 0)
	MyFileInfo.Println("Result of search for token[", token, "] was: Unix-validity", storedTime.String(), " for retrieved user-id:", uid)
	if time.Now().Before(storedTime) {
		//check if uid is an admin
		userDetail := GetUserDetail(dbArg, "user", strconv.Itoa(uid))
		if userDetail != nil {
			isAdmin := userDetail[1]
			if strings.HasPrefix(isAdmin, "y") && strings.HasSuffix(isAdmin, "y") {
				MyFileInfo.Println("Result of admin check for [", token, "] was: true.")
				return true
			} else {
				//not an admin
				MyFileInfo.Println("Result of admin check for [", token, "] was: false.")
				return false
			}
		} else {
			MyFileInfo.Println("Result of admin check for [", token, "] was: false. User details not found!")
			return false
		}
	} else {
		//expired token
		MyFileInfo.Println("Result of admin check for [", token, "] was: false. Token was expired or invalid!")
		return false
	}
}

func genuuid() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	uuid := string(out[:])
	return uuid
}

func checkErr(err error, errorType int, db *sql.DB) {
	if err != nil {
		MyFileError.Println("Unrecoverable Error!", err)
		panic(err)
	}
}
