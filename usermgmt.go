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
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
	"strings"
)

func UserDetailsHandler(out http.ResponseWriter, in *http.Request) {
	id := mux.Vars(in)["id"]
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
			userDetail := GetUserDetail(dbArg, "user", id)
			if userDetail != nil {
				var jsonbody = staticMsgs[14]
				jsonbody = strings.Replace(jsonbody, "xxx", userDetail[0], 1)
				jsonbody = strings.Replace(jsonbody, "yyy", userDetail[1], 1)
				jsonbody = strings.Replace(jsonbody, "zzz", userDetail[2], 1)
				out.WriteHeader(http.StatusOK) //200 status code
				fmt.Fprintln(out, jsonbody)
			} else {
				out.WriteHeader(http.StatusNotFound) //404 status code
				var jsonbody = staticMsgs[15]
				fmt.Fprintln(out, jsonbody)
			}
		} else {
			var jsonbody = staticMsgs[18]
			out.WriteHeader(http.StatusUnauthorized) //401 status code
			fmt.Fprintln(out, jsonbody)
		}
	}
	MyFileInfo.Println("Received request on URI:/admin/user/{id} GET for uid:", id)
}

func UserUpdateHandler(out http.ResponseWriter, in *http.Request) {
	id := mux.Vars(in)["id"]
	decoder := json.NewDecoder(in.Body)
	var u user_struct
	err := decoder.Decode(&u)
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
			if err != nil {
				out.WriteHeader(http.StatusBadRequest) //status 400 Bad Request
				var jsonbody = staticMsgs[1]
				fmt.Fprintln(out, jsonbody)
				MyFileInfo.Println("Received malformed request on URI:/admin/user/{id} PUT for uid:", id)
			} else if len(u.AdminFlag) == 0 && len(u.CapabilityList) == 0 {
				out.WriteHeader(http.StatusBadRequest) //status 400 Bad Request
				var jsonbody = staticMsgs[1]
				fmt.Fprintln(out, jsonbody)
				MyFileInfo.Println("Received malformed request on URI:/admin/user/{id} PUT for uid:", id)
			} else {
				status := 0
				if len(u.CapabilityList) == 0 {
					//update just the admin-flag
					status = UpdateUser(dbArg, "user", "isadmin", u.AdminFlag, id)
				} else if len(u.AdminFlag) == 0 {
					//update just the capability list
					status = UpdateUser(dbArg, "user", "capability", u.CapabilityList, id)
				} else {
					//update both the fields
					status = UpdateUser(dbArg, "user", "isadmin", u.AdminFlag, id)
					status = UpdateUser(dbArg, "user", "capability", u.CapabilityList, id)
				}
				var jsonbody = ""
				if status == 1 {
					jsonbody = staticMsgs[16]
					out.WriteHeader(http.StatusOK) //200 status code
				} else {
					jsonbody = staticMsgs[17]
					out.WriteHeader(http.StatusNotModified) //304 status code
				}
				fmt.Fprintln(out, jsonbody)
				MyFileInfo.Println("Received request on URI:/admin/user/{id} PUT for uid:", id)
			}
		} else {
			var jsonbody = staticMsgs[18]
			out.WriteHeader(http.StatusUnauthorized) //401 status code
			fmt.Fprintln(out, jsonbody)
			MyFileInfo.Println("Received unauthorized request on URI:/admin/user/{id} PUT for uid:", id)
		}
	}

}

func UserDeleteHandler(out http.ResponseWriter, in *http.Request) {
	out.Header().Set("Content-Type", "application/json")
}

func UserListHandler(out http.ResponseWriter, in *http.Request) {
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
			userList := GetUserList(dbArg, "user", "username")
			var jsonbody = staticMsgs[4]
			var buffer bytes.Buffer
			for i := 0; i < len(userList); i++ {
				if i == 0 {
					buffer.WriteString("\"")
					buffer.WriteString(userList[i])
					buffer.WriteString("\"")
				} else {
					buffer.WriteString(",\"")
					buffer.WriteString(userList[i])
					buffer.WriteString("\"")
				}
			}
			jsonbody = strings.Replace(jsonbody, "xxx", buffer.String(), 1)
			out.WriteHeader(http.StatusOK) //200 status code
			fmt.Fprintln(out, jsonbody)
		} else {
			var jsonbody = staticMsgs[18]
			out.WriteHeader(http.StatusUnauthorized) //401 status code
			fmt.Fprintln(out, jsonbody)
		}
	}

	MyFileInfo.Println("Received request on URI:/admin/user/ GET")
}

