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
		zap.S().Panicw("user.Current()", "error", err)
	}
	return n.Username
}

// CurrentExecutable retuns the full path of the currently running executable,
// or panics if it cannot be retrieved from os.Executable()
func CurrentExecutable() string {
	n, err := os.Executable()
	if err != nil {
		zap.S().Panicw("os.Executable()", "error", err)
	}
	return n
}

// RunCommand runs the command with arguments.  If an error occurs running
// the command, it will panic.
func RunCommand(path string, args []string) {
	cmd := exec.Command(path, args...)
	err := cmd.Run()
	if err != nil && cmd.ProcessState.ExitCode() == 0 {
		zap.S().Panicw("RunCommand", "error", err)
	}
	zap.S().Infow("RunCommand",
		"action", "RunCommand",
		"cmdPath", path,
		"cmdArgs", args,
		"cmdPID", cmd.ProcessState.Pid(),
		"cmdExitStatus", cmd.ProcessState.ExitCode())
}

// CreateFile will create an empty file using the 0644 umask at
// the specified path.  All directories need to exist prior to calling.
// The file is not deleted automatically.
// If the file already exists (and can be opened for reading at least)
// this function will panic.
// If the file cannot be created, it will panic.
func CreateFile(path string) {
	f, err := os.Open(path)
	if err == nil {
		f.Close()
		zap.S().Panicw("CreateFile", "error", "file already exists")
	}
	f, err = os.OpenFile(path, os.O_CREATE, 0644)
	if err != nil {
		zap.S().Panicw("CreateFile", "error", err)
	}
	f.Close()
	zap.S().Infow("CreateFile",
		"action", "CreateFile",
		"fileAction", "create",
		"path", path,
	)
}

// DeleteFile will delete the file at the specified path.
// If an error occurs, it will panic.
func DeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		zap.S().Panicw("DeleteFile", "error", err)
	}
	zap.S().Infow("DeleteFile",
		"action", "DeleteFile",
		"fileAction", "delete",
		"path", path,
	)
}

// ModifyFile will add some text to the end of a file.  The file must
// already exist.
// If an error occurs, it will panic.
func ModifyFile(path string, content string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		zap.S().Panicw("ModifyFile", "error", err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		zap.S().Panicw("ModifyFile", "error", err)
	}
	zap.S().Infow("ModifyFile",
		"action", "ModifyFile",
		"fileAction", "modify",
		"path", path,
	)
}
