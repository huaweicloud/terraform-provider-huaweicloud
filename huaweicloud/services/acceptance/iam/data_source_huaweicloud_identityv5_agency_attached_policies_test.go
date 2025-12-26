package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5AgencyAttachedPolicies_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_agency_attached_policies.test"
	rName := acceptance.RandomAccResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5AgencyAttachedPolicies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "agency_id"),
					resource.TestCheckResourceAttr(dataSourceName, "attached_policies.0.policy_id", "NATReadOnlyPolicy"),
					resource.TestCheckResourceAttr(dataSourceName, "attached_policies.0.policy_name", "NATReadOnlyPolicy"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5AgencyAttachedPolicies_basic(agencyName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identityv5_agency_attached_policies" "test" {
  agency_id = huaweicloud_identity_trust_agency.test.id
}
`, testAccIdentityTrustAgency_basic(agencyName))
}
