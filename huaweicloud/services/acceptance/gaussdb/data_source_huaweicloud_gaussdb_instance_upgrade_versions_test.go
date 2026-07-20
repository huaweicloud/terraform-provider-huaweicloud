package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBKernelVersionUpgradeDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_instance_upgrade_versions.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBKernelVersionUpgradeDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "upgrade_type_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "upgrade_type_list.0.upgrade_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "upgrade_type_list.0.enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "upgrade_type_list.0.is_parallel_upgrade"),
					resource.TestCheckResourceAttrSet(dataSourceName, "upgrade_type_list.0.upgrade_action_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "upgrade_type_list.0.upgrade_action_list.0.upgrade_action"),
					resource.TestCheckResourceAttrSet(dataSourceName, "upgrade_type_list.0.upgrade_action_list.0.enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hotfix_upgrade_infos.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hotfix_upgrade_infos.0.version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hotfix_upgrade_infos.0.common_patch"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hotfix_upgrade_infos.0.backup_sensitive"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hotfix_upgrade_infos.0.default_upgrade"),
				),
			},
		},
	})
}

func testAccGaussDBKernelVersionUpgradeDataSource_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_upgrade_versions" "test" {
  instance_ids = split(",", "%s")
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
