package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/terminals"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTerminalBindingFunc(conf *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating Workspace v2 client: %s", err)
	}
	opts := terminals.ListOpts{
		Offset: utils.Int(0), // Offset is the required query parameter.
		Limit:  1000,         // Limit is the required query parameter.
	}
	resp, err := terminals.List(client, opts)
	if len(resp) > 0 || err != nil {
		return resp, err
	}
	return nil, golangsdk.ErrDefault404{}
}

func TestAccTerminalBinding_basic(t *testing.T) {
	var (
		bindings     []terminals.TerminalBindingResp
		resourceName = "huaweicloud_workspace_terminal_binding.test"
		name         = acceptance.RandomAccResourceNameWithDash()
		rc           = acceptance.InitResourceCheck(resourceName, &bindings, getTerminalBindingFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTerminalBinding_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "disabled_after_delete", "false"),
					resource.TestCheckResourceAttr(resourceName, "bindings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0.mac", "FA-16-3E-E2-3A-1D"),
					resource.TestCheckResourceAttrPair(resourceName, "bindings.0.desktop_name", "huaweicloud_workspace_desktop.test1", "name"),
				),
			},
			{
				Config: testAccTerminalBinding_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "disabled_after_delete", "true"),
					resource.TestCheckResourceAttr(resourceName, "bindings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0.mac", "FA-16-3E-E2-3A-1D"),
					resource.TestCheckResourceAttrPair(resourceName, "bindings.0.desktop_name", "huaweicloud_workspace_desktop.test1", "name"),
				),
			},
			{
				Config: testAccTerminalBinding_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "disabled_after_delete", "true"),
					resource.TestCheckResourceAttr(resourceName, "bindings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0.mac", "FA-16-3E-E2-3A-1E"),
					resource.TestCheckResourceAttrPair(resourceName, "bindings.0.desktop_name", "huaweicloud_workspace_desktop.test2", "name"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disabled_after_delete"},
			},
		},
	})
}

func testAccTerminalBinding_base(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
  ]
}

resource "huaweicloud_workspace_desktop" "test1" {
  flavor_id         = "workspace.x86.ultimate.large2"
  image_type        = "market"
  image_id          = "8451dedf-b353-43aa-b5fb-5bccadda2207"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  security_groups   = [
    huaweicloud_workspace_service.test.desktop_security_group.0.id,
    huaweicloud_networking_secgroup.test.id,
  ]

  nic {
    network_id = huaweicloud_vpc_subnet.test.id
  }

  name        = "%[2]s-1"
  user_name   = "%[2]s-user-1" // The user_name cannot be same as the desktop_name
  user_email  = "basic@example.com"
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

resource "huaweicloud_workspace_desktop" "test2" {
  // Concurrent creation crashes the service, making it impossible to delete cloud desktops.
  depends_on = [huaweicloud_workspace_desktop.test1]

  flavor_id         = "workspace.x86.ultimate.large2"
  image_type        = "market"
  image_id          = "8451dedf-b353-43aa-b5fb-5bccadda2207"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  security_groups   = [
    huaweicloud_workspace_service.test.desktop_security_group.0.id,
    huaweicloud_networking_secgroup.test.id,
  ]

  nic {
    network_id = huaweicloud_vpc_subnet.test.id
  }

  name        = "%[2]s-2"
  user_name   = "%[2]s-user-2" // The user_name cannot be same as the desktop_name
  user_email  = "basic@example.com"
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

func testAccTerminalBinding_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_terminal_binding" "test" {
  enabled               = false
  disabled_after_delete = false

  bindings {
    desktop_name = huaweicloud_workspace_desktop.test1.name
    mac          = "FA-16-3E-E2-3A-1D"
  }
}
`, testAccTerminalBinding_base(name))
}

func testAccTerminalBinding_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_terminal_binding" "test" {
  enabled               = true
  disabled_after_delete = true

  bindings {
    desktop_name = huaweicloud_workspace_desktop.test1.name
    mac          = "FA-16-3E-E2-3A-1D"
  }
}
`, testAccTerminalBinding_base(name))
}

func testAccTerminalBinding_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_terminal_binding" "test" {
  enabled               = true
  disabled_after_delete = true

  bindings {
    desktop_name = huaweicloud_workspace_desktop.test2.name
    mac          = "FA-16-3E-E2-3A-1E"
  }
}
`, testAccTerminalBinding_base(name))
}
