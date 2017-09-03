// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Source: https://thenewstack.io/make-a-restful-json-api-go/

package main

import (
	"net/http"
	"log"
	"os"
)

var (
	db *Database
	config *Config
	accessLog *MultiWriter
	errorLog *MultiWriter
)

func init() {

	var err error

	db, err = NewDatabase("mysql.json")

	if err != nil {
		log.Fatalf("%v", err)
	}

	config, err = NewConfig("config.json")

	if err != nil {
		log.Fatalf("%v", err)
	}

	accessLog = NewMultiWriter()
	errorLog = NewMultiWriter()

	if config.UseLogFiles() {
		if accessLogFile, errorLogFile, err := config.LogFileInfo(); err == nil {
			accessLog.AddFiles(accessLogFile)
			errorLog.AddFiles(errorLogFile)
		} else {
			log.Printf("%v", err)
		}
	}

	if config.UseSyslog() {
		proto, raddr, tag := config.SyslogInfo()
		accessLog.AddSyslog(proto, raddr, tag, AccessPriority)
		accessLog.AddSyslog(proto, raddr, tag, ErrorPriority)
	}

	if accessLog.Count() == 0 {
		accessLog.Add(os.Stdout)
	}

	if errorLog.Count() == 0 {
		errorLog.Add(os.Stderr)
	}
}

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

	accessLog.Close()
	errorLog.Close()
	db.Close()
}
