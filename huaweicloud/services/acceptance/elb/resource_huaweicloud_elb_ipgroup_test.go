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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getELBIpGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/ipgroups/{ipgroup_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, err
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{ipgroup_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func TestAccElbV3IpGroup_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_ipgroup.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getELBIpGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3IpGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.0.ip", "192.168.10.10"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.0.description", "ECS01"),
				),
			},
			{
				Config: testAccElbV3IpGroupConfig_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test updated"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
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

func testAccElbV3IpGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_ipgroup" "test"{
  name                  = "%s"
  description           = "terraform test"
  enterprise_project_id = "0"

  ip_list {
    ip          = "192.168.10.10"
    description = "ECS01"
  }
}
`, name)
}

func testAccElbV3IpGroupConfig_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_ipgroup" "test"{
  name                  = "%s"
  description           = "terraform test updated"
  enterprise_project_id = "0"

  ip_list {
    ip          = "192.168.10.10"
    description = "ECS01"
  }

  ip_list {
    ip          = "192.168.10.11"
    description = "ECS02"
  }
}
`, name)
}
