package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsConfigurationHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_configuration_histories.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsConfigurationHistories_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.parameter_name"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.old_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.new_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.update_result"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.applied"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.apply_time"),
					resource.TestCheckOutput("param_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsConfigurationHistories_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.medium.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
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

  parameters {
    name  = "deadlock_timeout"
    value = "10001"
  }
}
`, testAccRdsInstance_base(name), name)
}

func testDataSourceRdsConfigurationHistories_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_configuration_histories" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}

data "huaweicloud_rds_configuration_histories" "param_name_filter" {
  instance_id = huaweicloud_rds_instance.test.id
  param_name  = "deadlock_timeout"
}

output "param_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_configuration_histories.param_name_filter.histories) > 0
}
`, testDataSourceRdsConfigurationHistories_base(name))
}
