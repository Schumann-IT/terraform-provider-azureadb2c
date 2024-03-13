default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

check:
	@if [[ "" == "$(GPG_FINGERPRINT)" ]]; then echo "please provide GPG_FINGERPRINT"; exit 1; fi
	@if [[ "" == "$(GITHUB_TOKEN)" ]]; then echo "please provide GITHUB_TOKEN"; exit 1; fi

release: check
	@goreleaser release --clean --timeout 2h --verbose --parallelism 4