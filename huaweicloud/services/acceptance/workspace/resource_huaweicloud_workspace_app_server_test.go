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

func getResourceAppServerFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	return workspace.GetServerById(client, state.Primary.ID)
}

// Before running this test, please enable a service that connects to LocalAD and the corresponding OU is created.
func TestAccResourceAppServer_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_server.test"
		name         = acceptance.RandomAccResourceName()

		server interface{}
		rc     = acceptance.InitResourceCheck(
			resourceName,
			&server,
			getResourceAppServerFunc,
		)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceAppServer_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "server_group_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id",
						"data.huaweicloud_workspace_app_server_groups.test", "server_groups.0.product_id"),
					resource.TestCheckResourceAttrPair(resourceName, "root_volume.0.type",
						"data.huaweicloud_workspace_app_server_groups.test", "server_groups.0.system_disk_type"),
					resource.TestCheckResourceAttrPair(resourceName, "root_volume.0.size",
						"data.huaweicloud_workspace_app_server_groups.test", "server_groups.0.system_disk_size"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"data.huaweicloud_vpc_subnets.test", "subnets.0.vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"data.huaweicloud_workspace_app_server_groups.test", "server_groups.0.subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "update_access_agent", "false"),
					resource.TestCheckResourceAttr(resourceName, "ou_name", acceptance.HW_WORKSPACE_OU_NAME),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "maintain_status", "true"),
				),
			},
			{
				Config: testResourceAppServer_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "maintain_status", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"type",
					"vpc_id",
					"subnet_id",
					"update_access_agent",
					"scheduler_hints",
				},
			},
		},
	})
}

func testResourceAppServer_base(name string) string {
	return fmt.Sprintf(`
variable "ou_name" {
  type    = string
  default = "%[1]s"
}

data "huaweicloud_workspace_app_server_groups" "test" {
  server_group_id = "%[2]s"
}

data "huaweicloud_vpc_subnets" "test" {
  id = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
}

`, acceptance.HW_WORKSPACE_OU_NAME, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID)
}

func testResourceAppServer_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server" "test" {
  name                = "%[2]s" 
  server_group_id     = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].id, null)
  type                = "createApps"
  flavor_id           = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].product_id, null)
  vpc_id              = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, null)
  subnet_id           = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
  update_access_agent = false
  ou_name             = var.ou_name != "" ? var.ou_name : null
  description         = "Created by terraform script"
  maintain_status     = true

  root_volume {
    type = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].system_disk_type, null)
    size = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].system_disk_size, null)
  }
}
`, testResourceAppServer_base(name), name, acceptance.HW_WORKSPACE_OU_NAME)
}

func testResourceAppServer_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server" "test" {
  name                = "%[2]s_update" 
  server_group_id     = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].id, null)
  type                = "createApps"
  flavor_id           = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].product_id, null)
  vpc_id              = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, null)
  subnet_id           = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
  update_access_agent = false
  ou_name             = var.ou_name != "" ? var.ou_name : null
  maintain_status     = false

  root_volume {
    type = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].system_disk_type, null)
    size = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].system_disk_size, null)
  }
}
`, testResourceAppServer_base(name), name, acceptance.HW_WORKSPACE_OU_NAME)
}

// Before running this test, please enable a service that connects to LocalAD and the corresponding OU is created.
func TestAccResourceAppServer_prepaid(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_server.test"
		name         = acceptance.RandomAccResourceName()

		server interface{}
		rc     = acceptance.InitResourceCheck(
			resourceName,
			&server,
			getResourceAppServerFunc,
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
				Config: testResourceAppServer_prepaid_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "server_group_id", "huaweicloud_workspace_app_server_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id", "huaweicloud_workspace_app_server_group.test", "flavor_id"),
					resource.TestCheckResourceAttrPair(resourceName, "root_volume.0.type",
						"huaweicloud_workspace_app_server_group.test", "system_disk_type"),
					resource.TestCheckResourceAttrPair(resourceName, "root_volume.0.size",
						"huaweicloud_workspace_app_server_group.test", "system_disk_size"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", acceptance.HW_WORKSPACE_AD_VPC_ID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", acceptance.HW_WORKSPACE_AD_NETWORK_ID),
					resource.TestCheckResourceAttr(resourceName, "description", "Created server by script"),
					resource.TestCheckResourceAttr(resourceName, "ou_name", acceptance.HW_WORKSPACE_OU_NAME),
					resource.TestCheckResourceAttr(resourceName, "maintain_status", "true"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
			{
				Config: testResourceAppServer_prepaid_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "maintain_status", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"type",
					"vpc_id",
					"subnet_id",
					"update_access_agent",
					"scheduler_hints",
					"period_unit",
					"period",
					"auto_renew",
				},
			},
		},
	})
}

func testResourceAppServer_prepaid_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server" "test" {
  name            = "%[2]s" 
  server_group_id = huaweicloud_workspace_app_server_group.test.id
  type            = "createApps"
  flavor_id       = huaweicloud_workspace_app_server_group.test.flavor_id

  root_volume {
    type = huaweicloud_workspace_app_server_group.test.system_disk_type
    size = huaweicloud_workspace_app_server_group.test.system_disk_size
  }

  vpc_id              = huaweicloud_workspace_app_server_group.test.vpc_id
  subnet_id           = huaweicloud_workspace_app_server_group.test.subnet_id
  os_type             = "Windows"
  update_access_agent = false
  ou_name             = "%[3]s"
  description         = "Created server by script"
  maintain_status     = true

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}
`, testResourceAppServer_base(name), name, acceptance.HW_WORKSPACE_OU_NAME)
}

func testResourceAppServer_prepaid_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server" "test" {
  name            = "%[2]s_update" 
  server_group_id = huaweicloud_workspace_app_server_group.test.id
  type            = "createApps"
  flavor_id       = huaweicloud_workspace_app_server_group.test.flavor_id

  root_volume {
    type = huaweicloud_workspace_app_server_group.test.system_disk_type
    size = huaweicloud_workspace_app_server_group.test.system_disk_size
  }

  vpc_id              = huaweicloud_workspace_app_server_group.test.vpc_id
  subnet_id           = huaweicloud_workspace_app_server_group.test.subnet_id
  update_access_agent = false
  ou_name             = "%[3]s"
  maintain_status     = false

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
`, testResourceAppServer_base(name), name, acceptance.HW_WORKSPACE_OU_NAME)
}
