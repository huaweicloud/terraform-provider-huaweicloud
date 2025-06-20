package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsInstanceConfigurations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_instance_configurations.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsInstanceConfigurations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datastore_version_name"),
					resource.TestCheckResourceAttrSet(dataSource, "datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "created"),
					resource.TestCheckResourceAttrSet(dataSource, "updated"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.restart_required"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.readonly"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.value_range"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.description"),
				),
			},
		},
	})
}

func testDataSourceRdsInstanceConfigurations_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mysql.x1.large.2"
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  charging_mode     = "postPaid"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]

  db {
    type     = "MySQL"
    version  = "8.0"
    password = "Terraform145!"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testDataSourceRdsInstanceConfigurations_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_instance_configurations" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}
`, testDataSourceRdsInstanceConfigurations_base(name))
}
