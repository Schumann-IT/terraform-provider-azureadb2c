package acceptance

import "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

func RandAlphanumericStrings(count, strLen int) []string {
	var res []string
	for i := 1; i <= count; i++ {
		res = append(res, acctest.RandStringFromCharSet(strLen, "ABCDEFGHIJKLMNOPQRSTXYZabcdefghijklmnopqrstxyz"))
	}

	return res
}
