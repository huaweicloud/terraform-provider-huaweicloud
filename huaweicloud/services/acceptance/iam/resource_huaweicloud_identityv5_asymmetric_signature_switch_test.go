package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityV5AsymmetricSignatureSwitch_basic(t *testing.T) {
	resourceName := "huaweicloud_identityv5_asymmetric_signature_switch.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV5AsymmetricSignatureSwitch_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "asymmetric_signature_switch", "true"),
				),
			},
			{
				Config: testAccIdentityV5AsymmetricSignatureSwitch_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "asymmetric_signature_switch", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccIdentityV5AsymmetricSignatureSwitch_basic = `
resource "huaweicloud_identityv5_asymmetric_signature_switch" "test" {
  asymmetric_signature_switch = true
}
`

var testAccIdentityV5AsymmetricSignatureSwitch_update = `
resource "huaweicloud_identityv5_asymmetric_signature_switch" "test" {
  asymmetric_signature_switch = false
}
`
