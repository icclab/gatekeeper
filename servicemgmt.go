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
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"encoding/json"
	"strconv"
	"bytes"
)

func ServiceRegisterHandler(out http.ResponseWriter, in *http.Request) {
	out.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(in.Body)
	var s service_struct   
    err := decoder.Decode(&s)
	
	if len(in.Header["X-Auth-Token"]) == 0 {
		MyFileWarning.Println("User List Module - Can't Proceed: Token Missing!")
		out.WriteHeader(http.StatusBadRequest) //400 status code
		var jsonbody = staticMsgs[5]
		fmt.Fprintln(out, jsonbody)
	} else {
		token := in.Header["X-Auth-Token"][0]
		//check if token is valid and belongs to an admin user
		isAdmin := CheckTokenAdmin(token)
		if isAdmin {
			if err != nil {
    			out.WriteHeader(http.StatusBadRequest) //status 400 Bad Request
    			var jsonbody = staticMsgs[1]
				fmt.Fprintln(out, jsonbody)
				MyFileInfo.Println("Received malformed request on URI:/admin/service/ POST")
    		    panic(err)
    		} else if len(s.Shortname) == 0 {
				MyFileInfo.Println("Received malformed request on URI:/admin/service/ POST")
    			out.WriteHeader(http.StatusBadRequest)
    			var jsonbody = staticMsgs[1] //status 400 Bad Request
				fmt.Fprintln(out, jsonbody)
			} else {
				MyFileInfo.Println("Received JSON: Struct value received for service [shortname, description]:", s.Shortname, ",", s.Description)
				serviceCount := GetCount(dbArg, "service", "shortname", s.Shortname)
				if serviceCount > 0 {
    				MyFileInfo.Println("Duplicate service create request on URI:/admin/service/ POST")
    				out.WriteHeader(http.StatusPreconditionFailed)
    				var jsonbody = staticMsgs[12] //service already exists
					fmt.Fprintln(out, jsonbody)
    			} else {
					//now store the new service in the table and return back the proper response
    				MyFileInfo.Println("Attempting to store new service:", s.Shortname, "into the table.")
					uuid := genuuid()
					uuid = strings.TrimSpace(uuid)
					MyFileInfo.Println("Generated a new uuid for the service:", uuid)
					MyFileInfo.Println("Attempting to store new service:", s.Shortname, "into the table.")
					status := InsertService(dbArg, "service", uuid, s.Shortname, s.Description)
    				MyFileInfo.Println("Status of the attempt to store new service:", s.Shortname, "into the table was:", status)

    				out.WriteHeader(http.StatusOK) //200 status code
    				var jsonbody = staticMsgs[13] //service registration msg, replace with actual content for xxx and yyy
					sId := LocateService(dbArg, "service", uuid)
					MyFileInfo.Println("The new id for service:", s.Shortname, "is:", sId)
			
    				//constructing the correct JSON response
    				jsonbody = strings.Replace(jsonbody, "xxx", strconv.Itoa(sId), 1)
    				jsonbody = strings.Replace(jsonbody, "yyy", uuid, 1)
					fmt.Fprintln(out, jsonbody)
				}
			}
		} else {
			var jsonbody = staticMsgs[18]
			out.WriteHeader(http.StatusUnauthorized) //401 status code
			fmt.Fprintln(out, jsonbody)
		}
	}
    	
	MyFileInfo.Println("Received request on URI:/admin/service/ POST")
}