func UserCreateHandler(out http.ResponseWriter, in *http.Request) {
	out.Header().Set("Content-Type", "application/json")
	if len(in.Header["X-Auth-Token"]) == 0 {
		MyFileWarning.Println("User Create Module - Can't Proceed: X-Auth-Token Missing!")
		out.WriteHeader(http.StatusBadRequest) //400 status code
		var jsonbody = staticMsgs[5]
		fmt.Fprintln(out, jsonbody)
	} else {
		token := in.Header["X-Auth-Token"][0]
		//check if token is valid and belongs to an admin user
		isAdmin := CheckTokenAdmin(token)
		if isAdmin {
			decoder := json.NewDecoder(in.Body)
			var u user_struct
			err := decoder.Decode(&u)

			if err != nil {
				out.WriteHeader(http.StatusBadRequest) //status 400 Bad Request
				var jsonbody = staticMsgs[1]
				fmt.Fprintln(out, jsonbody)
				MyFileInfo.Println("Received malformed request on URI:/admin/user/ POST")
				panic(err)
			} else if len(u.Username) == 0 {
				MyFileInfo.Println("Received malformed request on URI:/admin/user/ POST")
				out.WriteHeader(http.StatusBadRequest)
				var jsonbody = staticMsgs[1] //status 400 Bad Request
				fmt.Fprintln(out, jsonbody)
			} else {
				MyFileInfo.Println("Received JSON: Struct value received for user [pass hidden]:", u.Username)
				userCount := GetCount(dbArg, "user", "username", u.Username)
				if userCount > 0 {
					MyFileInfo.Println("Duplicate user create request on URI:/admin/user/ POST")
					out.WriteHeader(http.StatusPreconditionFailed)
					var jsonbody = staticMsgs[2] //user already exists
					fmt.Fprintln(out, jsonbody)
				} else {
					//now store the new user in the table and return back the proper response
					MyFileInfo.Println("Attempting to store new user:", u.Username, "into the table.")
					status := InsertUser(dbArg, "user", u.Username, u.Password, u.AdminFlag, u.CapabilityList) //inserting capability now
					MyFileInfo.Println("Status of the attempt to store new user:", u.Username, "into the table was:", status)

					out.WriteHeader(http.StatusOK) //200 status code
					var jsonbody = staticMsgs[3]   //user user creation msg, replace with actual content for xxx and yyy
					uId := LocateUser(dbArg, "user", u.Username)
					MyFileInfo.Println("The new id for user:", u.Username, "is:", uId)
					//constructing the correct JSON response
					jsonbody = strings.Replace(jsonbody, "xxx", strconv.Itoa(uId), 1)
					jsonbody = strings.Replace(jsonbody, "yyy", strconv.Itoa(uId), 1)
					jsonbody = strings.Replace(jsonbody, "zzz", strconv.Itoa(uId), 1)
					fmt.Fprintln(out, jsonbody)
				}
			}
		} else {
			var jsonbody = staticMsgs[18]
			out.WriteHeader(http.StatusUnauthorized) //401 status code
			fmt.Fprintln(out, jsonbody)
		}
	}
	MyFileInfo.Println("Received request on URI:/admin/user/ POST")
}

func GetUserDetail(filePath string, tableName string, userId string) []string {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT username, isadmin, capability FROM tablename WHERE uid=val;"
	queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
	queryStmt = strings.Replace(queryStmt, "val", userId, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in user-detail method.")
		checkErr(err, 1, db)
	}
	defer rows.Close()
	var udetail []string
	if rows.Next() {
		var userName string
		var isAdmin string
		var capabilityList string
		err = rows.Scan(&userName, &isAdmin, &capabilityList)
		checkErr(err, 1, db)
		udetail = append(udetail, userName)
		udetail = append(udetail, isAdmin)
		udetail = append(udetail, capabilityList)
	} else {
		udetail = nil
	}
	return udetail
}

