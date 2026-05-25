package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getDevServerResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetDevServerById(client, state.Primary.ID)
}

func TestAccDevServer_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_modelarts_devserver.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDevServerResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelartsDevServer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDevServer_basic(name, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "resource_flavor", acceptance.HW_MODELARTS_DEVSERVER_FLAVOR),
					resource.TestCheckResourceAttrSet(resourceName, "architecture"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_MODELARTS_DEVSERVER_IMAGE_ID),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "PRE_PAID"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Deprecated parameters.
					resource.TestCheckResourceAttrSet(resourceName, "flavor"),
				),
			},
			{
				Config: testAccDevServer_basic(name, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			// Stopping the DevServer.
			{
				Config: testAccDevServer_doAction(name, "stop", false),
			},
			// Stopping the stopped DevServer.
			{
				Config:      testAccDevServer_doAction(name, "stop", true),
				ExpectError: regexp.MustCompile(`unable to stop DevServer \([a-f0-9-]+\)`),
			},
			// Reinstalling the DevServer OS.
			{
				Config: testAccDevServer_doAction(name, "reinstallos", false),
			},
			// After reinstalling the DevServer OS, the status be changed to "RUNNING", and if we want to test the start
			// action, we need to stop the DevServer first.
			{
				Config: testAccDevServer_doAction(name, "stop", false),
			},
			// Starting the DevServer (after reinstall OS image).
			{
				Config: testAccDevServer_doAction(name, "start", false),
			},
			// Starting the running DevServer.
			{
				Config:      testAccDevServer_doAction(name, "start", true),
				ExpectError: regexp.MustCompile(`unable to start DevServer \([a-f0-9-]+\)`),
			},
			// Rebooting the DevServer.
			{
				Config: testAccDevServer_doAction(name, "reboot", false),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"subnet_id",
					"security_group_id",
					"admin_pass",
					"root_volume",
					"period_unit",
					"period",
					"auto_renew",
				},
			},
		},
	})
}

func testAccDevServer_basic(name string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "random_password" "test" {
  length           = 12
  min_numeric      = 1
  min_upper        = 1
  min_lower        = 1
  special          = true
  min_special      = 1
  override_special = "!@"
}

resource "huaweicloud_modelarts_devserver" "test" {
  name              = "%[2]s"
  resource_flavor   = "%[3]s"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  image_id          = "%[4]s"
  admin_pass        = random_password.test.result

  root_volume {
    size = 100
    type = "SSD"
  }

  charging_mode = "PRE_PAID"
  period_unit   = "MONTH"
  period        = 1
  auto_renew    = "%[5]v"
}
`, common.TestBaseNetwork(name), name,
		acceptance.HW_MODELARTS_DEVSERVER_FLAVOR, acceptance.HW_MODELARTS_DEVSERVER_IMAGE_ID, isAutoRenew)
}

func testAccDevServer_doAction(name, actionType string, doRetryAction bool) string {
	return fmt.Sprintf(`
%[1]s

variable "action_type" {
  type    = string
  default = "%[2]s"
}

resource "huaweicloud_modelarts_devserver_action" "test" {
  devserver_id = huaweicloud_modelarts_devserver.test.id
  action       = var.action_type
  admin_pass   = var.action_type == "reinstallos" ? random_password.test.result : null

  enable_force_new = "true"
}

variable "is_retry_devserver_action" {
  type    = bool
  default = %[3]v
}

resource "huaweicloud_modelarts_devserver_action" "expect_err" {
  count = var.is_retry_devserver_action ? 1 : 0

  depends_on = [huaweicloud_modelarts_devserver_action.test]

  devserver_id = huaweicloud_modelarts_devserver.test.id
  action       = var.action_type

  enable_force_new = "true"
}
`, testAccDevServer_basic(name, false), actionType, doRetryAction)
}
