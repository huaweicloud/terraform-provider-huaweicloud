package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/workspace/v2/policygroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPolicyGroupFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace v2 client: %s", err)
	}
	return policygroups.Get(client, state.Primary.ID)
}

func TestAccPolicyGroup_basic(t *testing.T) {
	var (
		policyGroup  policygroups.PolicyGroup
		resourceName = "huaweicloud_workspace_policy_group.test"
		baseConfig   = testAccPolicyGroup_base()
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&policyGroup,
		getPolicyGroupFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyGroup_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.type", "USER"),
					resource.TestCheckResourceAttrPair(resourceName, "targets.0.id", "huaweicloud_workspace_user.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "targets.0.name", "huaweicloud_workspace_user.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.access_control.0.ip_access_control",
						"112.20.53.1|255.255.240.0;112.20.53.2|255.255.240.0"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccPolicyGroup_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "targets.0.type", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.id", "default-apply-all-targets"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.name", "All-Targets"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.access_control.0.ip_access_control",
						"112.20.53.2|255.255.240.0;112.20.53.3|255.255.240.0"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
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

func testAccPolicyGroup_base() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/20"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
  ]
}

resource "huaweicloud_workspace_user" "test" {
  depends_on = [huaweicloud_workspace_service.test]

  name  = "%[1]s"
  email = "basic@example.com"

  password_never_expires = false
  disabled               = false
}
`, name)
}

func testAccPolicyGroup_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_policy_group" "test" {
  name        = "%[2]s"
  priority    = 1
  description = "Created by script"

  targets {
    type = "USER"
    id   = huaweicloud_workspace_user.test.id
    name = huaweicloud_workspace_user.test.name
  }
  policy {
    access_control {
      ip_access_control = "112.20.53.1|255.255.240.0;112.20.53.2|255.255.240.0"
    }
  }
}
`, baseConfig, name)
}

func testAccPolicyGroup_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_policy_group" "test" {
  name     = "%[2]s"
  priority = 1

  targets {
    type = "ALL"
    id   = "default-apply-all-targets"
    name = "All-Targets"
  }
  policy {
    access_control {
      ip_access_control = "112.20.53.2|255.255.240.0;112.20.53.3|255.255.240.0"
    }
  }
}
`, baseConfig, name)
}
