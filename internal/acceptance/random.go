package acceptance

import "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

func RandAlphanumericString(count, len int) []string {
	var res []string
	for i := 1; i <= count; i++ {
		res = append(res, acctest.RandStringFromCharSet(len, "ABCDEFGHIJKLMNOPQRSTXYZabcdefghijklmnopqrstxyz"))
	}

	return res
}
