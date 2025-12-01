package workspace

import (
	"fmt"
	"regexp"
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

// Before running this test, make sure the Workspace service have at least two networks.
func TestAccDesktop_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDesktopFunc)

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
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_WORKSPACE_DESKTOP_IMAGE_ID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "user_name", "user-"+name),
					resource.TestCheckResourceAttr(resourceName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.size", "70"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(resourceName, "nic.0.network_id", "data.huaweicloud_workspace_service.test", "network_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "email_notification", "true"),
				),
			},
			{
				Config: testAccDesktop_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "user_group", "administrators"),
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
					resource.TestCheckResourceAttrPair(resourceName, "nic.0.network_id", "data.huaweicloud_workspace_service.test", "network_ids.1"),
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
		obj interface{}

		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDesktopFunc)

		name = acceptance.RandomAccResourceNameWithDash()

		srcEPS  = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		destEPS = acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopImageId(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDesktop_withEPSId(srcEPS, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "user_name", "user-"+name),
					resource.TestCheckResourceAttr(resourceName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccDesktop_withEPSId(destEPS, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "user_name", "user-"+name),
					resource.TestCheckResourceAttr(resourceName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
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
		obj interface{}

		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDesktopFunc)

		name = acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccdesktop_powerAction(name, "os-stop", "SOFT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "power_action", "os-stop"),
					resource.TestCheckResourceAttr(resourceName, "power_action_type", "SOFT"),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
			{
				Config: testAccdesktop_powerAction(name, "os-start", "SOFT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "power_action", "os-start"),
					resource.TestCheckResourceAttr(resourceName, "power_action_type", "SOFT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccdesktop_powerAction(name, "reboot", "SOFT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "power_action", "reboot"),
					resource.TestCheckResourceAttr(resourceName, "power_action_type", "SOFT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccdesktop_powerAction(name, "reboot", "HARD"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "power_action", "reboot"),
					resource.TestCheckResourceAttr(resourceName, "power_action_type", "HARD"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccdesktop_powerAction(name, "os-hibernate", "HARD"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
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

func TestAccDesktop_rootVolumeModify(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDesktopFunc)

		name = acceptance.RandomAccResourceNameWithDash()
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
				// Test for create GPSSD2 root volume.
				Config: testAccDesktop_rootVolumeModify_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "root_volume.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "120"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.throughput", "125"),
				),
			},
			{
				// Test for modify the size of GPSSD2 root volume.
				Config: testAccDesktop_rootVolumeModify_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "root_volume.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "160"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.throughput", "125"),
				),
			},
			{
				// Test for modify the QoS of GPSSD2 root volume.
				Config: testAccDesktop_rootVolumeModify_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "root_volume.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "160"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.throughput", "125"),
				),
			},
		},
	})
}

func testAccDesktop_volumeChange_resource(varBlocks, name string) string {
	return fmt.Sprintf(`
%[1]s

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
  image_id          = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = data.huaweicloud_workspace_service.test.desktop_security_group[*].id

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
  }

  name       = "%[3]s"
  user_name  = "user-%[3]s"
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type       = var.root_volume.type
    size       = var.root_volume.size
    iops       = try(var.root_volume.iops, null)
    throughput = try(var.root_volume.throughput, null)
  }

  dynamic "data_volume" {
    for_each = var.data_volumes

    content {
      type       = data_volume.value.type
      size       = data_volume.value.size
      iops       = try(data_volume.value.iops, null)
      throughput = try(data_volume.value.throughput, null)
    }
  }

  email_notification = true
  delete_user        = true

  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zone,
      user_name,
    ]
  }
}
`, varBlocks, acceptance.HW_WORKSPACE_DESKTOP_IMAGE_ID, name)
}

func testAccDesktop_rootVolumeModify(name string, size, iops, throughput int) string {
	varBlocks := fmt.Sprintf(`
variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number, null)
    throughput = optional(number, null)
  })

  default = {
    type       = "GPSSD2"
    size       = %[1]d
    iops       = %[2]d
    throughput = %[3]d
  }
}

variable "data_volumes" {
  type = list(object({
    type = string
    size = number
  }))

  default = [
    {
      type = "SSD"
      size = 80
    }
  ]
}
`, size, iops, throughput)
	return testAccDesktop_volumeChange_resource(varBlocks, name)
}

func testAccDesktop_rootVolumeModify_step1(name string) string {
	return testAccDesktop_rootVolumeModify(name, 120, 3000, 125)
}

func testAccDesktop_rootVolumeModify_step2(name string) string {
	return testAccDesktop_rootVolumeModify(name, 160, 3000, 125)
}

func testAccDesktop_rootVolumeModify_step3(name string) string {
	return testAccDesktop_rootVolumeModify(name, 160, 4000, 125)
}

func TestAccDesktop_dataVolumeModify(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDesktopFunc)

		name = acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccDesktop_dataVolumeModify_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "data_volume.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.type", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.throughput", "125"),
				),
			},
			{
				Config: testAccDesktop_dataVolumeModify_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "data_volume.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.type", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.throughput", "125"),
				),
			},
			{
				Config: testAccDesktop_dataVolumeModify_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "data_volume.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.type", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.throughput", "250"),
				),
			},
			{
				Config: testAccDesktop_dataVolumeModify_step4(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "data_volume.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(resourceName, "data_volume.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.type", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.size", "120"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.size", "120"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.iops", "5000"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.throughput", "250"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.2.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.2.size", "120"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.3.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.3.size", "120"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.3.iops", "5000"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.3.throughput", "250"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.4.type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.4.size", "120"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.5.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.5.size", "120"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.5.iops", "5000"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.5.throughput", "250"),
				),
			},
		},
	})
}

