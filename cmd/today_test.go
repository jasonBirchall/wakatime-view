/*
Copyright © 2022 jason birchall.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetWakaData(t *testing.T) {
	expected := "{'data': 'dummy'}"

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expected)
	}))

	defer svr.Close()

	c := NewClient(svr.URL, "dummy")
	res, err := c.getWakaData("fake")
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	res = strings.TrimSpace(res)

	if res != expected {
		t.Errorf("expected res to be %s got %s", expected, res)
	}
}
