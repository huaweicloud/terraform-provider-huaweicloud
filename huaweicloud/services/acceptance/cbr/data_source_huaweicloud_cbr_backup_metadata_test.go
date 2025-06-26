package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBackupMetadata_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cbr_backup_metadata.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCBRBackupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBackupMetadata_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "backup_id", acceptance.HW_CBR_BACKUP_ID),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavor"),
					resource.TestCheckResourceAttrSet(dataSourceName, "floatingips.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "interface"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ports.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.#"),
				),
			},
		},
	})
}

func testAccDataSourceBackupMetadata_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cbr_backup_metadata" "test" {
  backup_id = "%s"
}
`, acceptance.HW_CBR_BACKUP_ID)
}
