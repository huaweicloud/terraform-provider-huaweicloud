package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsPgRoles_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_pg_roles.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsPgRoles_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "roles.#"),
					resource.TestCheckResourceAttr(dataSource, "roles.#", "1"),
					resource.TestCheckResourceAttrPair(dataSource, "roles.0",
						"huaweicloud_rds_pg_account.test", "name"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsPgRoles_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name               = "%[2]s"
  flavor             = "rds.pg.n1.large.2"
  availability_zone  = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id  = data.huaweicloud_networking_secgroup.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  vpc_id             = data.huaweicloud_vpc.test.id
  time_zone          = "UTC+08:00"

  db {
    type    = "PostgreSQL"
    version = "12"
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testAccRdsInstance_base(), name)
}

func testDataSourceDataSourceRdsPgRoles_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "Test@12345678"
}

data "huaweicloud_rds_pg_roles" "test" {
  depends_on = [huaweicloud_rds_pg_account.test]

  instance_id = huaweicloud_rds_instance.test.id
  account     = "root"
}
`, testDataSourceDataSourceRdsPgRoles_base(name), name)
}
