package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCBRVaultChangeChargeMode_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceVaultChangeChargeMode_basic(),
			},
		},
	})
}

func testResourceVaultChangeChargeMode_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name            = "%[1]s"
  type            = "disk"
  protection_type = "backup"
  size            = 50
}

resource "huaweicloud_cbr_vault_change_charge_mode" "test" {
  vault_ids     = [huaweicloud_cbr_vault.test.id]
  charging_mode = "pre_paid"
  period_type   = "month"
  period_num    = 1
  is_auto_renew = true
}
`, name)
}
