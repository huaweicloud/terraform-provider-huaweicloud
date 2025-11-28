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

func getDesktopPoolFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating Workspace client: %s", err)
	}
	return workspace.GetDesktopPoolById(client, state.Primary.ID)
}

func TestAccDesktopPool_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()

		desktopPool  interface{}
		resourceName = "huaweicloud_workspace_desktop_pool.test"
		rc           = acceptance.InitResourceCheck(resourceName, &desktopPool, getDesktopPoolFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopPoolImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopPool_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "DYNAMIC"),
					resource.TestCheckResourceAttr(resourceName, "size", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "product_id",
						"data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "80"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_ids.0",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "security_groups.0.id"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "authorized_objects.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "authorized_objects.0.object_id"),
					resource.TestCheckResourceAttrSet(resourceName, "authorized_objects.0.object_name"),
					resource.TestCheckResourceAttr(resourceName, "authorized_objects.0.object_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "authorized_objects.0.user_group", "administrators"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "disconnected_retention_period", "10"),
					resource.TestCheckResourceAttr(resourceName, "enable_autoscale", "true"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_policy.0.autoscale_type", "AUTO_CREATED"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_policy.0.min_idle", "1"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_policy.0.max_auto_created", "2"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_policy.0.once_auto_created", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "desktop_name_policy_id",
						"huaweicloud_workspace_desktop_name_rule.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "in_maintenance_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", "Create a dynamic desktop pool"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
			{
				Config: testAccDesktopPool_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", ""),
					resource.TestCheckResourceAttr(resourceName, "disconnected_retention_period", "43200"),
					resource.TestCheckResourceAttr(resourceName, "enable_autoscale", "false"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "desktop_name_policy_id", ""),
					resource.TestCheckResourceAttr(resourceName, "in_maintenance_mode", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(resourceName, "root_volume.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "data_volumes.0.id"),
					resource.TestCheckResourceAttr(resourceName, "status", "STEADY"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_used"),
					resource.TestCheckResourceAttrSet(resourceName, "product.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "product.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "product.0.cpu"),
					resource.TestCheckResourceAttrSet(resourceName, "product.0.memory"),
					resource.TestCheckResourceAttrSet(resourceName, "product.0.descriptions"),
					resource.TestCheckResourceAttrSet(resourceName, "product.0.charging_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "image_os_type"),
					resource.TestCheckResourceAttrSet(resourceName, "image_os_version"),
					resource.TestCheckResourceAttrSet(resourceName, "image_os_platform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"image_type",
					"vpc_id",
					"tags",
				},
			},
		},
	})
}

func testAccDesktopPool_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  os_type           = "Windows"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_desktop_name_rule" "test" {
  name                         = replace("%[1]s", "-", "_")
  name_prefix                  = "pool"
  digit_number                 = 1
  start_number                 = 1
  single_domain_user_increment = 0
}
`, name)
}

func testAccDesktopPool_basic_base(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_user" "test" {
  count = 2
  name  = "%[3]s${count.index}"
  email = "user@test.com"
}
`, common.TestBaseNetwork(name), testAccDesktopPool_base(name), name)
}

func testAccDesktopPool_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  data_volume_sizes = [50, 70]
}

