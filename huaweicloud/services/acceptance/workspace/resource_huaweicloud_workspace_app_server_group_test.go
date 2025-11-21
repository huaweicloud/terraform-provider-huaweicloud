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

func getResourceAppServerGroupFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	return workspace.GetServerGroupById(client, state.Primary.ID)
}

// Before running this test, please enable a service that connects to LocalAD and the corresponding OU is created.
func TestAccResourceAppServerGroup_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_server_group.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()

		serverGroup interface{}
		rc          = acceptance.InitResourceCheck(
			resourceName,
			&serverGroup,
			getResourceAppServerGroupFunc,
		)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
			acceptance.TestAccPreCheckWorkspaceOUName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceAppServerGroup_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"data.huaweicloud_workspace_service.test", "network_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "90"),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID),
					resource.TestCheckResourceAttr(resourceName, "image_type", "gold"),
					resource.TestCheckResourceAttr(resourceName, "image_product_id",
						acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID),
					resource.TestCheckResourceAttr(resourceName, "is_vdi", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "app_type", "COMMON_APP"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "ip_virtual.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.0.max_session", "3"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.0.cpu_threshold", "80"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.0.mem_threshold", "80"),
					resource.TestCheckResourceAttr(resourceName, "ou_name", acceptance.HW_WORKSPACE_OU_NAME),
					resource.TestCheckResourceAttr(resourceName, "extra_session_type", "CPU"),
					resource.TestCheckResourceAttr(resourceName, "extra_session_size", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "primary_server_group_id",
						"huaweicloud_workspace_app_server_group.primary", "id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "storage_mount_policy", "USER"),
				),
			},
			{
				Config: testResourceAppServerGroup_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "system_disk_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "80"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_type", "SESSION_DESKTOP_APP"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.0.max_session", "2"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.0.cpu_threshold", "85"),
					resource.TestCheckResourceAttr(resourceName, "route_policy.0.mem_threshold", "85"),
					resource.TestCheckResourceAttr(resourceName, "ou_name", ""),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_mount_policy", "ANY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"vpc_id", "image_type", "image_product_id", "availability_zone", "ip_virtual", "route_policy",
				},
			},
		},
	})
}

func testResourceAppServerGroup_base(name, appType, amountingPolicy string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "primary" {
  name             = "%[1]s_primary"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 80
  is_vdi           = false
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"

  ip_virtual {
    enable = true
  }

  app_type             = "%[5]s"
  storage_mount_policy = "%[6]s"
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		appType,
		amountingPolicy)
}

func testResourceAppServerGroup_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_group" "test" {
  name                    = "%[2]s"
  os_type                 = "Windows"
  flavor_id               = "%[3]s"
  vpc_id                  = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id               = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type        = "SAS"
  system_disk_size        = 90
  is_vdi                  = false
  image_id                = "%[4]s"
  image_type              = "gold"
  image_product_id        = "%[5]s"
  description             = "Created by script"
  ou_name                 = "%[6]s"
  primary_server_group_id = huaweicloud_workspace_app_server_group.primary.id

  tags = {
    foo = "bar"
  }

  ip_virtual {
    enable = true
  }

  extra_session_type = "CPU"
  extra_session_size = 2

  route_policy {
    max_session   = 3
    cpu_threshold = 80
    mem_threshold = 80
  }

  enabled              = false
  storage_mount_policy = "USER"
}
`, testResourceAppServerGroup_base(name, "COMMON_APP", "USER"), name,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		acceptance.HW_WORKSPACE_OU_NAME)
}

func testResourceAppServerGroup_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_group" "test" {
  name                    = "%[2]s"
  os_type                 = "Windows"
  flavor_id               = "%[3]s"
  vpc_id                  = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id               = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type        = "SAS"
  system_disk_size        = 80
  app_type                = "SESSION_DESKTOP_APP"
  is_vdi                  = false
  image_id                = "%[4]s"
  image_type              = "gold"
  image_product_id        = "%[5]s"
  primary_server_group_id = huaweicloud_workspace_app_server_group.primary.id

  tags = {
    foo = "bar"
  }

  ip_virtual {
    enable = true
  }

  extra_session_type = "CPU"
  extra_session_size = 2

  route_policy {
    max_session   = 2
    cpu_threshold = 85
    mem_threshold = 85
  }

  enabled              = true
  storage_mount_policy = "ANY"
}
`, testResourceAppServerGroup_base(name, "SESSION_DESKTOP_APP", "ANY"), name,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func TestAccResourceAppServerGroup_singleSession(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj          interface{}
		resourceName = "huaweicloud_workspace_app_server_group.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getResourceAppServerGroupFunc)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppServerGroup_singleSession_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"data.huaweicloud_workspace_service.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"data.huaweicloud_workspace_service.test", "network_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "90"),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID),
					resource.TestCheckResourceAttr(resourceName, "image_type", "gold"),
					resource.TestCheckResourceAttr(resourceName, "image_product_id",
						acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID),
					resource.TestCheckResourceAttr(resourceName, "is_vdi", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "app_type", "SESSION_DESKTOP_APP"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_mount_policy", "USER"),
				),
			},
			{
				Config: testAccResourceAppServerGroup_singleSession_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "system_disk_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "80"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_type", "SESSION_DESKTOP_APP"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "storage_mount_policy", "SHARE"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"vpc_id", "image_type", "image_product_id", "availability_zone",
				},
			},
		},
	})
}

func testAccResourceAppServerGroup_singleSession_step1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name                  = "%[1]s"
  os_type               = "Windows"
  flavor_id             = "%[2]s"
  vpc_id                = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id             = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type      = "SAS"
  system_disk_size      = 90
  app_type              = "SESSION_DESKTOP_APP"
  is_vdi                = true
  image_id              = "%[3]s"
  image_type            = "gold"
  image_product_id      = "%[4]s"
  description           = "Created by script"
  storage_mount_policy  = "USER"
  enterprise_project_id = "%[5]s"

  tags = {
    foo = "bar"
  }
}
`, name,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
	)
}

func testAccResourceAppServerGroup_singleSession_step2(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name                  = "%[1]s"
  os_type               = "Windows"
  flavor_id             = "%[2]s"
  vpc_id                = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id             = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type      = "SAS"
  system_disk_size      = 80
  app_type              = "SESSION_DESKTOP_APP"
  is_vdi                = true
  image_id              = "%[3]s"
  image_type            = "gold"
  image_product_id      = "%[4]s"
  storage_mount_policy  = "SHARE"
  enabled               = false
  enterprise_project_id = "%[5]s"

  tags = {
    foo   = "bar_update"
    owner = "terraform"
  }
}
`, name,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
	)
}
