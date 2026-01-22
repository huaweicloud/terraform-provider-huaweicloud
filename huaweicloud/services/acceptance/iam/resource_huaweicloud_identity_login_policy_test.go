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

func getLoginPolicyResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	httpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/login-policy"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	if utils.PathSearch("login_policy.account_validity_period", respBody, float64(0)).(float64) == 0 &&
		utils.PathSearch("login_policy.custom_info_for_login", respBody, "").(string) == "" &&
		utils.PathSearch("login_policy.lockout_duration", respBody, float64(0)).(float64) == 15 &&
		utils.PathSearch("login_policy.login_failed_times", respBody, float64(0)).(float64) == 5 &&
		utils.PathSearch("login_policy.period_with_login_failures", respBody, float64(0)).(float64) == 15 &&
		utils.PathSearch("login_policy.session_timeout", respBody, float64(0)).(float64) == 60 &&
		!utils.PathSearch("login_policy.show_recent_login_info", respBody, false).(bool) {
		return nil, golangsdk.ErrDefault404{}
	}
	return respBody, nil
}

func TestAccLoginPolicy_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_login_policy.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getLoginPolicyResourceFunc)
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
				Config: testAccLoginPolicy_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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
				Config: testAccLoginPolicy_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

const testAccLoginPolicy_basic_step1 = `
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

const testAccLoginPolicy_basic_step2 = `
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
