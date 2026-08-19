// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"strings"
	"time"
)

var crockford = base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").WithPadding(base32.NoPadding)

func NewID(prefix string) string {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:8], uint64(time.Now().UTC().UnixMilli()))
	if _, err := rand.Read(buf[8:]); err != nil {
		binary.BigEndian.PutUint64(buf[8:], uint64(time.Now().UnixNano()))
	}
	return prefix + "_" + strings.ToLower(crockford.EncodeToString(buf))
}
