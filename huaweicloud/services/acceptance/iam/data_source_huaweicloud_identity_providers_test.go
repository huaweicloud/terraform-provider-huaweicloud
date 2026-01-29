package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceProviders_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_identity_providers.all"
		dcAll  = acceptance.InitDataSourceCheck(dcName)

		dcNameByName = "data.huaweicloud_identity_providers.filter_by_name"
		dcByName     = acceptance.InitDataSourceCheck(dcNameByName)

		dcNameByType = "data.huaweicloud_identity_providers.filter_by_type"
		dcByType     = acceptance.InitDataSourceCheck(dcNameByType)

		dcNameByStatus = "data.huaweicloud_identity_providers.filter_by_status"
		dcByStatus     = acceptance.InitDataSourceCheck(dcNameByStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProviders_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dcAll.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "identity_providers.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccProviders_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "test" {
  name     = "%[1]s"
  protocol = "saml"
  status   = true
  sso_type = "virtual_user_sso"
}
`, name)
}

func testAccDataSourceProviders_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# All
data "huaweicloud_identity_providers" "all" {
  depends_on = [huaweicloud_identity_provider.test]
}

# Filter by name
data "huaweicloud_identity_providers" "filter_by_name" {
  name = "%[2]s"

  depends_on = [huaweicloud_identity_provider.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_identity_providers.filter_by_name.identity_providers : v.id == "%[2]s"]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) >= 1 && alltrue(local.name_filter_result)
}

# Filter by type
data "huaweicloud_identity_providers" "filter_by_type" {
  sso_type = "virtual_user_sso"

  depends_on = [huaweicloud_identity_provider.test]
}

locals {
  type_filter_result = [for v in data.huaweicloud_identity_providers.filter_by_type.identity_providers : v.sso_type == "virtual_user_sso"]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) >= 1 && alltrue(local.type_filter_result)
}

# Filter by status
data "huaweicloud_identity_providers" "filter_by_status" {
  status = true

  depends_on = [huaweicloud_identity_provider.test]
}

locals {
  status_filter_result = [for v in data.huaweicloud_identity_providers.filter_by_status.identity_providers : v.status == true]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) >= 1 && alltrue(local.status_filter_result)
}
`, testAccProviders_base(name), name)
}
