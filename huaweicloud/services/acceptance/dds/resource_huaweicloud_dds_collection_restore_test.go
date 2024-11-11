package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSCollectionRestore_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSCollectionRestore_basic(),
			},
		},
	})
}

func testAccDDSCollectionRestore_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_collection_restore" "test" {
  instance_id  = "%[2]s"

  restore_collections {
    database = try(data.huaweicloud_dds_restore_databases.test.databases[0], "")

    collections {
      old_name                = try(data.huaweicloud_dds_restore_collections.test.collections[0], "")
      restore_collection_time = local.end_time
    }
  }

  lifecycle {
    ignore_changes = [
	  restore_collections.0.collections.0.restore_collection_time,
    ]
  }
}`, testDataSourceDdsRestoreCollections_basic(), acceptance.HW_DDS_INSTANCE_ID)
}
