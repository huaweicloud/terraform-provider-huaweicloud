package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getDesktopFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating Workspace v2 client: %s", err)
	}
	return workspace.GetDesktopById(client, state.Primary.ID)
}

func TestAccDesktop_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceNameWithDash()

		desktop      interface{}
		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&desktop,
			getDesktopFunc,
		)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDesktop_basic_step1(rName),
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
					resource.TestCheckResourceAttrPair(resourceName, "nic.0.network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "email_notification", "true"),
				),
			},
			{
				Config: testAccDesktop_basic_step2(rName),
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
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.size", "70"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(resourceName, "nic.0.network_id", "huaweicloud_vpc_subnet.test", "id"),
				),
			},
			{
				Config: testAccDesktop_basic_step3(rName),
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
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "nic.0.network_id", "huaweicloud_vpc_subnet.standby", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_user",
					"image_type",
					"user_email",
					"vpc_id",
					"email_notification",
				},
			},
		},
	})
}

func testAccDesktop_basic_step1(rName string) string {
	return testAccDesktop_basic(rName, 80)
}

func testAccDesktop_basic_step2(rName string) string {
	return testAccDesktop_basic(rName, 100)
}

func testAccDesktop_basic_step3(rName string) string {
	return testAccDesktop_basic_update(rName)
}

func TestAccDesktop_UpdateWithEpsId(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceNameWithDash()

		desktop      interface{}
		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&desktop,
			getDesktopFunc,
		)

		srcEPS  = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		destEPS = acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDesktop_withEPSId(rName, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "workspace.x86.ultimate.large2"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "user_name", "user-"+rName),
					resource.TestCheckResourceAttr(resourceName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccDesktop_withEPSId(rName, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "workspace.x86.ultimate.large2"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "user_name", "user-"+rName),
					resource.TestCheckResourceAttr(resourceName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func TestAccDesktop_powerAction(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceNameWithDash()

		desktop      interface{}
		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&desktop,
			getDesktopFunc,
		)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccdesktop_powerAction(rName, "os-stop", "SOFT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "os-stop"),
					resource.TestCheckResourceAttr(resourceName, "power_action_type", "SOFT"),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
			{
				Config: testAccdesktop_powerAction(rName, "os-start", "SOFT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "os-start"),
					resource.TestCheckResourceAttr(resourceName, "power_action_type", "SOFT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccdesktop_powerAction(rName, "reboot", "SOFT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "reboot"),
					resource.TestCheckResourceAttr(resourceName, "power_action_type", "SOFT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccdesktop_powerAction(rName, "reboot", "HARD"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "reboot"),
					resource.TestCheckResourceAttr(resourceName, "power_action_type", "HARD"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccdesktop_powerAction(rName, "os-hibernate", "HARD"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "os-hibernate"),
					resource.TestCheckResourceAttr(resourceName, "power_action_type", "HARD"),
					resource.TestCheckResourceAttr(resourceName, "status", "HIBERNATED"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_user",
					"image_type",
					"user_email",
					"vpc_id",
					"power_action",
					"power_action_type",
				},
			},
		},
	})
}

func testAccDesktop_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_subnet" "standby" {
  name       = "%[2]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
    huaweicloud_vpc_subnet.standby.id,
  ]
}

data "huaweicloud_images_images" "test" {
  name_regex = "WORKSPACE"
  visibility = "market"
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccDesktop_basic(rName string, rootVolumeSize int) string {
	return fmt.Sprintf(`
%[1]s

locals {
  data_volume_sizes = [50, 70]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = "workspace.x86.ultimate.large2"
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, "")
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
    size = %[3]d
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

  email_notification = true
  delete_user        = true
}
`, testAccDesktop_base(rName), rName, rootVolumeSize)
}

func testAccDesktop_basic_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  data_volume_sizes = [50, 90, 20, 40]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = "workspace.x86.ultimate.large4"
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[1].id, "")
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  security_groups   = [
    huaweicloud_workspace_service.test.desktop_security_group.0.id,
    huaweicloud_networking_secgroup.test.id,
  ]

  nic {
    network_id = huaweicloud_vpc_subnet.standby.id
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

  email_notification = true
  delete_user        = true
}
`, testAccDesktop_base(rName), rName)
}

func testAccDesktop_withEPSId(rName, epsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id             = "workspace.x86.ultimate.large2"
  image_type            = "market"
  image_id              = try(data.huaweicloud_images_images.test.images[0].id, "")
  enterprise_project_id = "%[3]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = huaweicloud_vpc.test.id
  security_groups       = [
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

  data_volume {
    type = "SAS"
    size = 50
  }

  tags = {
    foo = "bar"
  }

  delete_user = true
}
`, testAccDesktop_base(rName), rName, epsId)
}

func testAccdesktop_powerAction(rName string, powerAction string, powerActionType string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = "workspace.x86.ultimate.large2"
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
  
  name               = "%[2]s"
  user_name          = "user-%[2]s"
  user_email         = "terraform@example.com"
  user_group         = "administrators"
  delete_user        = true
  power_action       = "%[3]s"
  power_action_type  = "%[4]s"

  root_volume {
    type = "SAS"
    size = 80
  }
  
  data_volume {
    type = "SAS"
    size = 50
  }
}
`, testAccDesktop_base(rName), rName, powerAction, powerActionType)
}
