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

package api

import (
	`bytes`
	`encoding/json`
	`io`
	`io/ioutil`
	`net/http`
)

const HttpBodySizeLimit = 1048576

// ReadBody returns the HTTP request body.
func ReadBody(r *http.Request) ([]byte, error) {

	if body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit)); err != nil {
		return nil, err
	} else if err := r.Body.Close(); err != nil {
		return nil, err
	} else {
		return body, nil
	}
}
// DecodeBody unmarshals the JSON object in the HTTP request body to an object.
func DecodeBody(i interface{}, r *http.Request) (error) {

	if body, err := ReadBody(r); err != nil {
		return err
	} else if err := json.Unmarshal(body, &i); err != nil {
		return err
	} else {
		return nil
	}
}

// WriteBody creates a new HTTP request body in an existing HTTP request.
func WriteBody(r *http.Request, body []byte) {
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	r.ContentLength = int64(len(body))
}
