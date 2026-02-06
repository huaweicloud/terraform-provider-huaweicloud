package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV5PasswordPolicyFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5PasswordPolicy(client)
}

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5PasswordPolicy_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_identityv5_password_policy.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5PasswordPolicyFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV5PasswordPolicy_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "maximum_consecutive_identical_chars", "1"),
					resource.TestCheckResourceAttr(rName, "minimum_password_age", "30"),
					resource.TestCheckResourceAttr(rName, "minimum_password_length", "10"),
					resource.TestCheckResourceAttr(rName, "password_reuse_prevention", "2"),
					resource.TestCheckResourceAttr(rName, "password_not_username_or_invert", "false"),
					resource.TestCheckResourceAttr(rName, "password_validity_period", "90"),
					resource.TestCheckResourceAttr(rName, "password_char_combination", "3"),
					resource.TestCheckResourceAttr(rName, "allow_user_to_change_password", "false"),
				),
			},
			{
				Config: testAccV5PasswordPolicy_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "maximum_consecutive_identical_chars", "2"),
					resource.TestCheckResourceAttr(rName, "minimum_password_age", "60"),
					resource.TestCheckResourceAttr(rName, "minimum_password_length", "20"),
					resource.TestCheckResourceAttr(rName, "password_reuse_prevention", "4"),
					resource.TestCheckResourceAttr(rName, "password_not_username_or_invert", "true"),
					resource.TestCheckResourceAttr(rName, "password_validity_period", "120"),
					resource.TestCheckResourceAttr(rName, "password_char_combination", "4"),
					resource.TestCheckResourceAttr(rName, "allow_user_to_change_password", "true"),
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

const testAccV5PasswordPolicy_basic_step1 = `
resource "huaweicloud_identityv5_password_policy" "test" {
  maximum_consecutive_identical_chars = 1
  minimum_password_age                = 30
  minimum_password_length             = 10
  password_reuse_prevention           = 2
  password_not_username_or_invert     = false
  password_validity_period            = 90
  password_char_combination           = 3
  allow_user_to_change_password       = false
}
`

const testAccV5PasswordPolicy_basic_step2 = `
resource "huaweicloud_identityv5_password_policy" "test" {
  maximum_consecutive_identical_chars = 2
  minimum_password_age                = 60
  minimum_password_length             = 20
  password_reuse_prevention           = 4
  password_not_username_or_invert     = true
  password_validity_period            = 120
  password_char_combination           = 4
  allow_user_to_change_password       = true
}
`
