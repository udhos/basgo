#!/bin/bash

basgo-build < gofunc.bas > a.go && go run a.go gofunc.go
