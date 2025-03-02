#!/bin/bash
set -eux

go mod tidy
go install -v golang.org/x/tools/gopls@latest