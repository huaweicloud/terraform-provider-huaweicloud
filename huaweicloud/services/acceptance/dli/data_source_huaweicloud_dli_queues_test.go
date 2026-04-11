package dli

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceQueues_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dli_queues.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dli_queues.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_dli_queues.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byTags   = "data.huaweicloud_dli_queues.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)

		byPriv   = "data.huaweicloud_dli_queues.filter_by_privilege"
		dcByPriv = acceptance.InitDataSourceCheck(byPriv)

		byCharge   = "data.huaweicloud_dli_queues.filter_by_charge_info"
		dcByCharge = acceptance.InitDataSourceCheck(byCharge)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliElasticResourcePoolName(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceQueues_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Query all queues in elastic resource pool
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "queues.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckOutput("is_existed_queues", "true"),

					// Filter by queue name
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byName, "queues.0.name"),
					resource.TestCheckResourceAttrSet(byName, "queues.0.type"),
					resource.TestCheckResourceAttrSet(byName, "queues.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(byName, "queues.0.created_at"),
					resource.TestCheckResourceAttrSet(byName, "queues.0.owner"),
					resource.TestMatchResourceAttr(byName, "queues.0.scaling_policies.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(byName, "queues.0.scaling_policies.0.priority"),
					resource.TestCheckResourceAttrSet(byName, "queues.0.scaling_policies.0.impact_start_time"),
					resource.TestCheckResourceAttrSet(byName, "queues.0.scaling_policies.0.impact_stop_time"),
					resource.TestCheckResourceAttrSet(byName, "queues.0.scaling_policies.0.min_cu"),
					resource.TestCheckResourceAttrSet(byName, "queues.0.scaling_policies.0.max_cu"),
					resource.TestCheckResourceAttrSet(byName,
						"queues.0.scaling_policies.0.inherit_elastic_resource_pool_max_cu"),

					// Filter by queue type (ListQueues API)
					dcByType.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byType, "queues.#"),

					// Filter by tags (ListQueues API)
					dcByTags.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byTags, "queues.#"),

					// Filter by privilege (ListQueues API)
					dcByPriv.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byPriv, "queues.#"),

					// Filter by charge info (ListQueues API)
					dcByCharge.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byCharge, "queues.#"),
				),
			},
		},
	})
}

func testDataSourceQueues_basic(name string) string {
	return fmt.Sprintf(`
locals {
  pool_names = split(",", "%[1]s")
}

// The CU of the relastic resource pool must be equal to or greater than 16
resource "huaweicloud_dli_queue" "test" {
  elastic_resource_pool_name = local.pool_names[0]
  resource_mode              = 1

  name       = "%[2]s"
  cu_count   = 8
  queue_type = "general"
}

// Query all queues in elastic resource pool
data "huaweicloud_dli_queues" "all" {
  elastic_resource_pool_name = local.pool_names[0]

  depends_on = [
    huaweicloud_dli_queue.test
  ]
}

locals {
  queue_names = data.huaweicloud_dli_queues.all.queues[*].name
}

output "is_existed_queues" {
  value = contains(local.queue_names, "%[2]s")
}

// Filter by queue name
data "huaweicloud_dli_queues" "filter_by_name" {
  elastic_resource_pool_name = local.pool_names[0]

  queue_name = huaweicloud_dli_queue.test.name

  depends_on = [
    huaweicloud_dli_queue.test
  ]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dli_queues.filter_by_name.queues[*].name : v == huaweicloud_dli_queue.test.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) == 1 && alltrue(local.name_filter_result)
}

// Filter by queue type (ListQueues API)
data "huaweicloud_dli_queues" "filter_by_type" {
  queue_type = "sql"

  depends_on = [
    huaweicloud_dli_queue.test
  ]
}

// Filter by tags (ListQueues API)
data "huaweicloud_dli_queues" "filter_by_tags" {
  tags = "key1=value1"

  depends_on = [
    huaweicloud_dli_queue.test
  ]
}

// Filter by privilege (ListQueues API)
data "huaweicloud_dli_queues" "filter_by_privilege" {
  with_privilege = true

  depends_on = [
    huaweicloud_dli_queue.test
  ]
}

// Filter by charge info (ListQueues API)
data "huaweicloud_dli_queues" "filter_by_charge_info" {
  with_charge_info = true

  depends_on = [
    huaweicloud_dli_queue.test
  ]
}
`, acceptance.HW_DLI_ELASTIC_RESOURCE_POOL_NAMES, name)
}
