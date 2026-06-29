package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksLtsConfigs_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_lts_configs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBHtapStarrocksLtsConfigs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.0.log_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.0.lts_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.0.lts_stream_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.0.enabled"),
					resource.TestCheckOutput("filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksLtsConfigs_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  count = 2

  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s_${count.index}"
}

resource "huaweicloud_taurusdb_htap_starrocks_lts_config" "test" {
  instance_id   = "%[2]s"
  log_type      = "error_log"
  lts_group_id  = huaweicloud_lts_group.test.id
  lts_stream_id = huaweicloud_lts_stream.test[0].id
}

resource "huaweicloud_taurusdb_htap_starrocks_lts_config" "test_slow" {
  instance_id   = "%[2]s"
  log_type      = "slow_log"
  lts_group_id  = huaweicloud_lts_group.test.id
  lts_stream_id = huaweicloud_lts_stream.test[1].id
}
`, rName, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}

func testAccDataSourceTaurusDBHtapStarrocksLtsConfigs_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_taurusdb_htap_starrocks_lts_configs" "test" {
  depends_on = [huaweicloud_taurusdb_htap_starrocks_lts_config.test, huaweicloud_taurusdb_htap_starrocks_lts_config.test_slow]
}

locals {
  instance_id    = data.huaweicloud_taurusdb_htap_starrocks_lts_configs.test.instance_lts_configs[0].instance[0].id
  instance_name  = data.huaweicloud_taurusdb_htap_starrocks_lts_configs.test.instance_lts_configs[0].instance[0].name
  eps_project_id = data.huaweicloud_taurusdb_htap_starrocks_lts_configs.test.instance_lts_configs[0].instance[0].enterprise_project_id
}

data "huaweicloud_taurusdb_htap_starrocks_lts_configs" "filter" {
  instance_id           = local.instance_id
  instance_name         = local.instance_name
  enterprise_project_id = local.eps_project_id
}

locals {
  filtered_lts_configs  = data.huaweicloud_taurusdb_htap_starrocks_lts_configs.filter.instance_lts_configs
  filtered_lts_instance = data.huaweicloud_taurusdb_htap_starrocks_lts_configs.filter.instance_lts_configs[0].instance[0]
}

output "filter_is_useful" {
  value = length(local.filtered_lts_configs) > 0 && (local.filtered_lts_instance.id == local.instance_id
  && local.filtered_lts_instance.name == local.instance_name && local.filtered_lts_instance.enterprise_project_id == local.eps_project_id)
}
`, testAccTaurusDBHtapStarrocksLtsConfigs_base(rName))
}
