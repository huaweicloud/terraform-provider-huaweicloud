package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBInstanceUpgrade_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBInstanceUpgrade_basic(),
			},
		},
	})
}

func testAccTaurusDBInstanceUpgrade_basic() string {
	return fmt.Sprintf(`


resource "huaweicloud_taurusdb_instance_upgrade" "test" {
  instance_id = "%[1]s"
  delay       = true
}`, acceptance.HW_TAURUSDB_INSTANCE_ID)
}
