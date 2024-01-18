package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/workspace/v2/accesspolicies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getAccessPolicyFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace v2 client: %s", err)
	}

	return workspace.GetAccessPolicyByPolicyName(client, state.Primary.Attributes["name"])
}

func TestAccAccessPolicy_basic(t *testing.T) {
	var (
		accessPolicy *accesspolicies.AccessPolicyDetailInfo
		resourceName = "huaweicloud_workspace_access_policy.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		rc           = acceptance.InitResourceCheck(resourceName, &accessPolicy, getAccessPolicyFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessPolicy_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "PRIVATE_ACCESS"),
					resource.TestCheckResourceAttr(resourceName, "blacklist_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceName, "blacklist.#", "2"),
				),
			},
			{
				Config: testAccAccessPolicy_basic_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "PRIVATE_ACCESS"),
					resource.TestCheckResourceAttr(resourceName, "blacklist_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceName, "blacklist.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAccessPolicyImportStateFunc(resourceName),
			},
		},
	})
}

func testAccAccessPolicyImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<name>', but got '%s'", rs.Primary.Attributes["name"])
		}
		return rs.Primary.Attributes["name"], nil
	}
}

func testAccAccessPolicy_basic_config(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_service" "test" {
  access_mode = "BOTH"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
  ]
}

resource "huaweicloud_workspace_user" "test" {
  depends_on = [huaweicloud_workspace_service.test]

  count = 3

  name  = format("%[2]s-%%d", count.index)
  email = "basic@example.com"
}
`, common.TestBaseNetwork(name), name)
}

func testAccAccessPolicy_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_access_policy" "test" {
  blacklist_type = "INTERNET"
  name           = "PRIVATE_ACCESS"

  dynamic "blacklist" {
    for_each = slice(huaweicloud_workspace_user.test[*].id, 0, 2)

    content {
      object_type = "USER"
      object_id   = blacklist.value
    }
  }
}
`, testAccAccessPolicy_basic_config(name), name)
}

func testAccAccessPolicy_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_access_policy" "test" {
  blacklist_type = "INTERNET"
  name           = "PRIVATE_ACCESS"

  dynamic "blacklist" {
    for_each = slice(huaweicloud_workspace_user.test[*].id, 1, 3)

    content {
      object_type = "USER"
      object_id   = blacklist.value
    }
  }
}
`, testAccAccessPolicy_basic_config(name), name)
}
