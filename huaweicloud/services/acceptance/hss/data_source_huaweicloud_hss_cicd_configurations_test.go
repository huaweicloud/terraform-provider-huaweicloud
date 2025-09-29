package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCiCdConfigurations_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_cicd_configurations.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCiCdConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cicd_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cicd_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.associated_images_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.associated_config_num"),

					resource.TestCheckOutput("is_cicd_id_filter_useful", "true"),
					resource.TestCheckOutput("is_cicd_name_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceCiCdConfigurations_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_hss_cicd_configuration" "test" {
  cicd_name               = "%[1]s"
  vulnerability_whitelist = ["vulnerability_whitelist_test"]
  vulnerability_blocklist = ["vulnerability_blocklist_test"]
  image_whitelist         = ["image_whitelist_test"]
  enterprise_project_id   = "0"
}
`, name)
}

func testDataSourceCiCdConfigurations_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_cicd_configurations" "test" {
  depends_on = [huaweicloud_hss_cicd_configuration.test]
}

# Filter using cicd_id.
locals {
  cicd_id = data.huaweicloud_hss_cicd_configurations.test.data_list[0].cicd_id
}

data "huaweicloud_hss_cicd_configurations" "cicd_id_filter" {
  cicd_id = local.cicd_id
}

output "is_cicd_id_filter_useful" {
  value = length(data.huaweicloud_hss_cicd_configurations.cicd_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_cicd_configurations.cicd_id_filter.data_list[*].cicd_id : v == local.cicd_id]
  )
}

# Filter using cicd_name.
locals {
  cicd_name = data.huaweicloud_hss_cicd_configurations.test.data_list[0].cicd_name
}

data "huaweicloud_hss_cicd_configurations" "cicd_name_filter" {
  cicd_name = local.cicd_name
}

output "is_cicd_name_filter_useful" {
  value = length(data.huaweicloud_hss_cicd_configurations.cicd_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_cicd_configurations.cicd_name_filter.data_list[*].cicd_name : v == local.cicd_name]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_cicd_configurations" "enterprise_project_id_filter" {
  depends_on = [huaweicloud_hss_cicd_configuration.test]

  enterprise_project_id = "0"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_cicd_configurations.enterprise_project_id_filter.data_list) > 0
}

# Filter using non existent cicd_name.
data "huaweicloud_hss_cicd_configurations" "not_found" {
  cicd_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_cicd_configurations.not_found.data_list) == 0
}
`, testDataSourceCiCdConfigurations_base())
}
