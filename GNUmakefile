KEY_ID := 51E964C56F41CCAD

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

release:
	goreleaser release --clean --timeout 2h --verbose --parallelism 4