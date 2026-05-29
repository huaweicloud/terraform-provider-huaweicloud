package das

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccHistoryTransactionSwitch_basic(t *testing.T) {
	var (
		rName = "huaweicloud_das_history_transaction_switch.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccHistoryTransactionSwitch_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "status", "Enabled"),
					resource.TestCheckResourceAttr(rName, "datastore_type", "MySQL"),
				),
			},
			{
				Config: testAccHistoryTransactionSwitch_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "status", "Disabled"),
				),
			},
		},
	})
}

func testAccHistoryTransactionSwitch_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccHistoryTransactionSwitch_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_history_transaction_switch" "test" {
  instance_id    = local.instance_ids[0]
  status         = "Enabled"
  datastore_type = "MySQL"
}
`, testAccHistoryTransactionSwitch_base())
}

func testAccHistoryTransactionSwitch_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_history_transaction_switch" "test" {
  instance_id    = local.instance_ids[0]
  status         = "Disabled"
  datastore_type = "MySQL"

  enable_force_new = "true"
}
`, testAccHistoryTransactionSwitch_base())
}
