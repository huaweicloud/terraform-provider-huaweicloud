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
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
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
					resource.TestCheckResourceAttrPair(resourceName, "nic.0.network_id", "data.huaweicloud_workspace_service.test", "network_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "email_notification", "true"),
				),
			},
			{
				Config: testAccDesktop_basic_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
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
					resource.TestCheckResourceAttrPair(resourceName, "nic.0.network_id", "data.huaweicloud_workspace_service.test", "network_ids.0"),
				),
			},
			{
				Config: testAccDesktop_basic_step3(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
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
					resource.TestCheckResourceAttrPair(resourceName, "nic.0.network_id", "data.huaweicloud_workspace_service.test", "network_ids.0"),
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

func testAccDesktop_base() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_images" "test" {
  image_id   = "%[1]s"
  visibility = "market"
}
`, acceptance.HW_IMAGE_ID)
}

func testAccDesktop_basic(rName string, rootVolumeSize int) string {
	return fmt.Sprintf(`
%[1]s

locals {
  data_volume_sizes = [50, 70]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "NOT_FOUND")
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, "NOT_FOUND")
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    try(data.huaweicloud_workspace_service.test.desktop_security_group[0].id, "NOT_FOUND"), 
    try(data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id, "NOT_FOUND")
  ]

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
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
`, testAccDesktop_base(), rName, rootVolumeSize)
}

func testAccDesktop_basic_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  data_volume_sizes = [50, 90, 20, 40]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "NOT_FOUND")
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, "NOT_FOUND")
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    try(data.huaweicloud_workspace_service.test.desktop_security_group[0].id, "NOT_FOUND"), 
    try(data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id, "NOT_FOUND")
  ]

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
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
`, testAccDesktop_base(), rName)
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
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
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
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "user_name", "user-"+rName),
					resource.TestCheckResourceAttr(resourceName, "user_group", "administrators"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testAccDesktop_withEPSId(rName, epsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id             = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "NOT_FOUND")
  image_type            = "market"
  image_id              = try(data.huaweicloud_images_images.test.images[0].id, "NOT_FOUND")
  enterprise_project_id = "%[3]s" 
  availability_zone     = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  vpc_id                = data.huaweicloud_workspace_service.test.vpc_id
  security_groups       = [
    try(data.huaweicloud_workspace_service.test.desktop_security_group[0].id, "NOT_FOUND"), 
    try(data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id, "NOT_FOUND")
  ]

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
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
`, testAccDesktop_base(), rName, epsId)
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

func testAccdesktop_powerAction(rName string, powerAction string, powerActionType string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "NOT_FOUND")
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, "NOT_FOUND")
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    try(data.huaweicloud_workspace_service.test.desktop_security_group[0].id, "NOT_FOUND"), 
    try(data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id, "NOT_FOUND")
  ]

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
}
`, testAccDesktop_base(), rName, powerAction, powerActionType)
}

