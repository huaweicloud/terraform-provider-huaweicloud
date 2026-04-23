package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceBackupFileLinks_basic(t *testing.T) {
	rName := "data.huaweicloud_dcs_backup_file_links.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDCSInstanceID(t)
			acceptance.TestAccPreCheckDcsBackupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceBackupFileLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "links.#"),
					resource.TestCheckResourceAttrSet(rName, "links.0.file_name"),
					resource.TestCheckResourceAttrSet(rName, "links.0.link"),
					resource.TestCheckResourceAttrSet(rName, "bucket_name"),
					resource.TestCheckResourceAttrSet(rName, "file_path"),
				),
			},
		},
	})
}

func testAccDatasourceBackupFileLinks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dcs_backup_file_links" "test" {
  instance_id = "%[1]s"
  backup_id   = "%[2]s"
  expiration  = 3600
}
`, acceptance.HW_DCS_INSTANCE_ID, acceptance.HW_DCS_BACKUP_ID)
}
