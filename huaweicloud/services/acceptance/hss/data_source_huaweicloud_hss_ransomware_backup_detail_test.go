package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRansomwareBackupDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_ransomware_backup_detail.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a backup ID for the ransomware protection host.
			acceptance.TestAccPreCheckHSSBackupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRansomwareBackupDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vault_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.image_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.vault_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.resource_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.resource_name"),
				),
			},
		},
	})
}

func testDataSourceRansomwareBackupDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_ransomware_backup_detail" "test" {
  backup_id             = "%s"
  enterprise_project_id = "all_granted_eps"
}
`, acceptance.HW_HSS_BACKUP_ID)
}
