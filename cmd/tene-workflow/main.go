// Copyright 2026 Kay Kim (kay@agentkay.it)
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/tene-ai/tene-codex/internal/app"
)

var version = "0.1.0-dev"

func main() { os.Exit(app.Run(os.Args[1:], os.Stdout, os.Stderr, version)) }
