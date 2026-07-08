package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBInstanceDatabaseBatchUpgrade_basic(t *testing.T) {
	resourceName := "huaweicloud_taurusdb_instance_database_batch_upgrade.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBInstanceDatabaseBatchUpgrade_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccTaurusDBInstanceDatabaseBatchUpgrade_basic() string {
	return `
data "huaweicloud_taurusdb_instances" "test" {}

locals {
  databases_instance_infos = [
    for instance in data.huaweicloud_taurusdb_instances.test.instances : {
      instance_id     = instance.id
      current_version = "2.0.75.28"
    }
  ]
}

resource "huaweicloud_taurusdb_instance_database_batch_upgrade" "test" {
  dynamic "databases_instance_infos" {
    for_each = local.databases_instance_infos

    content {
      instance_id     = databases_instance_infos.value.instance_id
      current_version = databases_instance_infos.value.current_version
    }
  }

  delay = false
}
`
}
