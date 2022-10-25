package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/workspace/v2/desktops"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDesktopFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating Workspace v2 client: %s", err)
	}
	return desktops.Get(client, state.Primary.ID)
}

func TestAccDesktop_basic(t *testing.T) {
	var (
		desktop      desktops.Desktop
		resourceName = "huaweicloud_workspace_desktop.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&desktop,
		getDesktopFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDesktop_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "workspace.x86.ultimate.large2"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "user_name", "user-"+rName),
					resource.TestCheckResourceAttr(resourceName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.size", "70"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccDesktop_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "workspace.x86.ultimate.large4"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.size", "90"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.2.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.2.size", "20"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.3.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.3.size", "40"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_user",
					"image_type",
					"nic",
					"user_email",
					"vpc_id",
				},
			},
		},
	})
}

func testAccDesktop_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.128.0/18"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
  ]
}
`, rName)
}

func testAccDesktop_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  data_volume_sizes = [50, 70]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = "workspace.x86.ultimate.large2"
  image_type        = "market"
  image_id          = "63aa8670-27ad-4747-8c44-6d8919e785a7"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  security_groups   = [
    huaweicloud_workspace_service.test.desktop_security_group.0.id,
    huaweicloud_networking_secgroup.test.id,
  ]

  nic {
    network_id = huaweicloud_vpc_subnet.test.id
  }

  name       = "%[2]s"
  user_name  = "user-%[2]s"
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 80
  }

  dynamic "data_volume" {
    for_each = local.data_volume_sizes

    content {
      type = "SAS"
      size = data_volume.value
    }
  }

  tags = {
    foo = "bar"
  }

  delete_user = true
}
`, testAccDesktop_base(rName), rName)
}

func testAccDesktop_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  data_volume_sizes = [50, 90, 20, 40]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = "workspace.x86.ultimate.large4"
  image_type        = "market"
  image_id          = "63aa8670-27ad-4747-8c44-6d8919e785a7"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  security_groups   = [
    huaweicloud_workspace_service.test.desktop_security_group.0.id,
    huaweicloud_networking_secgroup.test.id,
  ]

  nic {
    network_id = huaweicloud_vpc_subnet.test.id
  }

  name       = "%[2]s"
  user_name  = "user-%[2]s"
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 100
  }

  dynamic "data_volume" {
    for_each = local.data_volume_sizes

    content {
      type = "SAS"
      size = data_volume.value
    }
  }

  tags = {
    foo = "bar"
  }

  delete_user = true
}
`, testAccDesktop_base(rName), rName)
}
