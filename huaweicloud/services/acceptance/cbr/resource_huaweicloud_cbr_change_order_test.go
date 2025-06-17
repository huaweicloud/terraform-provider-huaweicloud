package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Because the API capabilities are unstable, sometimes the API will report an error "The system is busy. Please try later."
func TestAccResourceChangeOrder_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceChangeOrder_basic(),
			},
		},
	})
}

func testResourceChangeOrder_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name             = "%[1]s"
  type             = "server"
  consistent_level = "crash_consistent"
  protection_type  = "backup"
  size             = 100

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"

  lifecycle {
    ignore_changes = [
      size,
    ]
  }
}

resource "huaweicloud_cbr_change_order" "test" {
  resource_id = huaweicloud_cbr_vault.test.id

  product_info {
    product_id               = "00301-34090-0--0"
    resource_size            = 200
    resource_size_measure_id = 17
    resource_spec_code       = "vault.backup.server.normal"
  }
}
`, name)
}
