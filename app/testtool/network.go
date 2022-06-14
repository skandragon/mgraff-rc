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
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

// NetworkWrite will connect to a remote host and port over TCP (perhaps v4,
// perhaps v6) and write the specified data, then close the connection.
// No reply will be read.
// Only "tcp" and "udp" are supported, along with the special case variants
// of "tcp4", "tcp6", "udp4", and "udp6".
//
// Any error (connection, write timeout) will result in a panic.
// If the data is not fully written, no panic will occur, and the log
// message will indicate what has been written.
//
// Note this is all inherently unstable, as it would be better to use some
// server that sends a known response, so we can ensure everything has
// been sent as intended.  This mode is basically "fire and forget" and
// the only thing the Write() call ensures is that the data was written to
// the kernel network buffer.
func NetworkWrite(protocol string, host string, port int, data []byte) {
	const timeout = 10 * time.Second
	conn, err := net.DialTimeout(protocol, fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		zap.S().Panicw("NetworkWrite",
			"error", err,
			"protocol", protocol,
			"host", host,
			"port", port)
	}

	// Save away the local address and port, in case it doesn't last
	// after close.
	localAddress, localPort, localZone := getLocalAddress(conn.LocalAddr())

	conn.SetWriteDeadline(time.Now().Add(timeout))

	n, err := conn.Write(data)
	if err != nil {
		zap.S().Panicw("NetworkWrite",
			"error", err,
			"protocol", protocol,
			"localAddress", localAddress,
			"localPort", localPort,
			"localZone", localZone,
			"host", host,
			"port", port)
	}

	_ = conn.Close()

	zap.S().Infow("NetworkWrite",
		"host", host,
		"port", port,
		"protocol", protocol,
		"localAddress", localAddress,
		"localPort", localPort,
		"localZone", localZone,
		"nWritten", n)
}

func getLocalAddress(addr net.Addr) (string, int, string) {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if ok {
		return tcpAddr.IP.String(), tcpAddr.Port, tcpAddr.Zone
	}

	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		zap.S().Panicw("getLocalAddress",
			"error", "address is not tcp or udp")
	}
	return udpAddr.IP.String(), udpAddr.Port, udpAddr.Zone
}
