package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIamIdentityProvider_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_identity_providers.basic"
	dataSource2 := "data.huaweicloud_identity_providers.filter_by_name"
	dataSource3 := "data.huaweicloud_identity_providers.filter_by_type"
	dataSource4 := "data.huaweicloud_identity_providers.filter_by_status"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceIamIdentityProvider_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceIamIdentityProvider_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_providers" "basic" {
  depends_on = [huaweicloud_identity_provider.test]
}

data "huaweicloud_identity_providers" "filter_by_name" {
  name = "%[2]s"

  depends_on = [huaweicloud_identity_provider.test]
}

data "huaweicloud_identity_providers" "filter_by_type" {
  sso_type = "virtual_user_sso"

  depends_on = [huaweicloud_identity_provider.test]
}

data "huaweicloud_identity_providers" "filter_by_status" {
  status = true

  depends_on = [huaweicloud_identity_provider.test]
}

locals {
  name_filter_result   = [for v in data.huaweicloud_identity_providers.filter_by_name.identity_providers[*].id : v == "%[2]s"]
  type_filter_result   = [for v in data.huaweicloud_identity_providers.filter_by_type.identity_providers[*].sso_type : v == "virtual_user_sso"]
  status_filter_result = [for v in data.huaweicloud_identity_providers.filter_by_status.identity_providers[*].status : v == true]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_identity_providers.basic.identity_providers) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testAccIdentityProvider_base(name), name)
}

func testAccIdentityProvider_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "test" {
  name     = "%s"
  protocol = "saml"
  status   = true
  sso_type = "virtual_user_sso"
}
`, name)
}
