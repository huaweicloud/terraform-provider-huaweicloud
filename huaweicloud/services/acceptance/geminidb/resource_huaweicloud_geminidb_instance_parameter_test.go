package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBInstanceParameter_basic(t *testing.T) {
	resourceName := "huaweicloud_geminidb_instance_parameter.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstanceParameter_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "AuthFailLockTime"),
					resource.TestCheckResourceAttr(resourceName, "value", "6"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "value_range"),
				),
			},
			{
				Config: testAccGeminiDBInstanceParameter_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "AuthFailLockTime"),
					resource.TestCheckResourceAttr(resourceName, "value", "5"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func testAccGeminiDBInstanceParameter_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_instance_parameter" "test" {
  instance_id = "%s"
  name        = "AuthFailLockTime"
  value       = "6"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}

func testAccGeminiDBInstanceParameter_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_instance_parameter" "test" {
  instance_id = "%s"
  name        = "AuthFailLockTime"
  value       = "5"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}
