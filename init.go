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
	"io"
	"log"
	"strings"
)

func CheckDB(filePath string) bool {
	db, err := sql.Open("sqlite3", filePath)
	var status bool
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		checkErr(err, 0, db)
	} else {
		defer rows.Close()
		for rows.Next() {
			var tablename string
			err = rows.Scan(&tablename)
			checkErr(err, 1, db)
			MyFileInfo.Println("While performing DB sanity checks: found table", tablename)
			if len("user") == len(tablename) {
				if strings.Count(tablename, "user") == 1 {
					status = true
				}
			}
		}
	}
	return status
}

func InitDB(filePath string) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	} else {

		var dbCmd = `
    			CREATE TABLE 'user' (
    			'uid' INTEGER PRIMARY KEY AUTOINCREMENT,
    			'username' VARCHAR(64) NULL,
    			'password' VARCHAR(64) NULL,
    			'isadmin' VARCHAR(1),
				'capability' VARCHAR(128)
				);
    		`
		stmt, err := db.Prepare(dbCmd)
		checkErr(err, 1, db)
		res, err := stmt.Exec()
		checkErr(err, 1, db)
		dbCmd = `
    			CREATE TABLE 'token' (
    			'uuid' VARCHAR(128) PRIMARY KEY,
    			'uid' INTEGER,
    			'validupto' VARCHAR(64) NULL,
    			'capability' VARCHAR(128)
				);
			`
		stmt, err = db.Prepare(dbCmd)
		checkErr(err, 1, db)
		res, err = stmt.Exec()
		checkErr(err, 1, db)
		dbCmd = `
    			CREATE TABLE 'service' (
				'sid' INTEGER PRIMARY KEY AUTOINCREMENT,
    			'key' VARCHAR(128),
    			'shortname' VARCHAR(16) NULL,
    			'description' VARCHAR(128)
				);
			`
		stmt, err = db.Prepare(dbCmd)
		checkErr(err, 1, db)
		res, err = stmt.Exec()
		checkErr(err, 1, db)
		MyFileInfo.Println("Created tables user, token, service. System ready. System response:", res)
		status := InsertUser(dbArg, "user", cfg.Tnova.Defaultadmin, cfg.Tnova.Adminpassword, "y", "ALL")
		if status {
			MyFileInfo.Println("Status of the attempt to store default admin-user:", cfg.Tnova.Defaultadmin, "into the table was:", status)
		} else {
			MyFileWarning.Println("Status of the attempt to store default admin-user:", cfg.Tnova.Defaultadmin, "into the table was:", status)
		}
	}
}

func Initlogger(traceHandle, infoHandle, warningHandle, errorHandle, fileHandle io.Writer) {
	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	MyFileTrace = log.New(fileHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	MyFileInfo = log.New(fileHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	MyFileWarning = log.New(fileHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	MyFileError = log.New(fileHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
