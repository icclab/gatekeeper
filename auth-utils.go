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
	"github.com/scalingdata/gcfg"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Gatekeeper struct {
		Port    string
		LogFile string
		DbFile  string
	}
	Tnova struct {
		Defaultadmin  string
		Adminpassword string
	}
}

type user_struct struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	AdminFlag      string `json:"isadmin"`
	CapabilityList string `json:"accesslist"`
}

type service_struct struct {
	Shortname   string `json:"shortname"`
	Description string `json:"description"`
}

var (
	Trace         *log.Logger
	Info          *log.Logger
	Warning       *log.Logger
	Error         *log.Logger
	MyFileTrace   *log.Logger
	MyFileInfo    *log.Logger
	MyFileWarning *log.Logger
	MyFileError   *log.Logger
	staticMsgs    [20]string
	cfg           Config
	dbArg         string
)

func main() {
	err := gcfg.ReadFileInto(&cfg, "gatekeeper.cfg")
	if err != nil {
		log.Fatalf("Failed to parse gcfg data: %s", err)
		os.Exit(1)
	}
	file, err := os.OpenFile(cfg.Gatekeeper.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", "auth-utils.log", ":", err)
	}
	multi := io.MultiWriter(file, ioutil.Discard)
	Initlogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, multi)
	//logger has been initialized at this point
	InitMsgs()
	dbArg = "file:foo.db?cache=shared&mode=rwc"
	dbArg = strings.Replace(dbArg, "foo.db", cfg.Gatekeeper.DbFile, 1)
	dbCheck := CheckDB(dbArg)

	if dbCheck {
		MyFileInfo.Println("Table already exists in DB, nothing to do, proceeding normally.")
	} else {
		InitDB(dbArg)
	}
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", HomeHandler)
	users := r.Path("/admin/user/").Subrouter()
	users.Methods("GET").HandlerFunc(UserListHandler)
	users.Methods("POST").HandlerFunc(UserCreateHandler)

	user := r.Path("/admin/user/{id}").Subrouter()
	user.Methods("GET").HandlerFunc(UserDetailsHandler)
	user.Methods("PUT").HandlerFunc(UserUpdateHandler)
	user.Methods("DELETE").HandlerFunc(UserDeleteHandler)

	auth := r.Path("/auth/{id}").Subrouter()
	auth.Methods("GET").HandlerFunc(UserAuthHandler)

	tokens := r.Path("/token/").Subrouter()
	tokens.Methods("POST").HandlerFunc(TokenGenHandler)

	tokenval := r.Path("/token/validate/{id}").Subrouter()
	tokenval.Methods("GET").HandlerFunc(TokenValidateHandler)

	services := r.Path("/admin/service/").Subrouter()
	services.Methods("GET").HandlerFunc(ServiceListHandler)
	services.Methods("POST").HandlerFunc(ServiceRegisterHandler)

	portArg := ":port"
	portArg = strings.Replace(portArg, "port", cfg.Gatekeeper.Port, 1)
	MyFileInfo.Println("Starting server on", portArg)
	http.ListenAndServe(portArg, r)
	MyFileInfo.Println("Stopping server on", portArg)
}

func HomeHandler(out http.ResponseWriter, in *http.Request) {
	out.Header().Set("Content-Type", "application/json")
	out.WriteHeader(http.StatusOK) //200 status code
	var jsonbody = staticMsgs[0]
	fmt.Fprintln(out, jsonbody)
	MyFileInfo.Println("Received request on URI:/ GET")
}
