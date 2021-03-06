#!/usr/bin/env bash
go get golang.org/x/tools/cmd/goimports
go get golang.org/x/lint/golint

go vet ./client/ ./server/
vet_rc=$?
[[ vet_rc -eq 0 ]] && echo "GO VET PASS" || echo "GO VET FAIL"

golint ./client/ ./server/
lint_rc=$?
[[ lint_rc -eq 0 ]] && echo "GOLINT PASS" || echo "GOLINT FAIL"

[[ -z "$(goimports -l ./client/ ./server/)" ]]
imports_rc=$?
[[ imports_rc -eq 0 ]] && echo "GOIMPORTS PASS" || echo "GOIMPORTS FAIL"

rc=$((vet_rc + lint_rc + imports_rc))
echo "====="
[[ rc -eq 0 ]] && echo "PASS" || echo "FAIL"
exit $rc
