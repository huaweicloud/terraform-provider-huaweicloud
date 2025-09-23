package ga

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceAccelerators_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_ga_accelerators.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_ga_accelerators.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byAcceleratorId   = "data.huaweicloud_ga_accelerators.filter_by_accelerator_id"
		dcByAcceleratorId = acceptance.InitDataSourceCheck(byAcceleratorId)

		byStatus   = "data.huaweicloud_ga_accelerators.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEps   = "data.huaweicloud_ga_accelerators.filter_by_eps"
		dcByEps = acceptance.InitDataSourceCheck(byEps)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccelerators_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.0.ip_sets.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.0.tags.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "accelerators.0.updated_at"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByAcceleratorId.CheckResourceExists(),
					resource.TestCheckOutput("accelerator_id_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByEps.CheckResourceExists(),
					resource.TestCheckOutput("eps_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAccelerators_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_accelerator" "test" {
  name        = "%s"
  description = "terraform test"

  ip_sets {
    area = "CM"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testAccDataSourceAccelerators_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ga_accelerators" "test" {
  depends_on = [
    huaweicloud_ga_accelerator.test
  ]  
}

# Filter by name
locals {
  name = data.huaweicloud_ga_accelerators.test.accelerators[0].name
}

data "huaweicloud_ga_accelerators" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_ga_accelerators.filter_by_name.accelerators[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by accelerator_id
locals {
  accelerator_id = data.huaweicloud_ga_accelerators.test.accelerators[0].id
}

data "huaweicloud_ga_accelerators" "filter_by_accelerator_id" {
  accelerator_id = local.accelerator_id
}

locals {
  accelerator_id_filter_result = [
    for v in data.huaweicloud_ga_accelerators.filter_by_accelerator_id.accelerators[*].id : v == local.accelerator_id
  ]
}

output "accelerator_id_filter_is_useful" {
  value = alltrue(local.accelerator_id_filter_result) && length(local.accelerator_id_filter_result) > 0
}

# Filter by status
locals {
  status = data.huaweicloud_ga_accelerators.test.accelerators[0].status
}

data "huaweicloud_ga_accelerators" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ga_accelerators.filter_by_status.accelerators[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

# Filter by enterprise_project_id
locals {
  enterprise_project_id = data.huaweicloud_ga_accelerators.test.accelerators[0].enterprise_project_id
}

data "huaweicloud_ga_accelerators" "filter_by_eps" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  eps_filter_result = [
    for v in data.huaweicloud_ga_accelerators.filter_by_eps.accelerators[*].enterprise_project_id : 
    v == local.enterprise_project_id
  ]
}

output "eps_filter_is_useful" {
  value = alltrue(local.eps_filter_result) && length(local.eps_filter_result) > 0
}
`, testAccDataSourceAccelerators_base(name))
}
