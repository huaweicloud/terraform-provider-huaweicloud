package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityCompliance_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_security_compliance.test"

		dc = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSecurityCompliance,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "password_regex", passwordRegex),
					resource.TestCheckResourceAttr(all, "password_regex_description", passwordRegexDescription),
				),
			},
			{
				Config: testAccDataSecurityComplianceByOption1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "password_regex", passwordRegex),
				),
			},
			{
				Config: testAccDataSecurityComplianceByOption2,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "password_regex_description", passwordRegexDescription),
				),
			},
		},
	})
}

const testAccDataSecurityCompliance = `
data "huaweicloud_identity_security_compliance" "test" {}
`

const testAccDataSecurityComplianceByOption1 = `
data "huaweicloud_identity_security_compliance" "test" {
  option = "password_regex"
}
`

const testAccDataSecurityComplianceByOption2 = `
data "huaweicloud_identity_security_compliance" "test" {
  option = "password_regex_description"
}
`

const passwordRegex = "^(?![A-Z]*$)(?![a-z]*$)(?![\\d]*$)(?![^\\W]*$)\\S{8,32}$"
const passwordRegexDescription = "The password must contain at least two of the following character types: uppercase " +
	"letters, lowercase letters, digits, and special characters, and be a length between 8 and 32."
