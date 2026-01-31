package lakeformation

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataInstances_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_lakeformation_instances.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byInRecycleBin   = "data.huaweicloud_lakeformation_instances.filter_by_in_recycle_bin"
		dcByInRecycleBin = acceptance.InitDataSourceCheck(byInRecycleBin)

		byAllGrantedEps   = "data.huaweicloud_lakeformation_instances.filter_by_all_granted_eps"
		dcByAllGrantedEps = acceptance.InitDataSourceCheck(byAllGrantedEps)

		byEpsId   = "data.huaweicloud_lakeformation_instances.filter_by_enterprise_project_id"
		dcByEpsId = acceptance.InitDataSourceCheck(byEpsId)

		byName   = "data.huaweicloud_lakeformation_instances.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstances_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByInRecycleBin.CheckResourceExists(),
					resource.TestCheckOutput("is_in_recycle_bin_filter_useful", "true"),
					dcByAllGrantedEps.CheckResourceExists(),
					resource.TestCheckOutput("is_all_granted_eps_filter_useful", "true"),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.instance_id",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.instance_id"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.name",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.name"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.description",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.description"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.enterprise_project_id",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.shared",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.shared"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.default_instance",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.default_instance"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.create_time",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.create_time"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.update_time",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.update_time"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.status",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.status"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.in_recycle_bin",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.in_recycle_bin"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.tags.%",
						"data.huaweicloud_lakeformation_instances.all", "instances.0.tags.%"),
				),
			},
		},
	})
}

const testAccDataSourceInstances_basic = `
data "huaweicloud_lakeformation_instances" "all" {}

# Filter by in_recycle_bin
locals {
  in_recycle_bin = data.huaweicloud_lakeformation_instances.all.instances[0].in_recycle_bin
}

data "huaweicloud_lakeformation_instances" "filter_by_in_recycle_bin" {
  in_recycle_bin = local.in_recycle_bin
}

locals {
  in_recycle_bin_filter_result = [
    for v in data.huaweicloud_lakeformation_instances.filter_by_in_recycle_bin.instances[*].in_recycle_bin :
      v == local.in_recycle_bin
  ]
}

output "is_in_recycle_bin_filter_useful" {
  value = length(local.in_recycle_bin_filter_result) > 0 && alltrue(local.in_recycle_bin_filter_result)
}

# Query all instances with all_granted_eps enterprise project ID
data "huaweicloud_lakeformation_instances" "filter_by_all_granted_eps" {
  enterprise_project_id = "all_granted_eps"
}

output "is_all_granted_eps_filter_useful" {
  value = length(data.huaweicloud_lakeformation_instances.filter_by_all_granted_eps.instances) > 0
}

# Filter by a specified enterprise project ID
locals {
  enterprise_project_id = data.huaweicloud_lakeformation_instances.all.instances[0].enterprise_project_id
}

data "huaweicloud_lakeformation_instances" "filter_by_enterprise_project_id" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  enterprise_project_id_filter_result = [
    for v in data.huaweicloud_lakeformation_instances.filter_by_enterprise_project_id.instances[*].enterprise_project_id :
      v == local.enterprise_project_id
  ]
}

output "is_enterprise_project_id_filter_useful" {
  value = length(local.enterprise_project_id_filter_result) > 0 && alltrue(local.enterprise_project_id_filter_result)
}

# Filter by name
locals {
  instance_name = try(data.huaweicloud_lakeformation_instances.all.instances[0].name, "NOT_FOUND")
}

data "huaweicloud_lakeformation_instances" "filter_by_name" {
  name = local.instance_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_lakeformation_instances.filter_by_name.instances[*].name :
      v == local.instance_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`
