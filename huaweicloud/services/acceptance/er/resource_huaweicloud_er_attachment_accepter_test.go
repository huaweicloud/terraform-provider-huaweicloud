package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAttachmentAccepter_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckErSharedAttachmentAccepter(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachmentAccepter_basic(),
			},
		},
	})
}

func testAccAttachmentAccepter_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_er_attachment_accepter" "test" {
  instance_id   = "%[1]s"
  attachment_id = "%[2]s"
  action        = "accept"
} 
`, acceptance.HW_ER_SHARED_INSTANCE_ID, acceptance.HW_ER_SHARED_ATTACHMENT_ID)
}
