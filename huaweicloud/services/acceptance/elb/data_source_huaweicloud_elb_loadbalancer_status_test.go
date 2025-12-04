package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceElbLoadBalancerStatus_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_loadbalancer_status.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbLoadBalancerStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.#"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.#"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.listeners.#"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.listeners.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.listeners.0.name"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.listeners.0.pools.#"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.listeners.0.pools.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.listeners.0.pools.0.name"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.healthmonitor.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.healthmonitor.0.id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.healthmonitor.0.type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.healthmonitor.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.members.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.members.0.id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.members.0.address"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.members.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.members.0.protocol_port"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.members.0.operating_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.pools.0.operating_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.l7policies.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.l7policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.l7policies.0.action"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.l7policies.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.l7policies.0.rules.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.l7policies.0.rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.l7policies.0.rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.l7policies.0.rules.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.listeners.0.operating_status"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.pools.#"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.pools.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.pools.0.name"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.pools.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.pools.0.healthmonitor.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.pools.0.healthmonitor.0.id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.pools.0.healthmonitor.0.type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.pools.0.healthmonitor.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.pools.0.members.#"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.pools.0.members.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.pools.0.members.0.address"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.pools.0.members.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.pools.0.members.0.protocol_port"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.pools.0.members.0.operating_status"),
					resource.TestCheckResourceAttrSet(dataSource,
						"statuses.0.loadbalancer.0.pools.0.operating_status"),
					resource.TestCheckResourceAttrSet(dataSource, "statuses.0.loadbalancer.0.operating_status"),
				),
			},
		},
	})
}

func testDataSourceElbLoadBalancerStatus_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_elb_loadbalancer_status" "test" {
  loadbalancer_id = "%[1]s"
}
`, acceptance.HW_ELB_LOADBALANCER_ID)
}
