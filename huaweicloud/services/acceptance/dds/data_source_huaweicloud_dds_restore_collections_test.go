package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsRestoreCollections_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_restore_collections.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdsRestoreCollections_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "collections.#"),
				),
			},
		},
	})
}

func testDataSourceDdsRestoreCollections_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dds_restore_collections" "test" {
  instance_id  = "%[2]s"
  db_name      = try(data.huaweicloud_dds_restore_databases.test.databases[0], "")
  restore_time = local.end_time
}
`, testDataSourceDdsRestoreDatabases_basic(), acceptance.HW_DDS_INSTANCE_ID)
}
