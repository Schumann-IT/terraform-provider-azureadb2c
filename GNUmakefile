KEY_ID := 51E964C56F41CCAD
VERSION := 0.1.0

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

create-provider:
	@echo '{"data":{"type":"registry-providers","attributes":{"name":"azureadb2c","namespace":"schumann-it","registry-name":"private"}}}' > provider.json
	@curl -s --header "Authorization: Bearer $(TOKEN)" --header "Content-Type: application/vnd.api+json" --request POST --data @provider.json https://app.terraform.io/api/v2/organizations/schumann-it/registry-providers | jq
	@rm -f provider.json

list-gpg-keys:
	@curl -s --header "Authorization: Bearer $(TOKEN)" --header "Content-Type: application/vnd.api+json" --request GET https://app.terraform.io/api/registry/private/v2/gpg-keys?filter%5Bnamespace%5D=schumann-it | jq

create-gpg-key:
	@curl -s --header "Authorization: Bearer $(TOKEN)" --header "Content-Type: application/vnd.api+json" --request POST --data @key.json https://app.terraform.io/api/registry/private/v2/gpg-keys

create-version:
	@echo '{"data":{"type":"registry-provider-versions","attributes":{"version":"$(VERSION)","key-id":"$(KEY_ID)","protocols":["6.0"]}}}' > version.json
	@curl -s --header "Authorization: Bearer $(TOKEN)" --header "Content-Type: application/vnd.api+json" --request POST --data @version.json https://app.terraform.io/api/v2/organizations/schumann-it/registry-providers/private/schumann-it/azureadb2c/versions > version-response.json
	@rm -f version.json

upload-sigs:
	@curl -s -T dist/terraform-provider-azureadb2c_$(VERSION)_SHA256SUMS $(shell cat version-response.json | jq '.data.links."shasums-upload"')
	@curl -s -T dist/terraform-provider-azureadb2c_$(VERSION)_SHA256SUMS.sig $(shell cat version-response.json | jq '.data.links."shasums-sig-upload"')
	@rm -f version-response.json

create-platforms:
	@echo '{"data":{"type":"registry-provider-version-platforms","attributes":{"os":"darwin","arch":"arm64","shasum":"$(shell cat dist/terraform-provider-azureadb2c_$(VERSION)_SHA256SUMS | grep terraform-provider-azureadb2c_$(VERSION)_darwin_arm64.zip | awk '{print $$1}')","filename":"terraform-provider-azureadb2c_$(VERSION)_darwin_arm64.zip"}}}' > platform.json
	@curl -s --header "Authorization: Bearer $(TOKEN)" --header "Content-Type: application/vnd.api+json" --request POST --data @platform.json https://app.terraform.io/api/v2/organizations/schumann-it/registry-providers/private/schumann-it/azureadb2c/versions/$(VERSION)/platforms > platform-response.json
	@rm -f platform.json

upload-binary:
	@curl -s -T dist/terraform-provider-azureadb2c_$(VERSION)_darwin_arm64.zip $(shell cat platform-response.json | jq '.data.links."provider-binary-upload"')
	@rm -f platform-response.json