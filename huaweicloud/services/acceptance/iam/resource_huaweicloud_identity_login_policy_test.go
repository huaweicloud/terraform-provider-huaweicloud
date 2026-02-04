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

func getV3LoginPolicyResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV3LoginPolicy(client, cfg.DomainID)
}

func TestAccV3LoginPolicy_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_login_policy.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3LoginPolicyResourceFunc)
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
				Config: testAccV3LoginPolicy_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "account_validity_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "lockout_duration", "30"),
					resource.TestCheckResourceAttr(resourceName, "login_failed_times", "10"),
					resource.TestCheckResourceAttr(resourceName, "period_with_login_failures", "30"),
					resource.TestCheckResourceAttr(resourceName, "session_timeout", "120"),
					resource.TestCheckResourceAttr(resourceName, "show_recent_login_info", "true"),
					resource.TestCheckResourceAttr(resourceName, "custom_info_for_login", "Hello Terraform"),
				),
			},
			{
				Config: testAccV3LoginPolicy_basic_step2,
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

const testAccV3LoginPolicy_basic_step1 = `
resource "huaweicloud_identity_login_policy" "test" {
  account_validity_period    = 20
  lockout_duration           = 30
  login_failed_times         = 10
  period_with_login_failures = 30
  session_timeout            = 120
  show_recent_login_info     = true
  custom_info_for_login      = "Hello Terraform"
}
`

const testAccV3LoginPolicy_basic_step2 = `
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
