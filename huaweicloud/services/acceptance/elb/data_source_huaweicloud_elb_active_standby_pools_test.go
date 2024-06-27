package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceActiveStandbyPools_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_active_standby_pools.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceActiveStandbyPools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "pools.#"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.id"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.name"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.type"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.protocol"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.any_port_enable"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.members.0.address"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.members.0.protocol_port"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.members.0.role"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.members.1.address"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.members.1.protocol_port"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.members.1.role"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.healthmonitor.0.delay"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.healthmonitor.0.expected_codes"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.healthmonitor.0.http_method"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.healthmonitor.0.max_retries"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.healthmonitor.0.max_retries_down"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.healthmonitor.0.timeout"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.healthmonitor.0.type"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.connection_drain_enabled"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.connection_drain_timeout"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("pool_id_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("healthmonitor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("member_address_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceActiveStandbyPools_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_active_standby_pools" "test" {
  depends_on = [huaweicloud_elb_active_standby_pool.test] 
}

data "huaweicloud_elb_active_standby_pools" "name_filter" {
  depends_on = [huaweicloud_elb_active_standby_pool.test]
  name       = "%[2]s"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_active_standby_pools.name_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_active_standby_pools.name_filter.pools[*].name :v == "%[2]s"]
  )
}

locals {
  pool_id = huaweicloud_elb_active_standby_pool.test.id
}

data "huaweicloud_elb_active_standby_pools" "pool_id_filter" {
  depends_on = [huaweicloud_elb_active_standby_pool.test]

  pool_id = huaweicloud_elb_active_standby_pool.test.id
}

output "pool_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_active_standby_pools.pool_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_active_standby_pools.pool_id_filter.pools[*].id : v == local.pool_id]
  )
}

locals {
  description = huaweicloud_elb_active_standby_pool.test.description
}

data "huaweicloud_elb_active_standby_pools" "description_filter" {
  depends_on = [huaweicloud_elb_active_standby_pool.test]

  description = huaweicloud_elb_active_standby_pool.test.description
}

output "description_filter_is_useful" {
  value = length(data.huaweicloud_elb_active_standby_pools.description_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_active_standby_pools.description_filter.pools[*].description : v == local.description]
  )
}

locals {
  protocol = huaweicloud_elb_active_standby_pool.test.protocol
}

data "huaweicloud_elb_active_standby_pools" "protocol_filter" {
  depends_on = [huaweicloud_elb_active_standby_pool.test]

  protocol = huaweicloud_elb_active_standby_pool.test.protocol
}

output "protocol_filter_is_useful" {
  value = length(data.huaweicloud_elb_active_standby_pools.protocol_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_active_standby_pools.protocol_filter.pools[*].protocol : v == local.protocol]
  )
}

locals {
  type = huaweicloud_elb_active_standby_pool.test.type
}

data "huaweicloud_elb_active_standby_pools" "type_filter" {
  depends_on = [huaweicloud_elb_active_standby_pool.test]

  type = huaweicloud_elb_active_standby_pool.test.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_elb_active_standby_pools.type_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_active_standby_pools.type_filter.pools[*].type : v == local.type]
  )
}

locals {
  healthmonitor_id = huaweicloud_elb_active_standby_pool.test.healthmonitor.0.id
}

data "huaweicloud_elb_active_standby_pools" "healthmonitor_id_filter" {
  depends_on = [huaweicloud_elb_active_standby_pool.test]

  healthmonitor_id = huaweicloud_elb_active_standby_pool.test.healthmonitor.0.id
}

output "healthmonitor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_active_standby_pools.healthmonitor_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_active_standby_pools.healthmonitor_id_filter.pools[*].healthmonitor.0.id :
  v == local.healthmonitor_id]
  )
}

locals {
  vpc_id = huaweicloud_elb_active_standby_pool.test.vpc_id
}

data "huaweicloud_elb_active_standby_pools" "vpc_id_filter" {
  depends_on = [huaweicloud_elb_active_standby_pool.test]

  vpc_id = huaweicloud_elb_active_standby_pool.test.vpc_id
}

output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_active_standby_pools.vpc_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_active_standby_pools.vpc_id_filter.pools[*].vpc_id : v == local.vpc_id]
  )
}

locals {
  member_address = tolist(huaweicloud_elb_active_standby_pool.test.members).0.address
}

data "huaweicloud_elb_active_standby_pools" "member_address_filter" {
  depends_on = [huaweicloud_elb_active_standby_pool.test]

  member_address = tolist(huaweicloud_elb_active_standby_pool.test.members).0.address
}

output "member_address_filter_is_useful" {
  value = length(data.huaweicloud_elb_active_standby_pools.member_address_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_active_standby_pools.member_address_filter.pools[*].members[*].address : contains(v, local.member_address)]
  )
}
`, testAccElbActiveStandbyPoolConfig_basic(name), name)
}
