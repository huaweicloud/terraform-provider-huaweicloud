package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataPolicyAttachedEntities_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_organizations_policy_attached_entities.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPolicyAttachedEntities_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "attached_entities.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "attached_entities.0.id"),
					resource.TestCheckResourceAttrSet(all, "attached_entities.0.type"),
					resource.TestCheckResourceAttrSet(all, "attached_entities.0.name"),
				),
			},
		},
	})
}

func testAccDataPolicyAttachedEntities_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%[1]s"
  parent_id = data.huaweicloud_organizations_organization.test.root_id
}

resource "huaweicloud_organizations_policy" "test" {
  name        = "%[1]s"
  type        = "service_control_policy"
  description = "Created by terraform script"
  content     = jsonencode({
    Version = "5.0"
    Statement = [
      {
        Effect = "Deny"
        Action = []
      }
    ]
  })
}

resource "huaweicloud_organizations_policy_attach" "test" {
  policy_id = huaweicloud_organizations_policy.test.id
  entity_id = huaweicloud_organizations_organizational_unit.test.id
}
`, name)
}

func testAccDataPolicyAttachedEntities_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_organizations_policy_attached_entities" "test" {
  policy_id  = huaweicloud_organizations_policy.test.id
  depends_on = [huaweicloud_organizations_policy_attach.test]
}
`, testAccDataPolicyAttachedEntities_base(name))
}
