package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "max_htap_instance_num_of_taurus"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.project_id"),
					resource.TestCheckResourceAttr(dataSource, "instances.0.instance_state.#", "1"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_state.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_state.0.wait_restart_for_params"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.create_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.is_frozen"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.ha_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.pay_model"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.data_vip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.readable_node_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.readable_node_infos.0.data_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.readable_node_infos.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.readable_node_infos.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.proxy_ips.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.available_zones.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.available_zones.0.az_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.available_zones.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.available_zones.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.current_actions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.server_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttr(dataSource, "instances.0.network.#", "1"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.network.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.network.0.sub_net_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.network.0.security_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.ch_master_node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.node_num"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_instances" "test" {
  instance_id = "%s"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID)
}
