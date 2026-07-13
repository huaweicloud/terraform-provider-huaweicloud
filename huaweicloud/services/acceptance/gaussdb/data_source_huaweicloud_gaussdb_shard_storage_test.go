package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBShardStorageDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_shard_storage.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBShardStorageDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "group_disk_infos.#"),
				),
			},
		},
	})
}

func testAccGaussDBShardStorageDataSource_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_shard_storage" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
