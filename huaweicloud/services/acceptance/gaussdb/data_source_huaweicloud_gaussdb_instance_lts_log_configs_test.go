package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBInstanceLtsLogConfigs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_instance_lts_log_configs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDBInstanceLtsLogConfigs_basic(name),
			},
			{
				Config: testAccDataSourceGaussDBInstanceLtsLogConfigs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.datastore.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.frozen_flag"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.0.log_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.0.lts_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.0.lts_stream_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.0.enabled"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDBInstanceLtsLogConfigs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_gaussdb_instance_lts_log_associate" "test" {
  instance_id    = huaweicloud_gaussdb_instance.test.id
  log_type       = "audit_log"
  lts_group_id   = huaweicloud_lts_group.test.id
  lts_stream_id  = huaweicloud_lts_stream.test.id
}

data "huaweicloud_gaussdb_instance_lts_log_configs" "test" {
}

data "huaweicloud_gaussdb_instance_lts_log_configs" "instance_id_filter" {
  instance_id   = huaweicloud_gaussdb_instance.test.id
  instance_mode = "Ha"
  instance_name = "%[2]s"
}
output "instance_id_filter" {
  value = length(data.huaweicloud_gaussdb_instance_lts_log_configs.instance_id_filter.instance_lts_configs) > 0
}
`, testDataSourceGaussdbInstanceMetrics_base(name), name)
}
