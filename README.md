# Terraform Provider for Azure B2C

## Release

```bash
VERSION=<the version> TOKEN=<the terraform cloud api token> make create-version
GPG_FINGERPRINT=<the fingerprint> VERSION=<the version> make release-version
VERSION=<the version> TOKEN=<the terraform cloud api token> make upload-sigs 
VERSION=<the version> TOKEN=<the terraform cloud api token> make create-platforms
VERSION=<the version> TOKEN=<the terraform cloud api token> make upload-binary  
```