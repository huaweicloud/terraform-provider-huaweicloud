package iam_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccV5Agencies_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identityv5_agencies.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byAgencyId   = "data.huaweicloud_identityv5_agencies.filter_by_agency_id"
		dcByAgencyId = acceptance.InitDataSourceCheck(byAgencyId)

		byPathPrefix   = "data.huaweicloud_identityv5_agencies.filter_by_path_prefix"
		dcByPathPrefix = acceptance.InitDataSourceCheck(byPathPrefix)

		byAgencyIdWithoutPath   = "data.huaweicloud_identityv5_agencies.filter_by_agency_id_without_path"
		dcByAgencyIdWithoutPath = acceptance.InitDataSourceCheck(byAgencyIdWithoutPath)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckServiceLinkedAgencyPrincipal(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccV5Agencies_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "agencies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'agency_id' parameter.
					dcByAgencyId.CheckResourceExists(),
					resource.TestCheckOutput("is_agency_id_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byAgencyId, "agencies.0.agency_id",
						"huaweicloud_identityv5_service_linked_agency.test", "id"),
					resource.TestCheckResourceAttrPair(byAgencyId, "agencies.0.agency_name",
						"huaweicloud_identityv5_service_linked_agency.test", "agency_name"),
					resource.TestCheckResourceAttrPair(byAgencyId, "agencies.0.path",
						"huaweicloud_identityv5_service_linked_agency.test", "path"),
					resource.TestCheckResourceAttrPair(byAgencyId, "agencies.0.urn",
						"huaweicloud_identityv5_service_linked_agency.test", "urn"),
					resource.TestCheckResourceAttrPair(byAgencyId, "agencies.0.description",
						"huaweicloud_identityv5_service_linked_agency.test", "description"),
					resource.TestCheckResourceAttrPair(byAgencyId, "agencies.0.max_session_duration",
						"huaweicloud_identityv5_service_linked_agency.test", "max_session_duration"),
					resource.TestCheckResourceAttrSet(byAgencyId, "agencies.0.created_at"),
					// Filter by 'path_prefix' parameter.
					dcByPathPrefix.CheckResourceExists(),
					resource.TestCheckOutput("is_path_prefix_filter_useful", "true"),
					// Filter by 'agency_id' parameter without path.
					dcByAgencyIdWithoutPath.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byAgencyIdWithoutPath, "agencies.0.trust_domain_id"),
					resource.TestCheckResourceAttrPair(byAgencyIdWithoutPath, "agencies.0.trust_domain_name",
						"huaweicloud_identity_agency.test", "delegated_domain_name"),
				),
			},
		},
	})
}

func testAccV5Agencies_base() string {
	return fmt.Sprintf(`
# Create a service-linked agency, it will return the 'path' field.
resource "huaweicloud_identityv5_service_linked_agency" "test" {
  service_principal = "%[1]s"
  description       = "Create by terraform script"
}

# Create a agency, it will not return the 'path' field.
resource "huaweicloud_identity_agency" "test" {
  name                  = "%[2]s"
  delegated_domain_name = "%[3]s"
  description           = "Create by terraform script"
  duration              = "30"
}
`, acceptance.HW_IAM_SERVICE_LINKED_AGENCY_PRINCIPAL, acceptance.RandomAccResourceName(), acceptance.HW_DOMAIN_NAME)
}

func testAccV5Agencies_basic() string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_identityv5_agencies" "test" {
  depends_on = [huaweicloud_identityv5_service_linked_agency.test]
}

# Filter by 'agency_id' parameter.
locals {
  agency_id = huaweicloud_identityv5_service_linked_agency.test.id
}

data "huaweicloud_identityv5_agencies" "filter_by_agency_id" {
  agency_id = local.agency_id
}

locals {
  agency_id_filter_result = [for v in data.huaweicloud_identityv5_agencies.filter_by_agency_id.agencies[*].agency_id :
  v == local.agency_id]
}

output "is_agency_id_filter_useful" {
  value = length(local.agency_id_filter_result) > 0 && alltrue(local.agency_id_filter_result)
}

# Filter by 'path_prefix' parameter.
locals {
  path_prefix = huaweicloud_identityv5_service_linked_agency.test.path
}

data "huaweicloud_identityv5_agencies" "filter_by_path_prefix" {
  path_prefix = huaweicloud_identityv5_service_linked_agency.test.path
}

locals {
  path_prefix_filter_result = [for v in data.huaweicloud_identityv5_agencies.filter_by_path_prefix.agencies[*].path :
  strcontains(v, local.path_prefix)]
}

output "is_path_prefix_filter_useful" {
  value = length(local.path_prefix_filter_result) > 0 && alltrue(local.path_prefix_filter_result)
}

# It will not return the 'path' field, but will return the 'trust_domain_id' and 'trust_domain_name' fields.
data "huaweicloud_identityv5_agencies" "filter_by_agency_id_without_path" {
  agency_id = huaweicloud_identity_agency.test.id
}
`, testAccV5Agencies_base())
}
