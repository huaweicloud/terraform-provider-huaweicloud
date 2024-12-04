package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAccessCredential_basic(t *testing.T) {
	rName := "huaweicloud_iotda_access_credential.test"

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessCredential_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "type", "AMQP"),
					resource.TestCheckResourceAttr(rName, "force_disconnect", "false"),
					resource.TestCheckResourceAttrSet(rName, "access_key"),
					resource.TestCheckResourceAttrSet(rName, "access_code"),
				),
			},
		},
	})
}

func testAccAccessCredential_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_access_credential" "test" {
  type             = "AMQP"
  force_disconnect = false
}
`, buildIoTDAEndpoint())
}
