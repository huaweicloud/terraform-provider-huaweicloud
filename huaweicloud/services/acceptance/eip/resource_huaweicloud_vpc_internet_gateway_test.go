package eip

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

func getIGWResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NetworkingV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC V3 client: %s", err)
	}

	getIGWHttpUrl := "v3/{project_id}/geip/vpc-igws/{vpc_igw_id}"
	getIGWPath := client.Endpoint + getIGWHttpUrl
	getIGWPath = strings.ReplaceAll(getIGWPath, "{project_id}", client.ProjectID)
	getIGWPath = strings.ReplaceAll(getIGWPath, "{vpc_igw_id}", state.Primary.ID)
	getIGWOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getIGWResp, err := client.Request("GET", getIGWPath, &getIGWOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VPC internet gateway: %s", err)
	}
	return utils.FlattenResponse(getIGWResp)
}

func TestAccIGW_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpc_internet_gateway.test"
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIGWResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIGW_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enable_ipv6", "false"),
					resource.TestCheckResourceAttr(rName, "add_route", "true"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccIGW_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
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

func testVPC(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}
`, name, name)
}

func testAccIGW_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_internet_gateway" "test" {
  depends_on = [huaweicloud_vpc_subnet.test]

  vpc_id      = huaweicloud_vpc.test.id
  name        = "%s"
  add_route   = true
  enable_ipv6 = false
}
`, testVPC(name), name)
}

func testAccIGW_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_internet_gateway" "test" {
  depends_on = [huaweicloud_vpc_subnet.test]

  vpc_id      = huaweicloud_vpc.test.id
  name        = "%s-update"
  add_route   = true
  enable_ipv6 = true
}
`, testVPC(name), name)
}
