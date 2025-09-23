package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstances_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		rName      = acceptance.RandomAccResourceName()

		byId   = "data.huaweicloud_apig_instances.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_instances.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_apig_instances.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEpsId   = "data.huaweicloud_apig_instances.filter_by_enterprise_project_id"
		dcByEpsId = acceptance.InitDataSourceCheck(byEpsId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.edition"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_id"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_instances" "test" {
  depends_on = [
    huaweicloud_apig_instance.test
  ]
}

# Filter by ID
locals {
  instance_id = huaweicloud_apig_instance.test.id
}

data "huaweicloud_apig_instances" "filter_by_id" {
  instance_id = local.instance_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_instances.filter_by_id.instances[*].id : v == local.instance_id
  ]
}

output "instance_id_filter_is_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  name = huaweicloud_apig_instance.test.name
}

data "huaweicloud_apig_instances" "filter_by_name" {
  depends_on = [
    huaweicloud_apig_instance.test
  ]

  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_instances.filter_by_name.instances[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by status
locals {
  status = huaweicloud_apig_instance.test.status
}

data "huaweicloud_apig_instances" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_apig_instances.filter_by_status.instances[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by enterprise_project_id
locals {
  enterprise_project_id = huaweicloud_apig_instance.test.enterprise_project_id
}

data "huaweicloud_apig_instances" "filter_by_enterprise_project_id" {
  depends_on = [
    huaweicloud_apig_instance.test
  ]

  enterprise_project_id = local.enterprise_project_id
}

locals {
  enterprise_project_id_filter_result = [
    for v in data.huaweicloud_apig_instances.filter_by_enterprise_project_id.instances[*].enterprise_project_id : v == local.enterprise_project_id
  ]
}

output "enterprise_project_id_filter_is_useful" {
  value = length(local.enterprise_project_id_filter_result) > 0 && alltrue(local.enterprise_project_id_filter_result)
}
`, testAccInstance_basic_step1(name))
}
