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
	"bufio"
	"encoding/json"
	"io"
	"os"

	"go.uber.org/zap"
)

func main() {
	l := setupLogging()
	defer l.Sync()

	zap.S().Infow("Test starting", "action", "TestStart")

	//CreateFile("/tmp/foo")
	//ModifyFile("/tmp/foo", "item one.")
	//ModifyFile("/tmp/foo", "item two.")
	//DeleteFile("/tmp/foo")
	//RunCommand("/bin/ls", []string{"/foo"})
	//
	//NetworkWrite("tcp", "blog.flame.org", 80, []byte("This is a test"))
	//NetworkWrite("udp", "blog.flame.org", 80, []byte("This is a test"))

	readLines(os.Stdin)

	zap.S().Infow("Test ended", "action", "TestEnd")
}

func readLines(f io.Reader) {
	scanner := bufio.NewScanner(f)
	block := ""
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			if len(block) > 0 {
				process(block)
			}
			block = ""
		} else {
			block += text
		}
	}
	err := scanner.Err()
	if err != nil {
		zap.S().Panicw("readLines", "error", err)
	}
	if len(block) > 0 {
		process(block)
	}
}

type actionJSON struct {
	Action string `json:"action,omitempty"`
}

func process(s string) {
	var action actionJSON
	err := json.Unmarshal([]byte(s), &action)
	if err != nil {
		zap.S().Panicw("process",
			"error", "unable to process JSON",
			"content", s)
	}

	switch action.Action {
	case "CreateFile":
		createFile(s)
	case "ModifyFile":
		modifyFile(s)
	case "DeleteFile":
		deleteFile(s)
	case "RunCommand":
		runCommand(s)
	case "NetworkWrite":
		networkWrite(s)
	default:
		zap.S().Panicw("process",
			"error", "unknown action",
			"content", s)
	}
}

type fileJSON struct {
	Path    string `json:"path,omitempty"`
	Content string `json:"content,omitempty"`
}

func createFile(s string) {
	var fj fileJSON
	err := json.Unmarshal([]byte(s), &fj)
	if err != nil {
		zap.S().Panicw("createfFile",
			"error", "unable to process JSON",
			"content", s)
	}

	CreateFile(fj.Path)
}

func modifyFile(s string) {
	var fj fileJSON
	err := json.Unmarshal([]byte(s), &fj)
	if err != nil {
		zap.S().Panicw("modifyFile",
			"error", "unable to process JSON",
			"content", s)
	}

	ModifyFile(fj.Path, fj.Content)
}

func deleteFile(s string) {
	var fj fileJSON
	err := json.Unmarshal([]byte(s), &fj)
	if err != nil {
		zap.S().Panicw("deleteFile",
			"error", "unable to process JSON",
			"content", s)
	}

	DeleteFile(fj.Path)
}

type commandJSON struct {
	Path string   `json:"path,omitempty"`
	Args []string `json:"args,omitempty"`
}

func runCommand(s string) {
	var j commandJSON
	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		zap.S().Panicw("runCommand",
			"error", "unable to process JSON",
			"content", s)
	}

	RunCommand(j.Path, j.Args)
}

type networkWriteJSON struct {
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Data     string `json:"data,omitempty"`
}

func networkWrite(s string) {
	var j networkWriteJSON
	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		zap.S().Panicw("networkWrite",
			"error", "unable to process JSON",
			"content", s)
	}

	NetworkWrite(j.Protocol, j.Host, j.Port, []byte(j.Data))
}