func GetUserList(filePath string, tableName string, columnName string) []string {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT column FROM tablename;"
	queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
	queryStmt = strings.Replace(queryStmt, "column", columnName, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in user-list method.")
		checkErr(err, 1, db)
	}
	defer rows.Close()
	var ulist []string
	for rows.Next() {
		var userName string
		err = rows.Scan(&userName)
		checkErr(err, 1, db)
		ulist = append(ulist, userName)
	}
	return ulist
}

func UpdateUser(filePath string, tableName string, columnName string, newValue string, userId string) int {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	queryStmt := "UPDATE table SET column = 'value' WHERE uid=filter;"
	queryStmt = strings.Replace(queryStmt, "table", tableName, 1)
	queryStmt = strings.Replace(queryStmt, "column", columnName, 1)
	queryStmt = strings.Replace(queryStmt, "value", newValue, 1)
	queryStmt = strings.Replace(queryStmt, "filter", userId, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	result, err := db.Exec(queryStmt)
	MyFileInfo.Println("User Update Operation Result for user-id:", userId, "is:", result)
	if err != nil {
		MyFileWarning.Println("Caught error in user-update method.")
		checkErr(err, 1, db)
		return 0
	} else {
		return 1
	}
}

func LocateUser(filePath string, tableName string, userName string) int {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT uid FROM tablename WHERE username='searchterm';"
	queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
	queryStmt = strings.Replace(queryStmt, "searchterm", userName, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in user-locate method.")
		checkErr(err, 1, db)
	}
	defer rows.Close()
	if rows.Next() {
		var userId int
		err = rows.Scan(&userId)
		checkErr(err, 1, db)
		return userId
	}

	return -1
}

func CheckUserAccess(filePath string, tableName string, uid string, shortcode string) bool {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	queryStmt := "SELECT capability FROM tablename WHERE uid=searchterm;"
	queryStmt = strings.Replace(queryStmt, "tablename", tableName, 1)
	queryStmt = strings.Replace(queryStmt, "searchterm", uid, 1)

	MyFileInfo.Println("SQLite3 Query:", queryStmt)

	rows, err := db.Query(queryStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in check-user-access method.")
		checkErr(err, 1, db)
	}
	defer rows.Close()
	if rows.Next() {
		var accessList string
		err = rows.Scan(&accessList)
		checkErr(err, 1, db)
		MyFileInfo.Println("Found access-list for user [uid:", uid, "]: ", accessList)
		services := strings.Split(accessList, ",")
		located := false
		for i := 0; i < len(services); i++ {
			if strings.HasPrefix(services[i], shortcode) && strings.HasSuffix(services[i], shortcode) {
				located = true
			}
			if strings.HasPrefix(services[i], "ALL") && strings.HasSuffix(services[i], "ALL") {
				located = true
			}
			if strings.HasPrefix(services[i], "*") && strings.HasSuffix(services[i], "*") {
				located = true
			}
		}
		if located {
			return true
		}
	}
	return false
}

func InsertUser(filePath string, tableName string, userName string, passWord string, isAdmin string, capability string) bool {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		checkErr(err, 1, db)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	insertStmt := "INSERT INTO tablename VALUES (NULL, 'username', 'passhash', 'isadmin', 'capa');"
	insertStmt = strings.Replace(insertStmt, "tablename", tableName, 1)
	insertStmt = strings.Replace(insertStmt, "username", userName, 1)
	insertStmt = strings.Replace(insertStmt, "capa", capability, 1)
	data := []byte(passWord)
	hash := sha1.Sum(data)
	sha1hash := hex.EncodeToString(hash[:])
	//sha1hash := string(hash[:])
	MyFileInfo.Println("SHA-1 Hash Generated for the incoming password:", sha1hash)

	insertStmt = strings.Replace(insertStmt, "passhash", sha1hash, 1)
	insertStmt = strings.Replace(insertStmt, "isadmin", isAdmin, 1)
	MyFileInfo.Println("SQLite3 Query:", insertStmt)

	res, err := db.Exec(insertStmt)
	if err != nil {
		MyFileWarning.Println("Caught error in insert-user method,", res)
		checkErr(err, 1, db)
	}

	return true
}
