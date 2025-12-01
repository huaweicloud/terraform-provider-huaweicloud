package esw

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

func getEswInstance(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ESW client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func TestAccEswInstance_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_esw_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getEswInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEswInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_ref",
						"data.huaweicloud_esw_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "ha_mode", "ha"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0.primary",
						"data.huaweicloud_esw_flavors.test", "flavors.0.available_zones.0"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0.standby",
						"data.huaweicloud_esw_flavors.test", "flavors.0.available_zones.1"),
					resource.TestCheckResourceAttrPair(resourceName, "tunnel_info.0.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "tunnel_info.0.virsubnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_info.0.tunnel_ip", "192.168.0.192"),
					resource.TestCheckResourceAttr(resourceName, "charge_infos.0.charge_mode", "postPaid"),

					resource.TestCheckResourceAttrSet(resourceName, "tunnel_info.0.tunnel_port"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel_info.0.tunnel_type"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccEswInstance_basic_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_ref",
						"data.huaweicloud_esw_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "ha_mode", "ha"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0.primary",
						"data.huaweicloud_esw_flavors.test", "flavors.0.available_zones.0"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0.standby",
						"data.huaweicloud_esw_flavors.test", "flavors.0.available_zones.1"),
					resource.TestCheckResourceAttrPair(resourceName, "tunnel_info.0.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "tunnel_info.0.virsubnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_info.0.tunnel_ip", "192.168.0.192"),
					resource.TestCheckResourceAttr(resourceName, "charge_infos.0.charge_mode", "postPaid"),

					resource.TestCheckResourceAttrSet(resourceName, "tunnel_info.0.tunnel_port"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel_info.0.tunnel_type"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
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

func testAccEswInstance_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_esw_flavors" "test" {}

resource "huaweicloud_esw_instance" "test" {
  name        = "%[2]s"
  flavor_ref  = data.huaweicloud_esw_flavors.test.flavors[0].name
  ha_mode     = "ha"
  description = "test description"

  availability_zones {
    primary = data.huaweicloud_esw_flavors.test.flavors.0.available_zones[0]
    standby = data.huaweicloud_esw_flavors.test.flavors.0.available_zones[1]
  }

  tunnel_info {
    vpc_id       = huaweicloud_vpc.test.id
    virsubnet_id = huaweicloud_vpc_subnet.test.id
    tunnel_ip    = "192.168.0.192"
  }

  charge_infos {
    charge_mode = "postPaid"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccEswInstance_basic_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_esw_flavors" "test" {}

resource "huaweicloud_esw_instance" "test" {
  name        = "%[2]s"
  flavor_ref  = data.huaweicloud_esw_flavors.test.flavors.0.name
  ha_mode     = "ha"
  description = "test description update"

  availability_zones {
    primary = data.huaweicloud_esw_flavors.test.flavors.0.available_zones[0]
    standby = data.huaweicloud_esw_flavors.test.flavors.0.available_zones[1]
  }

  tunnel_info {
    vpc_id       = huaweicloud_vpc.test.id
    virsubnet_id = huaweicloud_vpc_subnet.test.id
    tunnel_ip    = "192.168.0.192"
  }

  charge_infos {
    charge_mode = "postPaid"
  }
}
`, common.TestVpc(rName), updateName)
}
