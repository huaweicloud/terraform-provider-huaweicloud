package cbh

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

func getCBHInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region             = acceptance.HW_REGION_NAME
		getInstanceHttpUrl = "v2/{project_id}/cbs/instance/list"
		getInstanceProduct = "cbh"
	)
	client, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CBH client: %s", err)
	}

	getInstancePath := client.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", client.ProjectID)
	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getInstanceResp, err := client.Request("GET", getInstancePath, &getInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBH instance list: %s", err)
	}
	getCbhInstancesRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return nil, err
	}
	instances := utils.PathSearch("instance", getCbhInstancesRespBody, make([]interface{}, 0)).([]interface{})
	expression := fmt.Sprintf("[?server_id == '%s']|[0]", state.Primary.ID)
	instance := utils.PathSearch(expression, instances, nil)
	if instance == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return instance, nil
}

// Considering security vulnerabilities, the test cases do not include the binding and unbinding content of EIP.
func TestAccCBHInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cbh_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCBHInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCBHInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "flavor_id", "cbh.basic.10"),
					resource.TestCheckResourceAttr(rName, "password", "test_123456"),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(rName, "period_unit", "month"),
					resource.TestCheckResourceAttr(rName, "period", "1"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(rName, "subnet_address", "192.168.0.154"),
					resource.TestCheckResourceAttr(rName, "public_ip_id", ""),
					resource.TestCheckResourceAttr(rName, "public_ip", ""),

					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),

					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "version"),
				),
			},
			{
				Config: testCBHInstance_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "password", "test_147258"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"charging_mode",
					"password",
					"period",
					"period_unit",
					"auto_renew",
					"ipv6_enable",
				},
			},
		},
	})
}

func testCBHInstance_base(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}
`, common.TestBaseNetwork(name))
}

func testCBHInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbh_instance" "test" {
  flavor_id         = "cbh.basic.10"
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  subnet_address    = "192.168.0.154"
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  password          = "test_123456"
  charging_mode     = "prePaid"
  period_unit       = "month"
  auto_renew        = "false"
  period            = 1
}
`, testCBHInstance_base(name), name)
}

func testCBHInstance_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbh_instance" "test" {
  flavor_id         = "cbh.basic.10"
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  subnet_address    = "192.168.0.154"
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  password          = "test_147258"
  charging_mode     = "prePaid"
  period_unit       = "month"
  auto_renew        = "true"
  period            = 1
}
`, testCBHInstance_base(name), name)
}
