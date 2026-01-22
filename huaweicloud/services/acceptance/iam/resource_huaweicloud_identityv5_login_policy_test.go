package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getV5LoginPolicy(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v5/login-policy"
	getPath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	// If the login policy values are not the default values, means the login policy is not destroyed.
	if utils.PathSearch("login_policy.user_validity_period", respBody, float64(0)).(float64) != 0 ||
		utils.PathSearch("login_policy.custom_info_for_login", respBody, "").(string) != "" ||
		utils.PathSearch("login_policy.lockout_duration", respBody, float64(0)).(float64) != 15 ||
		utils.PathSearch("login_policy.login_failed_times", respBody, float64(0)).(float64) != 5 ||
		utils.PathSearch("login_policy.period_with_login_failures", respBody, float64(0)).(float64) != 15 ||
		utils.PathSearch("login_policy.session_timeout", respBody, float64(0)).(float64) != 60 ||
		utils.PathSearch("login_policy.show_recent_login_info", respBody, false).(bool) ||
		utils.PathSearch("length(login_policy.allow_address_netmasks)", respBody, float64(0)).(float64) != 0 ||
		utils.PathSearch("length(login_policy.allow_ip_ranges)", respBody, float64(0)).(float64) != 1 ||
		utils.PathSearch("login_policy.allow_ip_ranges[0].ip_range", respBody, "").(string) != "0.0.0.0-255.255.255.255" ||
		utils.PathSearch("login_policy.allow_ip_ranges[0].description", respBody, "").(string) != "" ||
		utils.PathSearch("length(login_policy.allow_ip_ranges_ipv6)", respBody, float64(0)).(float64) != 1 ||
		utils.PathSearch("login_policy.allow_ip_ranges_ipv6[0].ip_range", respBody, "").(string) !=
			"0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF" ||
		utils.PathSearch("login_policy.allow_ip_ranges_ipv6[0].description", respBody, "").(string) != "" {
		return respBody, nil
	}

	return nil, golangsdk.ErrDefault404{}
}

func getV5LoginPolicyFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return getV5LoginPolicy(client)
}

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5LoginPolicy_basic(t *testing.T) {
	var (
		rName = "huaweicloud_identityv5_login_policy.test"
		obj   interface{}
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5LoginPolicyFunc)
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
				Config: testAccV5LoginPolicy_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "user_validity_period", "20"),
					resource.TestCheckResourceAttr(rName, "lockout_duration", "30"),
					resource.TestCheckResourceAttr(rName, "login_failed_times", "10"),
					resource.TestCheckResourceAttr(rName, "period_with_login_failures", "30"),
					resource.TestCheckResourceAttr(rName, "session_timeout", "120"),
					resource.TestCheckResourceAttr(rName, "show_recent_login_info", "true"),
					resource.TestCheckResourceAttr(rName, "custom_info_for_login", "hello Terraform"),
					resource.TestCheckResourceAttr(rName, "allow_address_netmasks.#", "1"),
					resource.TestCheckResourceAttr(rName, "allow_address_netmasks.0.address_netmask", "255.0.0.0/1"),
					resource.TestCheckResourceAttr(rName, "allow_address_netmasks.0.description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges.#", "2"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges.0.ip_range", "0.0.0.0-255.255.255.254"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges.0.description", ""),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges.1.ip_range", "0.0.0.0-255.255.255.100"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges.1.description", "terraform test2"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges_ipv6.#", "1"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges_ipv6.0.ip_range",
						"0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFF0"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges_ipv6.0.description", "terraform test"),
				),
			},
			{
				Config: testAccV5LoginPolicy_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "user_validity_period", "0"),
					resource.TestCheckResourceAttr(rName, "lockout_duration", "20"),
					resource.TestCheckResourceAttr(rName, "login_failed_times", "6"),
					resource.TestCheckResourceAttr(rName, "period_with_login_failures", "20"),
					resource.TestCheckResourceAttr(rName, "session_timeout", "90"),
					resource.TestCheckResourceAttr(rName, "show_recent_login_info", "false"),
					resource.TestCheckResourceAttr(rName, "custom_info_for_login", ""),
					resource.TestCheckResourceAttr(rName, "allow_address_netmasks.#", "1"),
					resource.TestCheckResourceAttr(rName, "allow_address_netmasks.0.address_netmask", "100.0.0.0/1"),
					resource.TestCheckResourceAttr(rName, "allow_address_netmasks.0.description", ""),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges.0.ip_range", "0.0.0.0-255.255.255.100"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges.0.description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges_ipv6.#", "1"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges_ipv6.0.ip_range",
						"0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFF1"),
					resource.TestCheckResourceAttr(rName, "allow_ip_ranges_ipv6.0.description", ""),
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

const testAccV5LoginPolicy_basic_step1 = `
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

const testAccV5LoginPolicy_basic_step2 = `
resource "huaweicloud_identityv5_login_policy" "test" {
  user_validity_period       = 0
  lockout_duration           = 20
  login_failed_times         = 6
  period_with_login_failures = 20
  session_timeout            = 90
  show_recent_login_info     = false

  allow_address_netmasks {
    address_netmask = "100.0.0.0/1"
  }

  allow_ip_ranges {
    ip_range    = "0.0.0.0-255.255.255.100"
	description = "terraform test"
  }

  allow_ip_ranges_ipv6 {
    ip_range = "0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFF1"
  }
}
`
