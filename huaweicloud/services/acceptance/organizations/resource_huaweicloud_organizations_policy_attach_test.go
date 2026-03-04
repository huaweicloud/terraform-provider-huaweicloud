package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/organizations"
)

func getPolicyAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("organizations", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	return organizations.GetPolicyAttachedEntity(client, state.Primary.Attributes["policy_id"], state.Primary.Attributes["entity_id"])
}

func TestAccPolicyAttach_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_organizations_policy_attach.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getPolicyAttachResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyAttach_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_organizations_policy.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "entity_id",
						"huaweicloud_organizations_organizational_unit.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "entity_name"),
					resource.TestCheckResourceAttrSet(rName, "entity_type"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPolicyAttach_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_policy" "test" {
  name        = "%[1]s"
  type        = "service_control_policy"
  description = "Created by terraform script"
  content     = jsonencode(
    {
      Version : "5.0",
      Statement : [
        {
          Effect : "Deny",
          Action : []
        }
      ]
    }
  )
}

data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%[1]s"
  parent_id = data.huaweicloud_organizations_organization.test.root_id
}
`, name)
}

func testAccPolicyAttach_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_organizations_policy_attach" "test" {
  policy_id = huaweicloud_organizations_policy.test.id
  entity_id = huaweicloud_organizations_organizational_unit.test.id
}
`, testAccPolicyAttach_base(name))
}
