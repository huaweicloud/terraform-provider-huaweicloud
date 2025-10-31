package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentitySecurityCompliance_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_security_compliance.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTestDataSourceIdentitySecurityCompliance,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "password_regex", passwordRegex),
					resource.TestCheckResourceAttr(dataSourceName, "password_regex_description", passwordRegexDescription),
				),
			},
			{
				Config: testTestDataSourceIdentitySecurityComplianceByOption1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "password_regex", passwordRegex),
				),
			},
			{
				Config: testTestDataSourceIdentitySecurityComplianceByOption2,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "password_regex_description", passwordRegexDescription),
				),
			},
		},
	})
}

const testTestDataSourceIdentitySecurityCompliance = `
data "huaweicloud_identity_security_compliance" "test" {}
`

const testTestDataSourceIdentitySecurityComplianceByOption1 = `
data "huaweicloud_identity_security_compliance" "test" {
  option = "password_regex"
}
`

const testTestDataSourceIdentitySecurityComplianceByOption2 = `
data "huaweicloud_identity_security_compliance" "test" {
  option = "password_regex_description"
}
`

const passwordRegex = "^(?![A-Z]*$)(?![a-z]*$)(?![\\d]*$)(?![^\\W]*$)\\S{8,32}$"
const passwordRegexDescription = "The password must contain at least two of the following character types: uppercase " +
	"letters, lowercase letters, digits, and special characters, and be a length between 8 and 32."
