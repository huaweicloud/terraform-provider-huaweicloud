package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDliElasticResourcePools_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dli_elastic_resource_pools.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_dli_elastic_resource_pools.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNameNotFound   = "data.huaweicloud_dli_elastic_resource_pools.filter_by_name_not_found"
		dcByNameNotFound = acceptance.InitDataSourceCheck(byNameNotFound)

		byStatus   = "data.huaweicloud_dli_elastic_resource_pools.filter_by_name_not_found"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byTags   = "data.huaweicloud_dli_elastic_resource_pools.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDliElasticPools_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.id"),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.name", rName),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.max_cu", "64"),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.min_cu", "64"),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.current_cu", "64"),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.cidr", "172.16.0.0/12"),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.queues.0", rName),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.description", "Created by terraform script"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.resource_id"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.owner"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.manager"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.created_at"),

					dcByNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful_not_found", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful_not_found", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),

					waitForDeletionCooldownComplete(),
				),
			},
		},
	})
}

func testDataSourceDataSourceDliElasticPools_basic(name string) string {
	return fmt.Sprintf(`
locals {
  tags = {
    foo       = "bar"
    terraform = "elastic_resource_pool"
  }
}

resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = "%[1]s"
  max_cu                = 64
  min_cu                = 64
  enterprise_project_id = "0"
  description           = "Created by terraform script"
  tags                  = local.tags
}

resource "huaweicloud_dli_queue" "test" {
  depends_on = [
    huaweicloud_dli_elastic_resource_pool.test,
  ] 
  
  name                       = "%[1]s"
  cu_count                   = 16
  resource_mode              = 1
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name
}

data "huaweicloud_dli_elastic_resource_pools" "test" {
  depends_on = [
    huaweicloud_dli_elastic_resource_pool.test,
    huaweicloud_dli_queue.test
  ]
}
  
data "huaweicloud_dli_elastic_resource_pools" "filter_by_name" {
  depends_on = [
    huaweicloud_dli_elastic_resource_pool.test,
    huaweicloud_dli_queue.test
  ]

  name = "%[1]s"
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dli_elastic_resource_pools.filter_by_name.elastic_resource_pools[*].name : v == "%[1]s"
  ]
}
    
output "is_name_filter_useful" {
  value = length(local.name_filter_result) == 1 && alltrue(local.name_filter_result)
}
    
data "huaweicloud_dli_elastic_resource_pools" "filter_by_name_not_found" {
  depends_on = [
    huaweicloud_dli_elastic_resource_pool.test,
    huaweicloud_dli_queue.test
  ]

  name = "not_found"
}
    
output "is_name_filter_useful_not_found" {
  value = length(data.huaweicloud_dli_elastic_resource_pools.filter_by_name_not_found.elastic_resource_pools) == 0
}

locals {
  status = huaweicloud_dli_elastic_resource_pool.test.status
}
    
data "huaweicloud_dli_elastic_resource_pools" "filter_by_status" {
  depends_on = [
    huaweicloud_dli_elastic_resource_pool.test,
    huaweicloud_dli_queue.test
  ]
 
   status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dli_elastic_resource_pools.filter_by_status.elastic_resource_pools[*].status : v == local.status
  ]
}
    
output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

data "huaweicloud_dli_elastic_resource_pools" "filter_by_tags" {
  depends_on = [
    huaweicloud_dli_elastic_resource_pool.test,
    huaweicloud_dli_queue.test
  ]

  tags = local.tags
}

output "is_tags_filter_useful" {
  value = length(data.huaweicloud_dli_elastic_resource_pools.filter_by_tags) >= 1
}  
`, name)
}
