package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5Agencies_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_agencies.test"
	rName := acceptance.RandomAccResourceName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5Agencies_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.#"),
				),
			},
			{
				Config: testAccDataSourceIdentityV5AgenciesWithId_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "agencies.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "agencies.0.agency_name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "agencies.0.description", "test for terraform"),
				),
			},
		},
	})
}

var testAccDataSourceIdentityV5Agencies_basic = `
data "huaweicloud_identityv5_agencies" "test" {}
`

func testAccDataSourceIdentityV5AgenciesWithId_basic(trustAgencyName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_service_agency" "test" {
  name                   = "%s"
  delegated_service_name = "service.APIG"
  policy_names           = ["NATReadOnlyPolicy"]
  description            = "test for terraform"
}

data "huaweicloud_identityv5_agencies" "test" {
  agency_id = huaweicloud_identity_service_agency.test.id
}
`, trustAgencyName)
}
