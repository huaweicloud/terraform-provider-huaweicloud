package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceFunctionGraphApplication_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_fgs_applications.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceFunctionGraphApplications_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "applications.0.id"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.name"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.status"),

					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceFunctionGraphApplications_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_fgs_applications" "test" {
  depends_on = [
    huaweicloud_fgs_application.test
  ]
}

data "huaweicloud_fgs_applications" "application_id_filter" {
  application_id = local.application_id
}
  
locals {
  application_id = data.huaweicloud_fgs_applications.test.applications[0].id
}
  
output "application_id_filter_is_useful" {
  value = length(data.huaweicloud_fgs_applications.application_id_filter.applications) > 0 && alltrue(
    [for v in data.huaweicloud_fgs_applications.application_id_filter.applications[*].id : v == local.application_id]
  )
}

data "huaweicloud_fgs_applications" "name_filter" {
  name = local.name
}
  
locals {
  name = data.huaweicloud_fgs_applications.test.applications[0].name
}
  
output "name_filter_is_useful" {
  value = length(data.huaweicloud_fgs_applications.name_filter.applications) > 0 && alltrue(
    [for v in data.huaweicloud_fgs_applications.name_filter.applications[*].name : v == local.name]
  )
}

data "huaweicloud_fgs_applications" "status_filter" {
  status = local.status
}

locals {
  status = data.huaweicloud_fgs_applications.test.applications[0].status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_fgs_applications.status_filter.applications) > 0 && alltrue(
    [for v in data.huaweicloud_fgs_applications.status_filter.applications[*].status : v == local.status]
  )
}
`, testAccApplication_basic(name))
}
