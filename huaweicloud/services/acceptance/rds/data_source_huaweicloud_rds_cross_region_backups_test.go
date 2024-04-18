package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsCrossRegionBackups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_cross_region_backups.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsCrossRegionBackupInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsCrossRegionBackups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.version"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsCrossRegionBackups_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_cross_region_backups" "test" {
  instance_id = "%[1]s"
  backup_type = "auto"
}
`, acceptance.HW_RDS_CROSS_REGION_BACKUP_INSTANCE_ID)
}
