package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getDesktopFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}
	return workspace.GetDesktopById(client, state.Primary.ID)
}

func TestAccDesktop_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_workspace_desktop.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDesktopFunc)

		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDesktop_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(rName, "image_id", acceptance.HW_WORKSPACE_DESKTOP_IMAGE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "user_name", "user-"+name),
					resource.TestCheckResourceAttr(rName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(rName, "root_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(rName, "root_volume.0.size", "80"),
					resource.TestCheckResourceAttr(rName, "data_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(rName, "data_volume.0.size", "50"),
					resource.TestCheckResourceAttr(rName, "data_volume.1.type", "SAS"),
					resource.TestCheckResourceAttr(rName, "data_volume.1.size", "70"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(rName, "nic.0.network_id", "data.huaweicloud_workspace_service.test", "network_ids.0"),
					resource.TestCheckResourceAttr(rName, "email_notification", "true"),
				),
			},
			{
				Config: testAccDesktop_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(rName, "root_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(rName, "root_volume.0.size", "100"),
					resource.TestCheckResourceAttr(rName, "data_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(rName, "data_volume.0.size", "50"),
					resource.TestCheckResourceAttr(rName, "data_volume.1.type", "SAS"),
					resource.TestCheckResourceAttr(rName, "data_volume.1.size", "90"),
					resource.TestCheckResourceAttr(rName, "data_volume.2.type", "SAS"),
					resource.TestCheckResourceAttr(rName, "data_volume.2.size", "20"),
					resource.TestCheckResourceAttr(rName, "data_volume.3.type", "SAS"),
					resource.TestCheckResourceAttr(rName, "data_volume.3.size", "40"),
					resource.TestCheckResourceAttr(rName, "tags.%", "0"),
					resource.TestCheckResourceAttrPair(rName, "nic.0.network_id", "data.huaweicloud_workspace_service.test", "network_ids.1"),
				),
			},
			{
				ResourceName:      rName,
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

func testAccDesktop_basic_step1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  os_type           = "Windows"
}

locals {
  cpu_flavor_ids    = [for v in data.huaweicloud_workspace_flavors.test.flavors : v.id if !v.is_gpu]
  data_volume_sizes = [50, 70]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(local.cpu_flavor_ids[0], "NOT_FOUND")
  image_type        = "market"
  image_id          = "%[1]s"
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = data.huaweicloud_workspace_service.test.desktop_security_group[*].id

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
  }

  name               = "%[2]s"
  user_name          = "user-%[2]s"
  user_group         = "administrators"
  user_email         = "terraform@example.com"
  delete_user        = true
  email_notification = true

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

  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zone,
      user_name,
    ]
  }
}
`, acceptance.HW_WORKSPACE_DESKTOP_IMAGE_ID, name)
}

func testAccDesktop_basic_step2(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  os_type           = "Windows"
}

locals {
  cpu_flavor_ids    = [for v in data.huaweicloud_workspace_flavors.test.flavors : v.id if !v.is_gpu]
  data_volume_sizes = [50, 90, 20, 40]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(local.cpu_flavor_ids[0], "NOT_FOUND")
  image_type        = "market"
  image_id          = "%[1]s"
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = data.huaweicloud_workspace_service.test.desktop_security_group[*].id

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[1], "NOT_FOUND")
  }

  name               = "%[2]s"
  user_name          = "user-%[2]s"
  user_group         = "administrators"
  user_email         = "terraform@example.com"
  delete_user        = true
  email_notification = true

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


  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zone,
      user_name,
    ]
  }
}
`, acceptance.HW_WORKSPACE_DESKTOP_IMAGE_ID, name)
}

func TestAccDesktop_withEpsId(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_workspace_desktop.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDesktopFunc)

		name = acceptance.RandomAccResourceNameWithDash()

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
				Config: testAccDesktop_withEPSId(srcEPS, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "user_name", "user-"+name),
					resource.TestCheckResourceAttr(rName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccDesktop_withEPSId(destEPS, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "user_name", "user-"+name),
					resource.TestCheckResourceAttr(rName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testAccDesktop_withEPSId(epsId, name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  os_type           = "Windows"
}

locals {
  cpu_flavor_ids = [for v in data.huaweicloud_workspace_flavors.test.flavors : v.id if !v.is_gpu] 
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id             = try(local.cpu_flavor_ids[0], "NOT_FOUND")
  image_type            = "market"
  image_id              = "%[1]s"
  enterprise_project_id = "%[2]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = data.huaweicloud_workspace_service.test.vpc_id
  security_groups       = data.huaweicloud_workspace_service.test.desktop_security_group[*].id

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
  }

  name        = "%[3]s"
  user_name   = "user-%[3]s"
  user_group  = "administrators"
  user_email  = "terraform@example.com"
  delete_user = true

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }

  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zone,
      user_name,
    ]
  }
}
`, acceptance.HW_WORKSPACE_DESKTOP_IMAGE_ID, epsId, name)
}

func TestAccDesktop_powerAction(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_workspace_desktop.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDesktopFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccdesktop_powerAction(name, "os-stop", "SOFT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "power_action", "os-stop"),
					resource.TestCheckResourceAttr(rName, "power_action_type", "SOFT"),
					resource.TestCheckResourceAttr(rName, "status", "SHUTOFF"),
				),
			},
			{
				Config: testAccdesktop_powerAction(name, "os-start", "SOFT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "power_action", "os-start"),
					resource.TestCheckResourceAttr(rName, "power_action_type", "SOFT"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccdesktop_powerAction(name, "reboot", "SOFT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "power_action", "reboot"),
					resource.TestCheckResourceAttr(rName, "power_action_type", "SOFT"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccdesktop_powerAction(name, "reboot", "HARD"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "power_action", "reboot"),
					resource.TestCheckResourceAttr(rName, "power_action_type", "HARD"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccdesktop_powerAction(name, "os-hibernate", "HARD"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "power_action", "os-hibernate"),
					resource.TestCheckResourceAttr(rName, "power_action_type", "HARD"),
					resource.TestCheckResourceAttr(rName, "status", "HIBERNATED"),
				),
			},
			{
				ResourceName:      rName,
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

func testAccdesktop_powerAction(name string, powerAction string, powerActionType string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  os_type           = "Windows"
}

locals {
  cpu_flavor_ids = [for v in data.huaweicloud_workspace_flavors.test.flavors : v.id if !v.is_gpu] 
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(local.cpu_flavor_ids[0], "NOT_FOUND")
  image_type        = "market"
  image_id          = "%[1]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = data.huaweicloud_workspace_service.test.desktop_security_group[*].id

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
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

  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zone,
      user_name,
    ]
  }
}
`, acceptance.HW_WORKSPACE_DESKTOP_IMAGE_ID, name, powerAction, powerActionType)
}
