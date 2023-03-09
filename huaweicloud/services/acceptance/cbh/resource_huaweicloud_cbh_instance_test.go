package cbh

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCBHInstanceResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getInstance: Query CBH instance
	var (
		getInstanceHttpUrl = "v1/{project_id}/cbs/instance/list"
		getInstanceProduct = "cbh"
	)
	getInstanceClient, err := config.NewServiceClient(getInstanceProduct, region)
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
					resource.TestCheckResourceAttr(rName, "bastion_type", "OEM"),
				),
			},
			{
				Config: testCBHInstance_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "bastion_type", "OEM"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"availability_zone", "charging_mode", "cloud_service_type",
					"flavor_id", "hx_password", "password", "nics", "period", "period_unit", "product_info",
					"security_groups", "public_ip"},
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

func testCBHInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_cbh_instance" "test" {
  flavor_id          = "cbh.basic.50"
  name               = "%s"
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  region             = "%s"
  hx_password        = "test_123456"
  bastion_type       = "OEM"
  charging_mode      = "prePaid"
  period_unit        = "month"
  auto_renew         = "false"
  period             = "1"

  product_info {
    product_id         = "OFFI740586375358963717"
    resource_size      = "1"
  }
}
`, common.TestBaseNetwork(name), name, acceptance.HW_REGION_NAME)
}

func testCBHInstance_basic_update(name string) string {
	return fmt.Sprintf(`
%s

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

resource "huaweicloud_cbh_instance" "test" {
  flavor_id          = "cbh.basic.50"
  name               = "%s"
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id

  public_ip {
    id      = huaweicloud_vpc_eip_v1.test.id
    address = huaweicloud_vpc_eip_v1.test.address
  }

  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  region             = "%s"
  hx_password        = "test_123456"
  password           = "test_147258"
  bastion_type       = "OEM"
  charging_mode      = "prePaid"
  period_unit        = "month"
  auto_renew         = "false"
  period             = "1"
  
  product_info {
    product_id         = "OFFI740586375358963717"
    resource_size      = "1"
  }
}
`, common.TestBaseNetwork(name), name, acceptance.HW_REGION_NAME)
}
