package iam

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccV5PasswordPolicy_basic(t *testing.T) {
	resourceName := "huaweicloud_identityv5_password_policy.enhanced"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccV5CheckPasswordPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccV5PasswordPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccV5CheckPasswordPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "maximum_consecutive_identical_chars", "1"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_age", "30"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_length", "10"),
					resource.TestCheckResourceAttr(resourceName, "password_reuse_prevention", "2"),
					resource.TestCheckResourceAttr(resourceName, "password_not_username_or_invert", "false"),
					resource.TestCheckResourceAttr(resourceName, "password_validity_period", "90"),
					resource.TestCheckResourceAttr(resourceName, "password_char_combination", "3"),
					resource.TestCheckResourceAttr(resourceName, "allow_user_to_change_password", "false"),
				),
			},
			{
				Config: testAccV5PasswordPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "maximum_consecutive_identical_chars", "2"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_age", "60"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_length", "20"),
					resource.TestCheckResourceAttr(resourceName, "password_reuse_prevention", "4"),
					resource.TestCheckResourceAttr(resourceName, "password_not_username_or_invert", "true"),
					resource.TestCheckResourceAttr(resourceName, "password_validity_period", "120"),
					resource.TestCheckResourceAttr(resourceName, "password_char_combination", "4"),
					resource.TestCheckResourceAttr(resourceName, "allow_user_to_change_password", "true"),
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

func testAccV5CheckPasswordPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return errors.New("resource ID is not set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		region := acceptance.HW_REGION_NAME
		client, err := cfg.IAMNoVersionClient(region)
		if err != nil {
			return fmt.Errorf("error creating IAM client: %s", err)
		}

		getPasswordPolicyHttpUrl := "v5/password-policy"
		getPasswordPolicyPath := client.Endpoint + getPasswordPolicyHttpUrl
		getPasswordPolicyOpt := golangsdk.RequestOpts{KeepResponseBody: true}
		_, err = client.Request("GET", getPasswordPolicyPath, &getPasswordPolicyOpt)
		if err != nil {
			return fmt.Errorf("error retrieving IAM password policy: %s", err)
		}
		return err
	}
}

func testAccV5CheckPasswordPolicyDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	region := acceptance.HW_REGION_NAME
	client, err := cfg.IAMNoVersionClient(region)
	if err != nil {
		return fmt.Errorf("error creating IAM client: %s", err)
	}

	getPasswordPolicyHttpUrl := "v5/password-policy"
	getPasswordPolicyPath := client.Endpoint + getPasswordPolicyHttpUrl
	getPasswordPolicyOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identityv5_password_policy" {
			continue
		}

		getPasswordPolicyResp, err := client.Request("GET", getPasswordPolicyPath, &getPasswordPolicyOpt)
		if err != nil {
			return fmt.Errorf("error fetching the IAM account password policy: %s", err)
		}

		getPasswordPolicyBody, err := utils.FlattenResponse(getPasswordPolicyResp)
		if err != nil {
			return err
		}

		if utils.PathSearch("password_policy.maximum_consecutive_identical_chars", getPasswordPolicyBody, float64(0)).(float64) != 0 ||
			utils.PathSearch("password_policy.minimum_password_age", getPasswordPolicyBody, float64(0)).(float64) != 0 ||
			utils.PathSearch("password_policy.minimum_password_length", getPasswordPolicyBody, float64(0)).(float64) != 8 ||
			utils.PathSearch("password_policy.password_reuse_prevention", getPasswordPolicyBody, float64(0)).(float64) != 1 ||
			!utils.PathSearch("password_policy.password_not_username_or_invert", getPasswordPolicyBody, false).(bool) ||
			utils.PathSearch("password_policy.password_validity_period", getPasswordPolicyBody, float64(0)).(float64) != 0 ||
			utils.PathSearch("password_policy.password_char_combination", getPasswordPolicyBody, float64(0)).(float64) != 2 ||
			!utils.PathSearch("password_policy.allow_user_to_change_password", getPasswordPolicyBody, false).(bool) {
			return errors.New("the password policy failed to reset to defaults")
		}
	}
	return nil
}

func testAccV5PasswordPolicy_basic() string {
	return `
resource "huaweicloud_identityv5_password_policy" "enhanced" {
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
}

func testAccV5PasswordPolicy_update() string {
	return `
resource "huaweicloud_identityv5_password_policy" "enhanced" {
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
}
