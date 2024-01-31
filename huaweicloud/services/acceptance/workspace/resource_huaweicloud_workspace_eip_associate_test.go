package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/desktops"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getEipAssociateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace v2 client: %s", err)
	}
	resp, err := desktops.ListEips(client, state.Primary.Attributes["desktop_id"])
	if len(resp) > 0 || err != nil {
		return resp, err
	}
	return nil, golangsdk.ErrDefault404{}
}

func TestAccEipAssociate_basic(t *testing.T) {
	var (
		eips         []desktops.EipResp
		resourceName = "huaweicloud_workspace_eip_associate.test"
		name         = acceptance.RandomAccResourceNameWithDash()
		rc           = acceptance.InitResourceCheck(resourceName, &eips, getEipAssociateFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEipAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "huaweicloud_vpc_eip.test", "id"),
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

func testAccEipAssociate_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
  ]
}

data "huaweicloud_images_images" "test" {
  name_regex = "WORKSPACE"
  visibility = "market"
}

data "huaweicloud_workspace_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  os_type           = "Windows"
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id)
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id)
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  security_groups   = [
    huaweicloud_workspace_service.test.desktop_security_group.0.id,
    huaweicloud_networking_secgroup.test.id,
  ]

  nic {
    network_id = huaweicloud_vpc_subnet.test.id
  }

  name        = "%[2]s"
  user_name   = "user-%[2]s"
  user_email  = "terraform@example.com"
  user_group  = "administrators"
  delete_user = true

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccEipAssociate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 1
    share_type  = "PER"
    charge_mode = "bandwidth"
  }
}

resource "huaweicloud_workspace_eip_associate" "test" {
  desktop_id = huaweicloud_workspace_desktop.test.id
  eip_id     = huaweicloud_vpc_eip.test.id
}
`, testAccEipAssociate_base(name), name)
}
