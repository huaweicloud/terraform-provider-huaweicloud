package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsLtsConfigs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_lts_configs.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRdsLtsConfigsBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.instance.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_lts_configs.0.lts_configs.#"),

					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_lts_configs.enterprise_project_id_filter", "instance_lts_configs.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_lts_configs.instance_id_filter", "instance_lts_configs.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_lts_configs.instance_name_filter", "instance_lts_configs.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_lts_configs.instance_status_filter", "instance_lts_configs.#"),
				),
			},
		},
	})
}
func testDataSourceRdsLtsConfigs_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.medium.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"
  charging_mode     = "postPaid"

  db {
    type    = "PostgreSQL"
    version = "16"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccDataSourceRdsLtsConfigsBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_group" "test" {
	group_name = "test"
	ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
	group_id = huaweicloud_lts_group.test.id
	stream_name = "test"
}

resource "huaweicloud_rds_lts_config" "test" {
	instance_id   = huaweicloud_rds_instance.test.id
	engine        = "postgresql"
	log_type      = "error_log"
	lts_group_id  = huaweicloud_lts_group.test.id
	lts_stream_id = huaweicloud_lts_stream.test.id
}


data "huaweicloud_rds_lts_configs" "test" {
	engine = "postgresql"
}

data "huaweicloud_rds_lts_configs" "enterprise_project_id_filter" {
	engine                = "postgresql"
	enterprise_project_id = "%s"
}

output "enterprise_project_id_filter_is_useful" {
	value = length(data.huaweicloud_rds_lts_configs.enterprise_project_id_filter.instance_lts_configs) > 0
}

data "huaweicloud_rds_lts_configs" "instance_id_filter" {
	engine       = "postgresql"
	instance_id  = huaweicloud_rds_instance.test.id
}

output "instance_id_filter_is_useful" {
	value = length(data.huaweicloud_rds_lts_configs.instance_id_filter.instance_lts_configs) > 0
}

data "huaweicloud_rds_lts_configs" "instance_name_filter" {
	engine         = "postgresql"
	instance_name  = huaweicloud_rds_instance.test.name
}

output "instance_name_filter_is_useful" {
	value = length(data.huaweicloud_rds_lts_configs.instance_name_filter.instance_lts_configs) > 0
}

data "huaweicloud_rds_lts_configs" "instance_status_filter" {
	engine           = "postgresql"
	instance_status  = "normal"
}

output "instance_status_filter_is_useful" {
	value = length(data.huaweicloud_rds_lts_configs.instance_status_filter.instance_lts_configs) > 0
}
`, testDataSourceRdsLtsConfigs_base(name), acceptance.HW_ENTERPRISE_PROJECT_ID)
}
