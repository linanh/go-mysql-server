// Copyright 2020-2021 Dolthub, Inc.
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

// +build darwin

package sockstate

import (
	"net"

	"github.com/sirupsen/logrus"
)

// tcpSocks returns a slice of active TCP sockets containing only those
// elements that satisfy the accept function
func tcpSocks(accept AcceptFn) ([]sockTabEntry, error) {
	// (juanjux) TODO: not implemented
	logrus.Warn("mysql/server connection checking not implemented for Darwin")
	return nil, ErrSocketCheckNotImplemented.New()
}

func GetConnInode(c *net.TCPConn) (n uint64, err error) {
	return 0, ErrSocketCheckNotImplemented.New()
}
