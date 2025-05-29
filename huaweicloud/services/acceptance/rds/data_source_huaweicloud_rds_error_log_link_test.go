package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsErrorLogLink_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_rds_error_log_link.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testRdsErrorLogLink_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "file_size"),
					resource.TestCheckResourceAttrSet(dataSource, "file_link"),
					resource.TestCheckResourceAttrSet(dataSource, "created_at"),
				),
			},
		},
	})
}

func testRdsErrorLogLink_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id

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

func testRdsErrorLogLink_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_error_log_link" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}
`, testRdsErrorLogLink_base(name))
}
