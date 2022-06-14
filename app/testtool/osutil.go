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
	"os/exec"
	"os/user"

	"go.uber.org/zap"
)

// CurrentUser returns a string representing the current user running this
// process, or logs and exits if this cannot be determined.
func CurrentUser() string {
	n, err := user.Current()
	if err != nil {
		zap.S().Fatalw("user.Current()", "error", err)
	}
	return n.Username
}

func CurrentExecutable() string {
	n, err := os.Executable()
	if err != nil {
		zap.S().Fatalw("os.Executable()", "error", err)
	}
	return n
}

func RunCommand(path string, args []string) {
	cmd := exec.Command(path, args...)
	err := cmd.Run()
	if err != nil {
		zap.S().Fatalw("exec", "error", err)
	}
	zap.S().Infow("exec",
		"cmdPath", path,
		"cmdArgs", args,
		"cmdPID", cmd.ProcessState.Pid(),
		"cmdExitStatus", cmd.ProcessState.ExitCode())
}
