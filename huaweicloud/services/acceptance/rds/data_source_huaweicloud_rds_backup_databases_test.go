package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsBackupDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_backup_databases.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRdsBackupDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "bucket_name"),
					resource.TestCheckResourceAttrSet(dataSource, "database_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.database_name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.backup_file_size"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.backup_file_name"),
				),
			},
		},
	})
}

func testAccDataSourceRdsBackupDatabases_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_backup_databases" "test" {
  instance_id = "%s"
  backup_id   = "%s"
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_BACKUP_ID)
}
