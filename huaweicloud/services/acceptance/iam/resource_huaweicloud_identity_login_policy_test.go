package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccLoginPolicy_basic(t *testing.T) {
	resourceName := "huaweicloud_identity_login_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLoginPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoginPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoginPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_validity_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "lockout_duration", "30"),
					resource.TestCheckResourceAttr(resourceName, "login_failed_times", "10"),
					resource.TestCheckResourceAttr(resourceName, "period_with_login_failures", "30"),
					resource.TestCheckResourceAttr(resourceName, "session_timeout", "120"),
					resource.TestCheckResourceAttr(resourceName, "show_recent_login_info", "true"),
					resource.TestCheckResourceAttr(resourceName, "custom_info_for_login", "hello Terraform"),
				),
			},
			{
				Config: testAccLoginPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_validity_period", "0"),
					resource.TestCheckResourceAttr(resourceName, "lockout_duration", "20"),
					resource.TestCheckResourceAttr(resourceName, "login_failed_times", "6"),
					resource.TestCheckResourceAttr(resourceName, "period_with_login_failures", "20"),
					resource.TestCheckResourceAttr(resourceName, "session_timeout", "90"),
					resource.TestCheckResourceAttr(resourceName, "show_recent_login_info", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_info_for_login", ""),
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

func testAccCheckLoginPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		region := acceptance.HW_REGION_NAME
		getLoginPolicyProduct := "iam"
		getLoginPolicyClient, err := cfg.NewServiceClient(getLoginPolicyProduct, region)
		if err != nil {
			return fmt.Errorf("error creating IAM Client: %s", err)
		}

		getLoginPolicyHttpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy"
		getLoginPolicyPath := getLoginPolicyClient.Endpoint + getLoginPolicyHttpUrl
		getLoginPolicyPath = strings.ReplaceAll(getLoginPolicyPath, "{domain_id}", cfg.DomainID)
		getLoginPolicyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		_, err = getLoginPolicyClient.Request("GET", getLoginPolicyPath, &getLoginPolicyOpt)
		if err != nil {
			return fmt.Errorf("error retrieving IAM login policy: %s", err)
		}
		return err
	}
}

func testAccCheckLoginPolicyDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	region := acceptance.HW_REGION_NAME
	getLoginPolicyProduct := "iam"
	getLoginPolicyClient, err := cfg.NewServiceClient(getLoginPolicyProduct, region)
	if err != nil {
		return fmt.Errorf("error creating IAM Client: %s", err)
	}

	getLoginPolicyHttpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy"
	getLoginPolicyPath := getLoginPolicyClient.Endpoint + getLoginPolicyHttpUrl
	getLoginPolicyPath = strings.ReplaceAll(getLoginPolicyPath, "{domain_id}", cfg.DomainID)
	getLoginPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_login_policy" {
			continue
		}

		getLoginPolicyResp, err := getLoginPolicyClient.Request("GET", getLoginPolicyPath, &getLoginPolicyOpt)
		if err != nil {
			return fmt.Errorf("error retrieving IAM login policy: %s", err)
		}

		getLoginPolicyRespBody, err := utils.FlattenResponse(getLoginPolicyResp)
		if err != nil {
			return err
		}

		if utils.PathSearch("login_policy.account_validity_period", getLoginPolicyRespBody, float64(0)).(float64) != 0 ||
			utils.PathSearch("login_policy.custom_info_for_login", getLoginPolicyRespBody, "").(string) != "" ||
			utils.PathSearch("login_policy.lockout_duration", getLoginPolicyRespBody, float64(0)).(float64) != 15 ||
			utils.PathSearch("login_policy.login_failed_times", getLoginPolicyRespBody, float64(0)).(float64) != 5 ||
			utils.PathSearch("login_policy.period_with_login_failures", getLoginPolicyRespBody, float64(0)).(float64) != 15 ||
			utils.PathSearch("login_policy.session_timeout", getLoginPolicyRespBody, float64(0)).(float64) != 60 ||
			utils.PathSearch("login_policy.show_recent_login_info", getLoginPolicyRespBody, false).(bool) {
			return fmt.Errorf("the login policy failed to reset to defaults")
		}
	}

	return nil
}

func testAccLoginPolicy_basic() string {
	return `
resource "huaweicloud_identity_login_policy" "test" {
  account_validity_period    = 20
  lockout_duration           = 30
  login_failed_times         = 10
  period_with_login_failures = 30
  session_timeout            = 120
  show_recent_login_info     = true
  custom_info_for_login      = "hello Terraform"
}
`
}

func testAccLoginPolicy_update() string {
	return `
resource "huaweicloud_identity_login_policy" "test" {
  account_validity_period    = 0
  lockout_duration           = 20
  login_failed_times         = 6
  period_with_login_failures = 20
  session_timeout            = 90
  show_recent_login_info     = false
  custom_info_for_login      = ""
}
`
}
