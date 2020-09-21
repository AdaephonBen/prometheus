// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/prometheus/util/testutil"
)

func TestAddScheme(t *testing.T) {
	url := "google.com"
	urlWithScheme, err := addScheme(url)
	testutil.Ok(t, err)
	if urlWithScheme != "http://google.com" {
		t.Errorf("unexpected value %s for urlWithScheme", urlWithScheme)
	}
	url = "http://prometheus.io"
	urlWithScheme, err = addScheme(url)
	testutil.Ok(t, err)
	if urlWithScheme != "http://prometheus.io" {
		t.Errorf("unexpected value %s for urlWithScheme", urlWithScheme)
	}
}

func TestQueryRange(t *testing.T) {
	s, getRequest := mockServer(200, `{"status": "success", "data": {"resultType": "matrix", "result": []}}`)
	defer s.Close()

	p := &promqlPrinter{}
	exitCode := QueryRange(s.URL, map[string]string{}, "up", "0", "300", 0, p)
	testutil.Equals(t, "/api/v1/query_range", getRequest().URL.Path)
	form := getRequest().Form
	testutil.Equals(t, "up", form.Get("query"))
	testutil.Equals(t, "1", form.Get("step"))
	testutil.Equals(t, 0, exitCode)

	exitCode = QueryRange(s.URL, map[string]string{}, "up", "0", "300", 10*time.Millisecond, p)
	testutil.Equals(t, "/api/v1/query_range", getRequest().URL.Path)
	form = getRequest().Form
	testutil.Equals(t, "up", form.Get("query"))
	testutil.Equals(t, "0.01", form.Get("step"))
	testutil.Equals(t, 0, exitCode)
}

func TestQueryInstant(t *testing.T) {
	s, getRequest := mockServer(200, `{"status": "success", "data": {"resultType": "vector", "result": []}}`)
	defer s.Close()

	p := &promqlPrinter{}
	exitCode := QueryInstant(s.URL, "up", "300", p)
	testutil.Equals(t, "/api/v1/query", getRequest().URL.Path)
	form := getRequest().Form
	testutil.Equals(t, "up", form.Get("query"))
	testutil.Equals(t, "300", form.Get("time"))
	testutil.Equals(t, 0, exitCode)
}

func mockServer(code int, body string) (*httptest.Server, func() *http.Request) {
	var req *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		req = r
		w.WriteHeader(code)
		fmt.Fprintln(w, body)
	}))

	f := func() *http.Request {
		return req
	}
	return server, f
}
