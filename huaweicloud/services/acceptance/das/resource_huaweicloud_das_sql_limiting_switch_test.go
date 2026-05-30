package das

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSqlLimitingSwitch_basic(t *testing.T) {
	var (
		rName = "huaweicloud_das_sql_limiting_switch.test"
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
				Config: testAccSqlLimitingSwitch_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "status", "ON"),
					resource.TestCheckResourceAttr(rName, "datastore_type", "MySQL"),
				),
			},
			{
				Config: testAccSqlLimitingSwitch_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "status", "OFF"),
				),
			},
		},
	})
}

func testAccSqlLimitingSwitch_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccSqlLimitingSwitch_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_sql_limiting_switch" "test" {
  instance_id    = local.instance_ids[0]
  status         = "ON"
  datastore_type = "MySQL"
}
`, testAccSqlLimitingSwitch_base())
}

func testAccSqlLimitingSwitch_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_sql_limiting_switch" "test" {
  instance_id    = local.instance_ids[0]
  status         = "OFF"
  datastore_type = "MySQL"

  enable_force_new = "true"
}
`, testAccSqlLimitingSwitch_base())
}
