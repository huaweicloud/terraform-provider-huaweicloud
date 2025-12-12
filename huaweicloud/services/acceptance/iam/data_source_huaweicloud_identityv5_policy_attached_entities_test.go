package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5PolicyAttachedEntities_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_policy_attached_entities.test"
	rName := acceptance.RandomAccResourceName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5PolicyAttachedEntities_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_users.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_groups.#"),
					resource.TestCheckResourceAttr(dataSourceName, "policy_agencies.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5PolicyAttachedEntities_basic(trustAgencyName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_identityv5_policy_attached_entities" "test" {
  policy_id = data.huaweicloud_identityv5_policies.test.policies[0].policy_id

  depends_on = [
	huaweicloud_identity_trust_agency.test
  ]
}
`, testAccIdentityTrustAgency_basic(trustAgencyName), testAccDataSourceIdentityV5Policies_basic)
}
