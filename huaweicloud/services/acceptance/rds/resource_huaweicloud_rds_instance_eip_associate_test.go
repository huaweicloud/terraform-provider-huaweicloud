package rds

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

func getRdsInstanceEipAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances?id={instance_id}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	publicIP := utils.PathSearch("instances|[0].public_ips[0]", getRespBody, nil)
	if publicIP == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccRdsInstanceEipAssociate_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_instance_eip_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRdsInstanceEipAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceEipAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip",
						"huaweicloud_vpc_eip.test", "address"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip_id",
						"huaweicloud_vpc_eip.test", "id"),
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

func testAccRdsInstanceEipAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]

  db {
    type    = "MySQL"
    version = "8.0"
    port    = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    share_type  = "PER"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_rds_instance_eip_associate" "test" { 
  instance_id  = huaweicloud_rds_instance.test.id
  public_ip    = huaweicloud_vpc_eip.test.address
  public_ip_id = huaweicloud_vpc_eip.test.id
}`, common.TestBaseNetwork(rName), rName)
}
