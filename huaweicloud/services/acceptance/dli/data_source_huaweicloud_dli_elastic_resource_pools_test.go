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
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_dli_elastic_resource_pools.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNameNotFound   = "data.huaweicloud_dli_elastic_resource_pools.filter_by_name_not_found"
		dcByNameNotFound = acceptance.InitDataSourceCheck(byNameNotFound)

		byStatus   = "data.huaweicloud_dli_elastic_resource_pools.filter_by_name_not_found"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byTags   = "data.huaweicloud_dli_elastic_resource_pools.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)

		_, elasticResourceName = getElasticResourcePoolNames()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliElasticResourcePoolName(t)
			acceptance.TestAccPreCheckDliSQLQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDliElasticPools_basic(elasticResourceName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.id"),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.name", elasticResourceName),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.max_cu"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.min_cu"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.current_cu"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.cidr"),
					resource.TestCheckResourceAttr(byName, "elastic_resource_pools.0.enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.queues.0"),
					resource.TestCheckResourceAttrSet(byName, "elastic_resource_pools.0.description"),
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
				),
			},
		},
	})
}

func testDataSourceDataSourceDliElasticPools_basic(elasticResourceName string) string {
	return fmt.Sprintf(`
locals {
  tags = {
    foo       = "bar"
    terraform = "elastic_resource_pool"
  }
}
data "huaweicloud_dli_elastic_resource_pools" "test" {}

data "huaweicloud_dli_elastic_resource_pools" "filter_by_name" {
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
  name = "not_found"
}

output "is_name_filter_useful_not_found" {
  value = length(data.huaweicloud_dli_elastic_resource_pools.filter_by_name_not_found.elastic_resource_pools) == 0
}

locals {
  status = data.huaweicloud_dli_elastic_resource_pools.filter_by_name.elastic_resource_pools[0].status
}

data "huaweicloud_dli_elastic_resource_pools" "filter_by_status" {
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
  tags = local.tags
}

output "is_tags_filter_useful" {
  value = length(data.huaweicloud_dli_elastic_resource_pools.filter_by_tags) >= 1
}
`, elasticResourceName)
}
