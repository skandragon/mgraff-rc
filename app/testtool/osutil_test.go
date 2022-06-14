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
		name string
		args args
	}{
		{"/bin/ls", args{"/bin/ls", []string{"/"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunCommand(tt.args.path, tt.args.args)
		})
	}
}

func TestCurrentExecutable(t *testing.T) {
	t.Run("returns something", func(t *testing.T) {
		got := CurrentExecutable()
		assert.NotEmpty(t, got)
	})
}
