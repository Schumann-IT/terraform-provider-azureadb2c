package model

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var (
	KeySetIdValidator = stringvalidator.RegexMatches(
		regexp.MustCompile(`^B2C_1A_[a-zA-Z]+$`), "must be prefixed with B2C_1A_ and must contain only alphanumeric characters",
	)

	KeySetNameValidator = stringvalidator.RegexMatches(
		regexp.MustCompile(`^[a-zA-Z]+$`), "must contain only alphanumeric characters",
	)
)