resource "huaweicloud_workspace_desktop_pool" "test" {
  name                          = "%[2]s"
  type                          = "DYNAMIC"
  size                          = 2
  product_id                    = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "")
  image_type                    = "gold"
  image_id                      = "%[3]s"
  subnet_ids                    = [huaweicloud_vpc_subnet.test.id]
  vpc_id                        = huaweicloud_vpc.test.id
  availability_zone             = data.huaweicloud_availability_zones.test.names[0]
  disconnected_retention_period = 10
  enable_autoscale              = true
  desktop_name_policy_id        = huaweicloud_workspace_desktop_name_rule.test.id
  description                   = "Create a dynamic desktop pool"
  in_maintenance_mode           = true

  root_volume {
    type = "SAS"
    size = 80
  }

  security_groups {
    id = data.huaweicloud_workspace_service.test.desktop_security_group.0.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }

  dynamic "data_volumes" {
    for_each = local.data_volume_sizes

    content {
      type = "SAS"
      size = data_volumes.value
    }
  }

  dynamic "authorized_objects" {
    for_each = huaweicloud_workspace_user.test[*]

    content {
      object_id   = authorized_objects.value.id
      object_type = "USER"
      object_name = authorized_objects.value.name
      user_group  = "administrators"
    }
  }

  autoscale_policy {
    autoscale_type    = "AUTO_CREATED"
    min_idle          = 1
    max_auto_created  = 2
    once_auto_created = 1
  }

  tags = {
    foo = "bar"
  }
}
`, testAccDesktopPool_basic_base(name), name, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID)
}

func testAccDesktopPool_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  data_volume_sizes = [50, 70]
}

resource "huaweicloud_workspace_desktop_pool" "test" {
  name                          = "%[2]s"
  type                          = "DYNAMIC"
  size                          = 2
  product_id                    = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "")
  image_type                    = "gold"
  image_id                      = "%[3]s"
  subnet_ids                    = [huaweicloud_vpc_subnet.test.id]
  vpc_id                        = huaweicloud_vpc.test.id
  disconnected_retention_period = 43200

  root_volume {
    type = "SAS"
    size = 80
  }


  security_groups {
    id = data.huaweicloud_workspace_service.test.desktop_security_group.0.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }

  dynamic "data_volumes" {
    for_each = local.data_volume_sizes

    content {
      type = "SAS"
      size = data_volumes.value
    }
  }

  dynamic "authorized_objects" {
    for_each = huaweicloud_workspace_user.test[*]

    content {
      object_id   = authorized_objects.value.id
      object_type = "USER"
      object_name = authorized_objects.value.name
      user_group  = "administrators"
    }
  }
}
`, testAccDesktopPool_basic_base(name), updateName, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID)
}

// Before running this use case, you must add the OU on the AD server to the Workspace.
func TestAccDesktopPool_localAD(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		desktopPool  interface{}
		resourceName = "huaweicloud_workspace_desktop_pool.test"
		rc           = acceptance.InitResourceCheck(resourceName, &desktopPool, getDesktopPoolFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopPoolImageId(t)
			acceptance.TestAccPreCheckWorkspaceOUName(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopPool_localAD_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "STATIC"),
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "product_id",
						"data.huaweicloud_workspace_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "80"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_ids.0",
						"data.huaweicloud_workspace_service.test", "network_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "authorized_objects.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", ""),
					resource.TestCheckResourceAttr(resourceName, "disconnected_retention_period", "0"),
					resource.TestCheckResourceAttr(resourceName, "enable_autoscale", "false"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_policy.0.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "desktop_name_policy_id", ""),
					resource.TestCheckResourceAttr(resourceName, "ou_name", acceptance.HW_WORKSPACE_OU_NAME),
					resource.TestCheckResourceAttr(resourceName, "in_maintenance_mode", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccDesktopPool_localAD_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "enable_autoscale", "true"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_policy.0.autoscale_type", "ACCESS_CREATED"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_policy.0.max_auto_created", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "desktop_name_policy_id",
						"huaweicloud_workspace_desktop_name_rule.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "ou_name", ""),
					resource.TestCheckResourceAttr(resourceName, "in_maintenance_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated desktop pool"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"image_type",
					"tags",
					"ou_name",
				},
			},
		},
	})
}

func testAccDesktopPool_localAD_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_desktop_pool" "test" {
  name                  = "%[2]s"
  type                  = "STATIC"
  size                  = 1
  product_id            = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "")
  image_type            = "gold"
  image_id              = "%[3]s"
  subnet_ids            = try(slice(data.huaweicloud_workspace_service.test.network_ids, 0, 1), [])
  enterprise_project_id = "%[4]s"
  ou_name               = "%[5]s"

  root_volume {
    type = "SAS"
    size = 80
  }
}
`, testAccDesktopPool_base(name),
		name,
		acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
		acceptance.HW_WORKSPACE_OU_NAME)
}

func testAccDesktopPool_localAD_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_desktop_pool" "test" {
  name                   = "%[2]s"
  type                   = "STATIC"
  size                   = 1
  product_id             = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "")
  image_type             = "gold"
  image_id               = "%[3]s"
  subnet_ids             = try(slice(data.huaweicloud_workspace_service.test.network_ids, 0, 1), [])
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  enable_autoscale       = true
  desktop_name_policy_id = huaweicloud_workspace_desktop_name_rule.test.id
  description            = "Updated desktop pool"
  in_maintenance_mode    = true
  enterprise_project_id  = "%[4]s"
  
  root_volume {
    type = "SAS"
    size = 80
  }

  autoscale_policy {
    autoscale_type   = "ACCESS_CREATED"
    max_auto_created = 1
  }

  tags = {
    foo = "update"
  }
}
`, testAccDesktopPool_base(name),
		name,
		acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