func testAccDesktop_dataVolumeModify(name string, size, iops, throughput int) string {
	varBlocks := fmt.Sprintf(`
variable "root_volume" {
  type = object({
    type = string
    size = number
  })

  default = {
    type = "SSD"
    size = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number, null)
    throughput = optional(number, null)
  }))

  default = [
    { type = "SSD", size = %[1]d },
    { type = "GPSSD2", size = %[1]d, iops = %[2]d, throughput = %[3]d },
  ]
}
`, size, iops, throughput)
	return testAccDesktop_volumeChange_resource(varBlocks, name)
}

func testAccDesktop_dataVolumeModify_moreVolume(name string, size, iops, throughput int) string {
	varBlocks := fmt.Sprintf(`
variable "root_volume" {
  type = object({
    type = string
    size = number
  })

  default = {
    type = "SSD"
    size = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number, null)
    throughput = optional(number, null)
  }))

  default = [
    { type = "SSD", size = %[1]d },
    { type = "GPSSD2", size = %[1]d, iops = %[2]d, throughput = %[3]d },
    { type = "SAS", size = %[1]d },
    { type = "GPSSD2", size = %[1]d, iops = %[2]d, throughput = %[3]d },
    { type = "GPSSD", size = %[1]d },
    { type = "GPSSD2", size = %[1]d, iops = %[2]d, throughput = %[3]d },
  ]
}
`, size, iops, throughput)
	return testAccDesktop_volumeChange_resource(varBlocks, name)
}

func testAccDesktop_dataVolumeModify_step1(name string) string {
	return testAccDesktop_dataVolumeModify(name, 80, 3000, 125)
}

func testAccDesktop_dataVolumeModify_step2(name string) string {
	return testAccDesktop_dataVolumeModify(name, 100, 3000, 125)
}

func testAccDesktop_dataVolumeModify_step3(name string) string {
	return testAccDesktop_dataVolumeModify(name, 100, 4000, 250)
}

func testAccDesktop_dataVolumeModify_step4(name string) string {
	return testAccDesktop_dataVolumeModify_moreVolume(name, 120, 5000, 250)
}

func TestAccDesktop_volumeChangeErrors(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDesktopFunc)

		name = acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccDesktop_volumeChangeErrors_base(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_typeChange(name),
				ExpectError: regexp.MustCompile(`volume type does not support updates`),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_sizeDecrease(name),
				ExpectError: regexp.MustCompile(`volume.*size.*cannot be smaller`),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_nonGPSSD2QoS(name),
				ExpectError: regexp.MustCompile(`the type of the volume.*is not GPSSD2, cannot set QoS options`),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_GPSSD2MissingQoS(name),
				ExpectError: regexp.MustCompile(`the type of the volume \(index number: \d+\) is GPSSD2, iops and throughput cannot be empty`),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_volumeCountDecrease(name),
				ExpectError: regexp.MustCompile(`the number of volumes cannot be reduced`),
			},
		},
	})
}

func testAccDesktop_volumeChangeErrors_base(name string) string {
	varBlocks := `
variable "root_volume" {
  type = object({
    type = string
    size = number
  })

  default = {
    type = "SAS"
    size = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number, null)
    throughput = optional(number, null)
  }))

  default = [
    { type = "SSD", size = 100 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
  ]
}
`
	return testAccDesktop_volumeChange_resource(varBlocks, name)
}

func testAccDesktop_volumeChangeErrors_typeChange(name string) string {
	varBlocks := `
variable "root_volume" {
  type = object({
    type = string
    size = number
  })

  default = {
    type = "SSD"
    size = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number, null)
    throughput = optional(number, null)
  }))

  default = [
    { type = "SSD", size = 100 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
  ]
}
`
	return testAccDesktop_volumeChange_resource(varBlocks, name)
}

func testAccDesktop_volumeChangeErrors_sizeDecrease(name string) string {
	varBlocks := `
variable "root_volume" {
  type = object({
    type = string
    size = number
  })

  default = {
    type = "SAS"
    size = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number, null)
    throughput = optional(number, null)
  }))

  default = [
    { type = "SSD", size = 80 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
  ]
}
`
	return testAccDesktop_volumeChange_resource(varBlocks, name)
}

func testAccDesktop_volumeChangeErrors_nonGPSSD2QoS(name string) string {
	varBlocks := `
variable "root_volume" {
  type = object({
    type = string
    size = number
  })

  default = {
    type = "SAS"
    size = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number, null)
    throughput = optional(number, null)
  }))

  default = [
    { type = "SSD", size = 100, iops = 3000, throughput = 125 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
  ]
}
`
	return testAccDesktop_volumeChange_resource(varBlocks, name)
}

func testAccDesktop_volumeChangeErrors_GPSSD2MissingQoS(name string) string {
	varBlocks := `
variable "root_volume" {
  type = object({
    type = string
    size = number
  })

  default = {
    type = "SAS"
    size = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number, null)
    throughput = optional(number, null)
  }))

  default = [
    { type = "SSD", size = 100 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
    { type = "GPSSD2", size = 100 },
  ]
}
`
	return testAccDesktop_volumeChange_resource(varBlocks, name)
}

func testAccDesktop_volumeChangeErrors_volumeCountDecrease(name string) string {
	varBlocks := `
variable "root_volume" {
  type = object({
    type = string
    size = number
  })

  default = {
    type = "SAS"
    size = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number, null)
    throughput = optional(number, null)
  }))

  default = [
    { type = "SSD", size = 100 },
  ]
}
`
	return testAccDesktop_volumeChange_resource(varBlocks, name)
}
