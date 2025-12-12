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

func TestAccV5LoginPolicy_basic(t *testing.T) {
	resourceName := "huaweicloud_identityv5_login_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccV5CheckLoginPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccV5LoginPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccV5CheckLoginPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_validity_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "lockout_duration", "30"),
					resource.TestCheckResourceAttr(resourceName, "login_failed_times", "10"),
					resource.TestCheckResourceAttr(resourceName, "period_with_login_failures", "30"),
					resource.TestCheckResourceAttr(resourceName, "session_timeout", "120"),
					resource.TestCheckResourceAttr(resourceName, "show_recent_login_info", "true"),
					resource.TestCheckResourceAttr(resourceName, "custom_info_for_login", "hello Terraform"),
					resource.TestCheckResourceAttr(resourceName, "allow_address_netmasks.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allow_address_netmasks.0.address_netmask", "255.0.0.0/1"),
					resource.TestCheckResourceAttr(resourceName, "allow_address_netmasks.0.description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges_ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges_ipv6.0.ip_range",
						"0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFF0"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges_ipv6.0.description", "terraform test"),
				),
			},
			{
				Config: testAccV5LoginPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_validity_period", "0"),
					resource.TestCheckResourceAttr(resourceName, "lockout_duration", "20"),
					resource.TestCheckResourceAttr(resourceName, "login_failed_times", "6"),
					resource.TestCheckResourceAttr(resourceName, "period_with_login_failures", "20"),
					resource.TestCheckResourceAttr(resourceName, "session_timeout", "90"),
					resource.TestCheckResourceAttr(resourceName, "show_recent_login_info", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_info_for_login", ""),
					resource.TestCheckResourceAttr(resourceName, "allow_address_netmasks.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allow_address_netmasks.0.address_netmask", "100.0.0.0/1"),
					resource.TestCheckResourceAttr(resourceName, "allow_address_netmasks.0.description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges.0.ip_range", "0.0.0.0-255.255.255.100"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges.0.description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges_ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges_ipv6.0.ip_range",
						"0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFF1"),
					resource.TestCheckResourceAttr(resourceName, "allow_ip_ranges_ipv6.0.description", "terraform test"),
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

func testAccV5CheckLoginPolicyExists(n string) resource.TestCheckFunc {
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
		getLoginPolicyClient, err := cfg.IAMNoVersionClient(region)
		if err != nil {
			return fmt.Errorf("error creating IAM Client: %s", err)
		}

		getLoginPolicyHttpUrl := "v5/login-policy"
		getLoginPolicyPath := getLoginPolicyClient.Endpoint + getLoginPolicyHttpUrl
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

func testAccV5CheckLoginPolicyDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	region := acceptance.HW_REGION_NAME
	getLoginPolicyClient, err := cfg.IAMNoVersionClient(region)
	if err != nil {
		return fmt.Errorf("error creating IAM Client: %s", err)
	}

	getLoginPolicyHttpUrl := "v5/login-policy"
	getLoginPolicyPath := getLoginPolicyClient.Endpoint + getLoginPolicyHttpUrl
	getLoginPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identityv5_login_policy" {
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

		if utils.PathSearch("login_policy.user_validity_period", getLoginPolicyRespBody, float64(0)).(float64) != 0 ||
			utils.PathSearch("login_policy.custom_info_for_login", getLoginPolicyRespBody, "").(string) != "" ||
			utils.PathSearch("login_policy.lockout_duration", getLoginPolicyRespBody, float64(0)).(float64) != 15 ||
			utils.PathSearch("login_policy.login_failed_times", getLoginPolicyRespBody, float64(0)).(float64) != 5 ||
			utils.PathSearch("login_policy.period_with_login_failures", getLoginPolicyRespBody, float64(0)).(float64) != 15 ||
			utils.PathSearch("login_policy.session_timeout", getLoginPolicyRespBody, float64(0)).(float64) != 60 ||
			utils.PathSearch("login_policy.show_recent_login_info", getLoginPolicyRespBody, false).(bool) ||
			utils.PathSearch("length(login_policy.allow_address_netmasks)", getLoginPolicyRespBody, float64(0)).(float64) != 0 ||
			utils.PathSearch("length(login_policy.allow_ip_ranges)", getLoginPolicyRespBody, float64(0)).(float64) != 1 ||
			utils.PathSearch("login_policy.allow_ip_ranges[0].ip_range", getLoginPolicyRespBody, "").(string) != "0.0.0.0-255.255.255.255" ||
			utils.PathSearch("login_policy.allow_ip_ranges[0].description", getLoginPolicyRespBody, "").(string) != "" ||
			utils.PathSearch("length(login_policy.allow_ip_ranges_ipv6)", getLoginPolicyRespBody, float64(0)).(float64) != 1 ||
			utils.PathSearch("login_policy.allow_ip_ranges_ipv6[0].ip_range", getLoginPolicyRespBody, "").(string) !=
				"0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF" ||
			utils.PathSearch("login_policy.allow_ip_ranges_ipv6[0].description", getLoginPolicyRespBody, "").(string) != "" {
			return errors.New("the login policy failed to reset to defaults")
		}
	}
	return nil
}

func testAccV5LoginPolicy_basic() string {
	return `
resource "huaweicloud_identityv5_login_policy" "test" {
  user_validity_period       = 20
  lockout_duration           = 30
  login_failed_times         = 10
  period_with_login_failures = 30
  session_timeout            = 120
  show_recent_login_info     = true
  custom_info_for_login      = "hello Terraform"
  allow_address_netmasks {
    address_netmask = "255.0.0.0/1"
    description     = "terraform test"
  }
  allow_ip_ranges {
    ip_range    = "0.0.0.0-255.255.255.254"
    description = "terraform test1"
  }
  allow_ip_ranges {
    ip_range    = "0.0.0.0-255.255.255.100"
    description = "terraform test2"
  }
  allow_ip_ranges_ipv6 {
    ip_range    = "0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFF0"
    description = "terraform test"
  }
}
`
}

func testAccV5LoginPolicy_update() string {
	return `
resource "huaweicloud_identityv5_login_policy" "test" {
  user_validity_period       = 0
  lockout_duration           = 20
  login_failed_times         = 6
  period_with_login_failures = 20
  session_timeout            = 90
  show_recent_login_info     = false
  custom_info_for_login      = ""
  allow_address_netmasks {
    address_netmask = "100.0.0.0/1"
    description     = "terraform test"
  }
  allow_ip_ranges {
    ip_range    = "0.0.0.0-255.255.255.100"
    description = "terraform test"
  }
  allow_ip_ranges_ipv6 {
    ip_range    = "0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFF1"
    description = "terraform test"
  }
}
`
}
