package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSecurityCompliance_basic(t *testing.T) {
	var (
		byAll   = "data.huaweicloud_identity_security_compliance.all"
		dcByAll = acceptance.InitDataSourceCheck(byAll)

		byPasswordRegex   = "data.huaweicloud_identity_security_compliance.filter_by_password_regex"
		dcByPasswordRegex = acceptance.InitDataSourceCheck(byPasswordRegex)

		byPasswordRegexDescription   = "data.huaweicloud_identity_security_compliance.filter_by_password_regex_description"
		dcByPasswordRegexDescription = acceptance.InitDataSourceCheck(byPasswordRegexDescription)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentitySecurityCompliance_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcByAll.CheckResourceExists(),
					resource.TestCheckResourceAttr(byAll, "password_regex", passwordRegex),
					resource.TestCheckResourceAttr(byAll, "password_regex_description", passwordRegexDescription),

					// filter by password_regex
					dcByPasswordRegex.CheckResourceExists(),
					resource.TestCheckOutput("is_password_regex_filter_useful", "true"),

					// filter by password_regex_description
					dcByPasswordRegexDescription.CheckResourceExists(),
					resource.TestCheckOutput("is_password_regex_description_filter_useful", "true"),
				),
			},
		},
	})
}

const passwordRegex = "^(?![A-Z]*$)(?![a-z]*$)(?![\\d]*$)(?![^\\W]*$)\\S{8,32}$"

const passwordRegexDescription = "The password must contain at least two of the following character types: uppercase " +
	"letters, lowercase letters, digits, and special characters, and be a length between 8 and 32."

func testAccIdentitySecurityCompliance_basic() string {
	return fmt.Sprintf(`
# All
data "huaweicloud_identity_security_compliance" "all" {}

# Filter by option password_regex
data "huaweicloud_identity_security_compliance" "filter_by_password_regex" {
  option = "password_regex"
}

locals {
  password_regex = "%[1]s"
}

output "is_password_regex_filter_useful" {
  value = (
    data.huaweicloud_identity_security_compliance.filter_by_password_regex.password_regex == local.password_regex
  )
}

# Filter by option password_regex_description
data "huaweicloud_identity_security_compliance" "filter_by_password_regex_description" {
  option = "password_regex_description"
}

locals {
  password_regex_description = "%[2]s"
}

output "is_password_regex_description_filter_useful" {
  value = (
    data.huaweicloud_identity_security_compliance.filter_by_password_regex_description.password_regex_description
    == local.password_regex_description
  )
}
`, strings.ReplaceAll(passwordRegex, `\`, `\\`), passwordRegexDescription)
}
