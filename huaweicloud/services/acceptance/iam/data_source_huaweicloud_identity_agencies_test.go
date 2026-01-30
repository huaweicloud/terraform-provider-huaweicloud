package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAgencies_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identity_agencies.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_identity_agencies.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byTrustDomainId   = "data.huaweicloud_identity_agencies.filter_by_trust_domain_id"
		dcByTrustDomainId = acceptance.InitDataSourceCheck(byTrustDomainId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPrecheckDomainId(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAgencies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "agencies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "agencies.0.id"),
					resource.TestCheckResourceAttrSet(all, "agencies.0.name"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByTrustDomainId.CheckResourceExists(),
					resource.TestCheckOutput("is_trust_domain_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataAgencies_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%[1]s"
  description           = "This is a test agency"
  delegated_domain_name = "%[2]s"
  duration              = "ONEDAY"

  project_role {
    project = "%[3]s"
    roles   = ["CCE Administrator"]
  }
}
`, name, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME)
}

func testAccDataAgencies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# All
data "huaweicloud_identity_agencies" "all" {
  # The agency creation is not immediate, so the data source query may not get the latest agency.
  depends_on = [huaweicloud_identity_agency.test]
}

# Filter by name
locals {
  name = "%[2]s"
}

data "huaweicloud_identity_agencies" "filter_by_name" {
  name = local.name

  # The agency creation is not immediate, so the data source query may not get the latest agency.
  depends_on = [huaweicloud_identity_agency.test]
}

locals {
  name_filter_result = [for o in data.huaweicloud_identity_agencies.filter_by_name.agencies : o.name == local.name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) >= 1 && alltrue(local.name_filter_result)
}

# Filter by trust domain ID
locals {
  trust_domain_id = "%[3]s"
}

data "huaweicloud_identity_agencies" "filter_by_trust_domain_id" {
  trust_domain_id = local.trust_domain_id

  # The agency creation is not immediate, so the data source query may not get the latest agency.
  depends_on = [huaweicloud_identity_agency.test]
}

locals {
  trust_domain_id_filter_result = [for o in data.huaweicloud_identity_agencies.filter_by_trust_domain_id.agencies
  : o.trust_domain_id == local.trust_domain_id]
}

output "is_trust_domain_id_filter_useful" {
  value = length(local.trust_domain_id_filter_result) >= 1 && alltrue(local.trust_domain_id_filter_result)
}
`, testAccDataAgencies_base(name), name, acceptance.HW_DOMAIN_ID)
}
