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
		obj          interface{}
		resourceName = "huaweicloud_modelarts_devserver.test"
		name         = acceptance.RandomAccResourceName()
		password     = acceptance.RandomPassword("!@%-_=+[{}]:,./?")
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&obj,
			getDevServerResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelartsDevServer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDevServer_basic(true, name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor", acceptance.HW_MODELARTS_DEVSERVER_FLAVOR),
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
				),
			},
			{
				Config: testAccDevServer_basic(false, name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			// Stopping the DevServer.
			{
				Config: testAccDevServer_doAction(name, password, "stop", false),
			},
			// Stopping the stopped DevServer.
			{
				Config:      testAccDevServer_doAction(name, password, "stop", true),
				ExpectError: regexp.MustCompile(`Resource.Server '[a-f0-9-]+' is not allowed STOP: STOPPED`),
			},
			// Starting the DevServer.
			{
				Config: testAccDevServer_doAction(name, password, "start", false),
			},
			// Starting the running DevServer.
			{
				Config:      testAccDevServer_doAction(name, password, "start", true),
				ExpectError: regexp.MustCompile(`Resource.Server '[a-f0-9-]+' is not allowed START: RUNNING`),
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

func testAccDevServer_basic(isAutoRenew bool, name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_devserver" "test" {
  name              = "%[2]s"
  flavor            = "%[3]s"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  image_id          = "%[4]s"
  admin_pass        = "%[5]s"

  root_volume {
    size = 100
    type = "SSD"
  }

  charging_mode = "PRE_PAID"
  period_unit   = "MONTH"
  period        = 1
  auto_renew    = "%[6]v"
}
`, common.TestBaseNetwork(name), name,
		acceptance.HW_MODELARTS_DEVSERVER_FLAVOR, acceptance.HW_MODELARTS_DEVSERVER_IMAGE_ID, password, isAutoRenew)
}

func testAccDevServer_doAction(name, password, actionType string, doRetryAction bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_devserver_action" "test" {
  devserver_id = huaweicloud_modelarts_devserver.test.id
  action       = "%[2]s"
}

variable "is_retry_devserver_action" {
  type    = bool
  default = "%[3]v"
}

resource "huaweicloud_modelarts_devserver_action" "expect_err" {
  count = var.is_retry_devserver_action ? 1 : 0

  depends_on = [huaweicloud_modelarts_devserver_action.test]

  devserver_id = huaweicloud_modelarts_devserver.test.id
  action       = "%[2]s"
}
`, testAccDevServer_basic(false, name, password), actionType, doRetryAction)
}