func testAccDesktop_volumeChange_resource(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "NOT_FOUND")
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, "NOT_FOUND")
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    try(data.huaweicloud_workspace_service.test.desktop_security_group[0].id, "NOT_FOUND"), 
    try(data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id, "NOT_FOUND")
  ]

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
  }

  name       = "%[1]s"
  user_name  = "user-%[1]s"
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type       = var.root_volume.type
    size       = var.root_volume.size
    iops       = try(var.root_volume.iops, null)
    throughput = try(var.root_volume.throughput, null)
    kms_id     = try(var.root_volume.kms_id, null)
  }

  dynamic "data_volume" {
    for_each = var.data_volumes

    content {
      type       = data_volume.value.type
      size       = data_volume.value.size
      iops       = try(data_volume.value.iops, null)
      throughput = try(data_volume.value.throughput, null)
      kms_id     = try(data_volume.value.kms_id, null)
    }
  }

  tags = {
    foo = "bar"
  }

  email_notification = true
  delete_user        = true
}
`, rName)
}

func TestAccDesktop_rootVolumeModify(t *testing.T) {
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
				// Test for create GPSSD2 root volume.
				Config: testAccDesktop_rootVolumeModify_step1(rName),
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
				Config: testAccDesktop_rootVolumeModify_step2(rName),
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
				Config: testAccDesktop_rootVolumeModify_step3(rName),
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

func testAccDesktop_rootVolumeModify_step1(rName string) string {
	return testAccDesktop_rootVolumeModify(rName, 120, 3000, 125)
}

func testAccDesktop_rootVolumeModify_step2(rName string) string {
	return testAccDesktop_rootVolumeModify(rName, 160, 3000, 125)
}

func testAccDesktop_rootVolumeModify_step3(rName string) string {
	return testAccDesktop_rootVolumeModify(rName, 160, 4000, 125)
}

func testAccDesktop_rootVolumeModify(rName string, size, iops, throughput int) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  })

  default = {
    type       = "GPSSD2"
    size       = %[3]d
	iops       = %[4]d
	throughput = %[5]d
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = 80 },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName), size, iops, throughput)
}

func TestAccDesktop_dataVolumeModify(t *testing.T) {
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
				Config: testAccDesktop_dataVolumeModify_step1(rName),
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
				Config: testAccDesktop_dataVolumeModify_step2(rName),
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
				Config: testAccDesktop_dataVolumeModify_step3(rName),
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
				Config: testAccDesktop_dataVolumeModify_step4(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "data_volume.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(resourceName, "data_volume.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.type", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.0.size", "100"),

					resource.TestCheckResourceAttr(resourceName, "data_volume.1.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.1.throughput", "250"),

					resource.TestCheckResourceAttr(resourceName, "data_volume.2.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.2.size", "100"),

					resource.TestCheckResourceAttr(resourceName, "data_volume.3.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.3.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.3.iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.3.throughput", "250"),
				),
			},
			{
				Config: testAccDesktop_dataVolumeModify_step5(rName),
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

func testAccDesktop_dataVolumeModify_step1(rName string) string {
	return testAccDesktop_dataVolumeModify(rName, 80, 3000, 125)
}

func testAccDesktop_dataVolumeModify_step2(rName string) string {
	return testAccDesktop_dataVolumeModify(rName, 100, 3000, 125)
}

func testAccDesktop_dataVolumeModify_step3(rName string) string {
	return testAccDesktop_dataVolumeModify(rName, 100, 4000, 250)
}

func testAccDesktop_dataVolumeModify_step4(rName string) string {
	return testAccDesktop_dataVolumeModify_moreVolume(rName, 100, 4000, 250)
}

func testAccDesktop_dataVolumeModify_step5(rName string) string {
	return testAccDesktop_dataVolumeModify_complex(rName, 120, 5000, 250)
}

func testAccDesktop_dataVolumeModify(rName string, size, iops, throughput int) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  })

  default = {
    type       = "SSD"
    size       = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = %[3]d },
    { type = "GPSSD2", size = %[3]d, iops = %[4]d, throughput = %[5]d },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName), size, iops, throughput)
}

func testAccDesktop_dataVolumeModify_moreVolume(rName string, size, iops, throughput int) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  })

  default = {
    type       = "SSD"
    size       = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = %[3]d },
    { type = "GPSSD2", size = %[3]d, iops = %[4]d, throughput = %[5]d },
    { type = "SAS", size = %[3]d },
    { type = "GPSSD2", size = %[3]d, iops = %[4]d, throughput = %[5]d },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName), size, iops, throughput)
}

func testAccDesktop_dataVolumeModify_complex(rName string, size, iops, throughput int) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  })

  default = {
    type       = "SSD"
    size       = 80
  }
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = %[3]d },
    { type = "GPSSD2", size = %[3]d, iops = %[4]d, throughput = %[5]d },
    { type = "SAS", size = %[3]d },
    { type = "GPSSD2", size = %[3]d, iops = %[4]d, throughput = %[5]d },
    { type = "GPSSD", size = %[3]d },
    { type = "GPSSD2", size = %[3]d, iops = %[4]d, throughput = %[5]d },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName), size, iops, throughput)
}

func TestAccDesktop_volumeChangeErrors(t *testing.T) {
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
				Config: testAccDesktop_volumeChangeErrors_base(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_typeChange(rName),
				ExpectError: regexp.MustCompile(`volume type does not support updates`),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_sizeDecrease(rName),
				ExpectError: regexp.MustCompile(`volume.*size.*cannot be smaller`),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_nonGPSSD2QoS(rName),
				ExpectError: regexp.MustCompile(`the type of the volume.*is not GPSSD2, cannot set QoS options`),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_GPSSD2MissingQoS(rName),
				ExpectError: regexp.MustCompile(`the type of the .*th volume is GPSSD2, iops and throughput cannot be empty`),
			},
			{
				Config:      testAccDesktop_volumeChangeErrors_volumeCountDecrease(rName),
				ExpectError: regexp.MustCompile(`the number of volumes cannot be reduced`),
			},
		},
	})
}

func testAccDesktop_volumeChangeErrors_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
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
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = 100 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName))
}

func testAccDesktop_volumeChangeErrors_typeChange(rName string) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
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
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = 100 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName))
}

func testAccDesktop_volumeChangeErrors_sizeDecrease(rName string) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
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
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = 80 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName))
}

func testAccDesktop_volumeChangeErrors_nonGPSSD2QoS(rName string) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
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
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = 100, iops = 3000, throughput = 125 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName))
}

func testAccDesktop_volumeChangeErrors_GPSSD2MissingQoS(rName string) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
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
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = 100 },
    { type = "GPSSD2", size = 100, iops = 3000, throughput = 125 },
    { type = "GPSSD2", size = 100 },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName))
}

func testAccDesktop_volumeChangeErrors_volumeCountDecrease(rName string) string {
	return fmt.Sprintf(`
%[1]s

variable "root_volume" {
  type = object({
    type       = string
    size       = number
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
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
    iops       = optional(number)
    throughput = optional(number)
    kms_id     = optional(string)
  }))

  default = [
    { type = "SSD", size = 100 },
  ]
}

%[2]s
`, testAccDesktop_base(), testAccDesktop_volumeChange_resource(rName))
}
