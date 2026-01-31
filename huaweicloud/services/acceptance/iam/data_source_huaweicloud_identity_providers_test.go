package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataProviders_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identity_providers.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_identity_providers.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_identity_providers.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byStatus   = "data.huaweicloud_identity_providers.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataProviders_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "identity_providers.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
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

func testAccDataProviders_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "test" {
  name     = "%[1]s"
  protocol = "saml"
  status   = true
  sso_type = "virtual_user_sso"
}
`, name)
}

func testAccDataProviders_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# All
data "huaweicloud_identity_providers" "all" {
  # Waiting for the provider to be created
  depends_on = [huaweicloud_identity_provider.test]
}

# Filter by name
locals {
  name = "%[2]s"
}

data "huaweicloud_identity_providers" "filter_by_name" {
  name = local.name

  # Waiting for the provider to be created
  depends_on = [huaweicloud_identity_provider.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_identity_providers.filter_by_name.identity_providers : v.id == local.name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) >= 1 && alltrue(local.name_filter_result)
}

# Filter by type
locals {
  type = "virtual_user_sso"
}

data "huaweicloud_identity_providers" "filter_by_type" {
  sso_type = local.type

  # Waiting for the provider to be created
  depends_on = [huaweicloud_identity_provider.test]
}

locals {
  type_filter_result = [for v in data.huaweicloud_identity_providers.filter_by_type.identity_providers : v.sso_type == local.type]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) >= 1 && alltrue(local.type_filter_result)
}

# Filter by status
locals {
  status = true
}

data "huaweicloud_identity_providers" "filter_by_status" {
  status = local.status

  # Waiting for the provider to be created
  depends_on = [huaweicloud_identity_provider.test]
}

locals {
  status_filter_result = [for v in data.huaweicloud_identity_providers.filter_by_status.identity_providers : v.status == local.status]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) >= 1 && alltrue(local.status_filter_result)
}
`, testAccDataProviders_base(name), name)
}
