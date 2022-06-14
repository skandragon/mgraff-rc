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
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getLocalAddress(t *testing.T) {
	tests := []struct {
		name        string
		addr        net.Addr
		want        string
		want1       int
		want2       string
		expectPanic bool
	}{
		{
			"tcp address",
			&net.TCPAddr{IP: net.IPv4(1, 2, 3, 4), Port: 4321, Zone: "foo"},
			"1.2.3.4", 4321, "foo",
			false,
		},
		{
			"udp address",
			&net.UDPAddr{IP: net.IPv4(1, 2, 3, 4), Port: 4321, Zone: "foo"},
			"1.2.3.4", 4321, "foo",
			false,
		},
		{
			"unix address",
			&net.UnixAddr{Name: "foo", Net: "bar"},
			"", 0, "",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					_, _, _ = getLocalAddress(tt.addr)
				})
			} else {
				got, got1, got2 := getLocalAddress(tt.addr)
				if got != tt.want {
					t.Errorf("getLocalAddress() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t.Errorf("getLocalAddress() got1 = %v, want %v", got1, tt.want1)
				}
				if got2 != tt.want2 {
					t.Errorf("getLocalAddress() got2 = %v, want %v", got2, tt.want2)
				}
			}
		})
	}
}
