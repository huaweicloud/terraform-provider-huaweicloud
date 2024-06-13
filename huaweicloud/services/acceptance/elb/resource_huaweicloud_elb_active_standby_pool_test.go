package elb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getELBActiveStandbyPoolResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		getActiveStandbyPoolUrl     = "v3/{project_id}/elb/master-slave-pools/{pool_id}"
		getActiveStandbyPoolProduct = "elb"
	)
	elbClient, err := cfg.NewServiceClient(getActiveStandbyPoolProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, err
	}
	getActiveStandbyPoolPath := elbClient.Endpoint + getActiveStandbyPoolUrl
	getActiveStandbyPoolPath = strings.ReplaceAll(getActiveStandbyPoolPath, "{project_id}", elbClient.ProjectID)
	getActiveStandbyPoolPath = strings.ReplaceAll(getActiveStandbyPoolPath, "{pool_id}", state.Primary.ID)

	getActiveStandbyPoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getActiveStandbyPoolResp, err := elbClient.Request("GET", getActiveStandbyPoolPath, &getActiveStandbyPoolOpt)
	if err != nil {
		return nil, err
	}
	getActiveStandbyPoolBody, err := utils.FlattenResponse(getActiveStandbyPoolResp)
	if err != nil {
		return nil, err
	}
	return getActiveStandbyPoolBody, nil
}

func TestAccElbActiveStandbyPool_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_active_standby_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getELBActiveStandbyPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbActiveStandbyPoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "type", "instance"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "any_port_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "dualstack"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "healthmonitor.0.delay", "5"),
					resource.TestCheckResourceAttr(resourceName, "healthmonitor.0.expected_codes", "200"),
					resource.TestCheckResourceAttr(resourceName, "healthmonitor.0.http_method", "HEAD"),
					resource.TestCheckResourceAttr(resourceName, "healthmonitor.0.max_retries", "3"),
					resource.TestCheckResourceAttr(resourceName, "healthmonitor.0.max_retries_down", "3"),
					resource.TestCheckResourceAttr(resourceName, "healthmonitor.0.timeout", "3"),
					resource.TestCheckResourceAttr(resourceName, "healthmonitor.0.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "connection_drain_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "connection_drain_timeout", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "members.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "members.0.member_type"),
					resource.TestCheckResourceAttrSet(resourceName, "members.0.operating_status"),
					resource.TestCheckResourceAttrSet(resourceName, "members.0.ip_version"),
					resource.TestCheckResourceAttrSet(resourceName, "healthmonitor.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
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

func testAccElbActiveStandbyPoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_active_standby_pool" "test" {
  name                     = "%s"
  description              = "test"
  protocol                 = "TCP"
  vpc_id                   = huaweicloud_vpc.test.id
  type                     = "instance"
  any_port_enable          = false
  ip_version               = "dualstack"
  connection_drain_enabled = true
  connection_drain_timeout = 100

  members {
    address       = "192.168.0.1"
    role          = "master"
    protocol_port = 45
  }

  members {
    address       = "192.168.0.2"
    role          = "slave"
    protocol_port = 36
  }

  healthmonitor {
    delay            = 5
    expected_codes   = "200"
    http_method      = "HEAD"
    max_retries      = 3
    max_retries_down = 3
    timeout          = 3
    type             = "HTTP"
  }

}
`, common.TestVpc(rName), rName)
}
