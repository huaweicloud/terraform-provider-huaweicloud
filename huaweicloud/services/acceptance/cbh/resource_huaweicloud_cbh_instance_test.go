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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCBHInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getInstance: Query CBH instance
	var (
		getInstanceHttpUrl = "v1/{project_id}/cbs/instance/list"
		getInstanceProduct = "cbh"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CBH Client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBHInstance: %s", err)
	}
	getCbhInstancesRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBHInstance: %s", err)
	}
	instances := flattenGetInstancesResponseBodyInstanceTest(getCbhInstancesRespBody)
	if instance, ok := instances[state.Primary.ID]; ok {
		return instance, nil
	}
	return nil, fmt.Errorf("error get CBH instance")
}

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
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"data.huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"data.huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip_id",
						"huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip",
						"huaweicloud_vpc_eip.test", "address"),
				),
			},
			{
				Config: testCBHInstance_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "flavor_id", "cbh.basic.10"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"data.huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"data.huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip_id",
						"huaweicloud_vpc_eip.test_update", "id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip",
						"huaweicloud_vpc_eip.test_update", "address"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"charging_mode", "password", "period", "period_unit", "auto_renew"},
			},
		},
	})
}

func flattenGetInstancesResponseBodyInstanceTest(resp interface{}) map[string]interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("instance", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make(map[string]interface{})
	for _, v := range curArray {
		rst[utils.PathSearch("instanceId", v, "").(string)] = v
	}
	return rst
}

func testCBHInstance_base() string {
	return `
data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`
}

func testCBHInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbh_instance" "test" {
  flavor_id         = "cbh.basic.10"
  name              = "%s"
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  public_ip_id      = huaweicloud_vpc_eip.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  region            = "%s"
  password          = "test_123456"
  charging_mode     = "prePaid"
  period_unit       = "month"
  auto_renew        = "false"
  period            = 1
}
`, testCBHInstance_base(), name, acceptance.HW_REGION_NAME)
}

func testCBHInstance_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test_update" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cbh_instance" "test" {
  flavor_id         = "cbh.basic.10"
  name              = "%s"
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  public_ip_id      = huaweicloud_vpc_eip.test_update.id
  public_ip         = huaweicloud_vpc_eip.test_update.address
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  region            = "%s"
  password          = "test_147258"
  charging_mode     = "prePaid"
  period_unit       = "month"
  auto_renew        = "true"
  period            = 1
}
`, testCBHInstance_base(), name, acceptance.HW_REGION_NAME)
}
