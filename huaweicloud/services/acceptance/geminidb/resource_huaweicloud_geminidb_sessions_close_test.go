package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBSessionsClose_basic(t *testing.T) {
	resourceName := "huaweicloud_geminidb_sessions_close.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBSessionsClose_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_GEMINIDB_INSATNCE_ID),
				),
			},
		},
	})
}

func testAccGeminiDBSessionsClose_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_sessions_close" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}
