package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsBackupFiles_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_backup_files.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsBackupFiles_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "files.#"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.download_link"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.link_expired_time"),
					resource.TestCheckResourceAttrSet(dataSource, "bucket"),
				),
			},
		},
	})
}

func testDataSourceRdsBackupFiles_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 8630
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "huaweicloud_rds_backup" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_rds_instance.test.id
}
`, testAccRdsInstance_base(), name)
}

func testDataSourceRdsBackupFiles_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_backup_files" "test" {
  backup_id = huaweicloud_rds_backup.test.id
}
`, testDataSourceRdsBackupFiles_base(name))
}
