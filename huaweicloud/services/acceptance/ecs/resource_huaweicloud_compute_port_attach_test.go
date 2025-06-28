package ecs

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

func getPortAttachResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("ecs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating compute client: %s", err)
	}

	url := client.Endpoint + "v2.1/{project_id}/servers/{server_id}/os-interface/{port_id}"
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)
	url = strings.ReplaceAll(url, "{server_id}", state.Primary.Attributes["instance_id"])
	url = strings.ReplaceAll(url, "{port_id}", state.Primary.Attributes["port_id"])
	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	response, err := client.Request("GET", url, &reqOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ECS Port: %s", err)
	}
	body, err := utils.FlattenResponse(response)
	if err != nil {
		return nil, fmt.Errorf("error prasing ECS Port: %s", err)
	}

	port := utils.PathSearch("interfaceAttachment.port_id", body, nil)
	if port == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return port, nil
}

func TestAccComputePortAttach_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_compute_port_attach.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPortAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputePortAttachBase(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "port_state", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "port_id", "huaweicloud_vpc_network_interface.port", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccComputePortAttachImport(resourceName),
			},
		},
	})
}

func testAccComputePortAttachImport(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["port_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<port_id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["port_id"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["port_id"]), nil
	}
}

func testAccComputePortAttachBase(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  flavor_id = data.huaweicloud_compute_flavors.test.ids[0]

  os         = "Ubuntu"
  visibility = "public"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_images.test.images[0].id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_vpc_network_interface" "port" {
  name               = "%[1]s_port"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
}

resource "huaweicloud_compute_port_attach" "test" {
  instance_id = huaweicloud_compute_instance.test.id
  port_id     = huaweicloud_vpc_network_interface.port.id
}

`, rName)
}
