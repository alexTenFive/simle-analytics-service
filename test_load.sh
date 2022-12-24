#!/bin/zsh
go run ./cmd/vegeta -rps=300 -duration=60 && \
cat ./cmd/vegeta/reports/report.bin | vegeta plot > ./cmd/vegeta/reports/plot.html
