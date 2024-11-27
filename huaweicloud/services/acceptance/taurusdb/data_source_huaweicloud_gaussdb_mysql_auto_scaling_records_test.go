package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbMysqlAutoScalingRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_auto_scaling_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlAutoScalingRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlAutoScalingRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_mysql_flavors" "test" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = "multi"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name                  = "%[2]s"
  password              = "Test@12345678"
  flavor                = data.huaweicloud_gaussdb_mysql_flavors.test.flavors[0].name
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"
}

data "huaweicloud_gaussdb_mysql_auto_scaling_records" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
}
`, common.TestBaseNetwork(name), name)
}
