package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5AsymmetricSignatureSwitch_basic(t *testing.T) {
	rName := "huaweicloud_identityv5_asymmetric_signature_switch.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccV5AsymmetricSignatureSwitch_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "asymmetric_signature_switch", "true"),
				),
			},
			{
				Config: testAccV5AsymmetricSignatureSwitch_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "asymmetric_signature_switch", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccV5AsymmetricSignatureSwitch_basic_step1 = `
resource "huaweicloud_identityv5_asymmetric_signature_switch" "test" {
  asymmetric_signature_switch = true
}
`

const testAccV5AsymmetricSignatureSwitch_basic_step2 = `
resource "huaweicloud_identityv5_asymmetric_signature_switch" "test" {
  asymmetric_signature_switch = false
}
`
