package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCloudConnections_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_connections.filter_by_id"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)
	baseConfig := testAccDatasourceCreateCloudConnections(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcCloudConnections_basic(baseConfig, rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_connections.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_connections.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_connections.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_connections.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_connections.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_connections.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_connections.0.updated_at"),
					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_not_found", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("name_and_tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcCloudConnections_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id   = huaweicloud_cc_connection.test1.id
  name = huaweicloud_cc_connection.test1.name
  tags = huaweicloud_cc_connection.test2.tags
}	

data "huaweicloud_cc_connections" "filter_by_id" {
  cloud_connection_id = local.id
}

data "huaweicloud_cc_connections" "filter_by_name" {
  name = "%[2]s1"

  depends_on = [
    huaweicloud_cc_connection.test1,
    huaweicloud_cc_connection.test2,
    huaweicloud_cc_connection.test3,
  ]
}

data "huaweicloud_cc_connections" "filter_by_name_not_found" {
  name = "%[2]s_not_found"
  
  depends_on = [
    huaweicloud_cc_connection.test1,
    huaweicloud_cc_connection.test2,
    huaweicloud_cc_connection.test3,
  ]
}

data "huaweicloud_cc_connections" "filter_by_tags" {
  tags = local.tags

  depends_on = [
    huaweicloud_cc_connection.test1,
    huaweicloud_cc_connection.test2,
    huaweicloud_cc_connection.test3,
  ]
}

data "huaweicloud_cc_connections" "filter_by_name_and_tags" {
  name = local.name
  tags = local.tags

  depends_on = [
    huaweicloud_cc_connection.test1,
    huaweicloud_cc_connection.test2,
    huaweicloud_cc_connection.test3,
  ]
}

output "id_filter_is_useful" {
  value = length(data.huaweicloud_cc_connections.filter_by_id.cloud_connections) > 0 && alltrue(
    [for v in data.huaweicloud_cc_connections.filter_by_id.cloud_connections[*].id : v == local.id]
  )
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_cc_connections.filter_by_name.cloud_connections) > 0 && alltrue(
	[for v in data.huaweicloud_cc_connections.filter_by_name.cloud_connections[*].name : v == local.name]
  )
}

output "name_filter_not_found" {
  value = length(data.huaweicloud_cc_connections.filter_by_name_not_found.cloud_connections) == 0
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_cc_connections.filter_by_tags.cloud_connections) >= 1 && alltrue([
    for cc in data.huaweicloud_cc_connections.filter_by_tags.cloud_connections : alltrue([
      for k, v in local.tags : cc.tags[k] == v
    ])
  ])
}

output "name_and_tags_filter_is_useful" {
  value = length(data.huaweicloud_cc_connections.filter_by_name_and_tags.cloud_connections) >= 1 && alltrue([
    for cc in data.huaweicloud_cc_connections.filter_by_name_and_tags.cloud_connections : alltrue([
      for k, v in local.tags : cc.tags[k] == v
    ]) && cc.name == local.name
  ])
}
`, baseConfig, name)
}

func testAccDatasourceCreateCloudConnections(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test1" {
  name                  = "%[1]s1"
  description           = "demo 1"
	
  tags = {
    key = "value"
  }
}

resource "huaweicloud_cc_connection" "test2" {
  name                  = "%[1]s2"
  description           = "demo 2"
	  
  tags = {
    key = "value"
  }
}

resource "huaweicloud_cc_connection" "test3" {
  name                  = "%[1]s3"
  description           = "demo 3"
	  
  tags = {
    foo = "bar"
  }
}
`, rName)
}
