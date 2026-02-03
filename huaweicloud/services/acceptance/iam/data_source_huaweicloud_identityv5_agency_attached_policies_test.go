package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccDataV5AgencyAttachedPolicies_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identityv5_agency_attached_policies.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5AgencyAttachedPolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "attached_policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "attached_policies.0.policy_id"),
					resource.TestCheckResourceAttrSet(all, "attached_policies.0.policy_name"),
					resource.TestCheckResourceAttrSet(all, "attached_policies.0.urn"),
					resource.TestCheckResourceAttrSet(all, "attached_policies.0.attached_at"),
				),
			},
		},
	})
}

func testAccDataV5AgencyAttachedPolicies_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_trust_agency" "test" {
  name         = "%[1]s"
  policy_names = ["NATReadOnlyPolicy"]
  trust_policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "sts:agencies:assume",
          ]
          Effect = "Allow"
          Principal = {
            Service = [
              "service.OBS",
            ]
          }
        },
      ]
      Version = "5.0"
    }
  )
}
`, name)
}

func testAccDataV5AgencyAttachedPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identityv5_agency_attached_policies" "test" {
  agency_id = huaweicloud_identity_trust_agency.test.id
}
`, testAccDataV5AgencyAttachedPolicies_base(name))
}