func ServiceListHandler(out http.ResponseWriter, in *http.Request) {
	out.Header().Set("Content-Type", "application/json")
	if len(in.Header["X-Auth-Token"]) == 0 {
		MyFileWarning.Println("User List Module - Can't Proceed: Token Missing!")
		out.WriteHeader(http.StatusBadRequest) //400 status code
		var jsonbody = staticMsgs[5]
		fmt.Fprintln(out, jsonbody)
	} else {
		token := in.Header["X-Auth-Token"][0]
		//check if token is valid and belongs to an admin user
		isAdmin := CheckTokenAdmin(token)
		if isAdmin {
			uuidlist, snamelist := GetServiceList(dbArg, "service")
			var jsonbody = staticMsgs[11]
			var buffer1 bytes.Buffer
			var buffer2 bytes.Buffer
	
			for i := 0; i < len(uuidlist); i++ {
				if i == 0 {
					buffer1.WriteString("\"")
					buffer1.WriteString(uuidlist[i])
					buffer1.WriteString("\"")
					buffer2.WriteString("\"")
					buffer2.WriteString(snamelist[i])
					buffer2.WriteString("\"")
				} else {
					buffer1.WriteString(",\"")
					buffer1.WriteString(uuidlist[i])
					buffer1.WriteString("\"")
					buffer2.WriteString(",\"")
					buffer2.WriteString(snamelist[i])
					buffer2.WriteString("\"")
				}
			}
			jsonbody = strings.Replace(jsonbody, "uuid-xxx", buffer1.String(), 1)
			jsonbody = strings.Replace(jsonbody, "name-yyy", buffer2.String(), 1)
			out.WriteHeader(http.StatusOK) //200 status code
			fmt.Fprintln(out, jsonbody)
		} else {
			var jsonbody = staticMsgs[18]
			out.WriteHeader(http.StatusUnauthorized) //401 status code
			fmt.Fprintln(out, jsonbody)
		}
	}
	
	MyFileInfo.Println("Received request on URI:/admin/service/ GET")
}

func GetServiceList(filePath string, tableName string) ([]string, []string) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
        checkErr(err, 1, db)
    }
    defer db.Close()
    
    err = db.Ping()
	if err != nil {
    	panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT key, shortname FROM tablename;"
	queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
    if err != nil {
    	MyFileWarning.Println("Caught error in service-list method.")
    	checkErr(err, 1, db)
    }
    defer rows.Close()
    var uuidlist []string
	var namelist []string
    for rows.Next() {
    	var suuid string
		var sname string
        err = rows.Scan(&suuid, &sname)
        checkErr(err, 1, db)
        uuidlist = append(uuidlist, suuid)
		namelist = append(namelist, sname)
    }
    return uuidlist, namelist
}

func InsertService(filePath string, tableName string, serviceUuid string, serviceName string, serviceDesc string) bool {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
        checkErr(err, 1, db)
    }
    defer db.Close()
    
    err = db.Ping()
	if err != nil {
    	panic(err.Error()) // proper error handling instead of panic in your app
	}

    insertStmt := "INSERT INTO tablename VALUES (NULL, 'skey', 'sname', 'sdesc');"
    insertStmt = strings.Replace(insertStmt, "tablename", tableName, 1)
    insertStmt = strings.Replace(insertStmt, "skey", serviceUuid, 1)
	insertStmt = strings.Replace(insertStmt, "sname", serviceName, 1)
	insertStmt = strings.Replace(insertStmt, "sdesc", serviceDesc, 1)
    MyFileInfo.Println("SQLite3 Query:", insertStmt)

    res, err := db.Exec(insertStmt)
    if err != nil {
    	MyFileWarning.Println("Caught error in insert-service method,", res)
    	checkErr(err, 1, db)
    }
    
	return true
}

func LocateService(filePath string, tableName string, uuid string) int {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
        checkErr(err, 1, db)
    }
    defer db.Close()
    
    err = db.Ping()
	if err != nil {
    	panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT sid FROM tablename WHERE key='searchterm';"
    queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
    queryStmt = strings.Replace(queryStmt, "searchterm", uuid, 1)
    
    MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
    if err != nil {
    	MyFileWarning.Println("Caught error in service-locate method.")
    	checkErr(err, 1, db)
    }
    defer rows.Close()
    if rows.Next() {
    	var serviceId int
        err = rows.Scan(&serviceId)
        checkErr(err, 1, db)
        return serviceId
    }
    
	return -1
}

func LocateServiceCode(filePath string, tableName string, uuId string) string {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
        checkErr(err, 1, db)
    }
    defer db.Close()
    
    err = db.Ping()
	if err != nil {
    	panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT shortname FROM tablename WHERE key='searchterm';"
    queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
    queryStmt = strings.Replace(queryStmt, "searchterm", uuId, 1)
    
    MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
    if err != nil {
    	MyFileWarning.Println("Caught error in locate-service-code method.")
    	checkErr(err, 1, db)
    }
    defer rows.Close()
    if rows.Next() {
    	var shortcode string
        err = rows.Scan(&shortcode)
        checkErr(err, 1, db)
        return shortcode
    }
    
	return ""
}