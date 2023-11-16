package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEnterpriseProjectsDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_enterprise_projects.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProjectsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "enterprise_projects.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enterprise_projects.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enterprise_projects.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enterprise_projects.0.type"),

					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccEnterpriseProjectsDataSource_basic() string {
	return `
data "huaweicloud_enterprise_projects" "test" {}

locals {
  enterprise_project_id = data.huaweicloud_enterprise_projects.test.enterprise_projects[0].id
}
data "huaweicloud_enterprise_projects" "enterprise_project_id_filter" {
  enterprise_project_id = local.enterprise_project_id
}
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_enterprise_projects.enterprise_project_id_filter.enterprise_projects) > 0 && alltrue(
	[for v in data.huaweicloud_enterprise_projects.enterprise_project_id_filter.enterprise_projects[*].id : v == local.enterprise_project_id]
  )  
}

locals {
  name = data.huaweicloud_enterprise_projects.test.enterprise_projects[0].name
}
data "huaweicloud_enterprise_projects" "name_filter" {
  name = local.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_enterprise_projects.name_filter.enterprise_projects) > 0 && alltrue(
	[for v in data.huaweicloud_enterprise_projects.name_filter.enterprise_projects[*].name : v == local.name]
  )  
}

locals {
  status = data.huaweicloud_enterprise_projects.test.enterprise_projects[0].status
}
data "huaweicloud_enterprise_projects" "status_filter" {
  status = local.status
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_enterprise_projects.status_filter.enterprise_projects) > 0 && alltrue(
	[for v in data.huaweicloud_enterprise_projects.status_filter.enterprise_projects[*].status : v == local.status]
  )  
}

locals {
  type = data.huaweicloud_enterprise_projects.test.enterprise_projects[0].type
}
data "huaweicloud_enterprise_projects" "type_filter" {
  type = local.type
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_enterprise_projects.type_filter.enterprise_projects) > 0 && alltrue(
	[for v in data.huaweicloud_enterprise_projects.type_filter.enterprise_projects[*].type : v == local.type]
  )  
}`
}
