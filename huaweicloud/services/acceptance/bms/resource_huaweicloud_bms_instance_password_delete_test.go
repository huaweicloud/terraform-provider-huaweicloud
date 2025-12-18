package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBmsInstancePasswordDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckBmsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstancePasswordDelete_basic(),
			},
		},
	})
}

func testAccBmsInstancePasswordDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_bms_instance_password_delete" "test" {
  server_id = "%s"
}
`, acceptance.HW_BMS_INSTANCE_ID)
}
