package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceRdsCrossRegionRestoreTimeRanges_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_cross_region_restore_time_ranges.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsCrossRegionRestoreTimeRanges_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "restore_time.#"),
					resource.TestCheckResourceAttrSet(dataSource, "restore_time.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "restore_time.0.end_time"),

					resource.TestCheckOutput("date_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsCrossRegionRestoreTimeRanges_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "PostgreSQL"
  db_version    = "12"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    type    = "PostgreSQL"
    version = "12"
  }
    
  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_backup" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_rds_instance.test.id
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceRdsCrossRegionRestoreTimeRanges_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_cross_region_restore_time_ranges" "test" {
  depends_on = [huaweicloud_rds_backup.test]

  instance_id = huaweicloud_rds_instance.test.id
}

data "huaweicloud_rds_cross_region_restore_time_ranges" "date_filter" {
  depends_on = [huaweicloud_rds_backup.test]

  instance_id = huaweicloud_rds_instance.test.id
  date        = split("T", huaweicloud_rds_instance.test.created)[0]
}
output "date_filter_is_useful" {
  value = length(data.huaweicloud_rds_cross_region_restore_time_ranges.date_filter.restore_time) > 0
}
`, testDataSourceRdsCrossRegionRestoreTimeRanges_base(name))
}
