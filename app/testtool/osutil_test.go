/*
 * Copyright 2022 Michael Graff.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrentUser(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// likely to fail on Windows, for for now this will work on Unix-like systems
		{"returns same as environment", os.Getenv("USER")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CurrentUser(); got != tt.want {
				t.Errorf("CurrentUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunCommand(t *testing.T) {
	type args struct {
		path string
		args []string
	}
	tests := []struct {
		name        string
		args        args
		expectPanic bool
	}{
		{"/bin/ls", args{"/bin/ls", []string{"/"}}, false},
		{"/bin/lsxxxxxasda", args{"/bin/lsxxxxxasda", []string{"/"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					RunCommand(tt.args.path, tt.args.args)
				})
			} else {
				assert.NotPanics(t, func() {
					RunCommand(tt.args.path, tt.args.args)
				})
			}
		})
	}
}

func TestCurrentExecutable(t *testing.T) {
	t.Run("returns something", func(t *testing.T) {
		assert.NotPanics(t, func() {
			got := CurrentExecutable()
			assert.NotEmpty(t, got)
		})
	})
}

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectPanic bool
	}{
		{"/tmp/foo", "/tmp/foo", false},
		{"/DoesnOTeXiST", "/DoesnOTeXiST", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					CreateFile(tt.path)
				})
			} else {
				assert.NotPanics(t, func() {
					CreateFile(tt.path)
				})
				assert.FileExists(t, tt.path)
				assert.NoError(t, os.Remove(tt.path))
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		path        string
		expectPanic bool
	}{
		{"/tmp/foo2", false},
		{"/DoesnOTeXiST", true},
	}
	os.Create("/tmp/foo2") // hack to make sure this exists
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					DeleteFile(tt.path)
				})
			} else {
				assert.NotPanics(t, func() {
					DeleteFile(tt.path)
				})
			}
		})
	}
}
