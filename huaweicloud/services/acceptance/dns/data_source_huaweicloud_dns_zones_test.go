package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataZones_basic(t *testing.T) {
	var (
		name = fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

		all           = "data.huaweicloud_dns_zones.test"
		dcForAllZones = acceptance.InitDataSourceCheck(all)

		byNameFuzzy      = "data.huaweicloud_dns_zones.filter_by_name_fuzzy"
		dcByNameFuzzy    = acceptance.InitDataSourceCheck(byNameFuzzy)
		byNameExact      = "data.huaweicloud_dns_zones.filter_by_name_exact"
		dcByNameExact    = acceptance.InitDataSourceCheck(byNameExact)
		byNotFoundName   = "data.huaweicloud_dns_zones.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byStatus   = "data.huaweicloud_dns_zones.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEpsId           = "data.huaweicloud_dns_zones.filter_by_eps_id"
		dcByEpsId         = acceptance.InitDataSourceCheck(byEpsId)
		byNotFoundEpsId   = "data.huaweicloud_dns_zones.filter_by_not_found_eps_id"
		dcByNotFoundEpsId = acceptance.InitDataSourceCheck(byNotFoundEpsId)

		byTags           = "data.huaweicloud_dns_zones.filter_by_tags"
		dcByTags         = acceptance.InitDataSourceCheck(byTags)
		byNotFoundTags   = "data.huaweicloud_dns_zones.filter_by_not_found_tags"
		dcByNotFoundTags = acceptance.InitDataSourceCheck(byNotFoundTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataZones_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcForAllZones.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "zones.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Fuzzy search by zone name.
					dcByNameFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_name_fuzzy_filter_useful", "true"),
					// Exactly search by zone name, but the name can be repeated for the private zone.
					dcByNameExact.CheckResourceExists(),
					resource.TestCheckOutput("is_exact_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
					// Filter by zone status.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Filter by zone enterprise project ID.
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					dcByNotFoundEpsId.CheckResourceExists(),
					resource.TestCheckOutput("eps_id_not_found_validation_pass", "true"),
					dcByTags.CheckResourceExists(),
					// Filter by zone tags.
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					dcByNotFoundTags.CheckResourceExists(),
					resource.TestCheckOutput("tags_not_found_validation_pass", "true"),
					// Check attributes
					// Cannot query unique data by name.
					// A VPC cannot be bound to a private network with the same name, so after the router_id query parameter is provided,
					// the router_id and name can be used together to query and determine the unique data.
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.id"),
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.name"),
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.description"),
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.email"),
					resource.TestCheckResourceAttr(byNameExact, "zones.0.zone_type", "private"),
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.routers.0.router_id"),
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.routers.0.router_region"),
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.ttl"),
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.record_num"),
					resource.TestCheckResourceAttr(byNameExact, "zones.0.masters.#", "0"),
				),
			},
		},
	})
}

func testAccDataZones_base(name string) string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dns_zones" "test" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
}

# Filter by zone name (Fuzzy search).
data "huaweicloud_dns_zones" "filter_by_name_fuzzy" {
  depends_on = [huaweicloud_dns_zone.test]

  zone_type = huaweicloud_dns_zone.test.0.zone_type
  name      = "%[1]s"
}

locals {
  name_fuzzy_filter_result = [for v in data.huaweicloud_dns_zones.filter_by_name_fuzzy.zones[*].name :
  strcontains(v, "%[1]s")]
}

output "is_name_fuzzy_filter_useful" {
  value = length(local.name_fuzzy_filter_result) >= 2 && alltrue(local.name_fuzzy_filter_result)
}

# Filter by zone name (exact search).
locals {
  zone_name = huaweicloud_dns_zone.test.0.name
}

data "huaweicloud_dns_zones" "filter_by_name_exact" {
  depends_on = [huaweicloud_dns_zone.test]

  zone_type   = huaweicloud_dns_zone.test.0.zone_type
  name        = local.zone_name
  search_mode = "equal"
}

locals {
  exact_name_filter_result = [for v in data.huaweicloud_dns_zones.filter_by_name_exact.zones[*].name : v == local.zone_name]
}

# The name can be repeated for the private zone.
output "is_exact_name_filter_useful" {
  value = alltrue(local.exact_name_filter_result) && length(local.exact_name_filter_result) >= 1
}

data "huaweicloud_dns_zones" "filter_by_not_found_name" {
  depends_on = [huaweicloud_dns_zone.test]

  zone_type = huaweicloud_dns_zone.test.0.zone_type
  name      = "zone_name_not_found"
}

output "name_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_zones.filter_by_not_found_name.zones) == 0
}

# By filter zone status.
data "huaweicloud_dns_zones" "filter_by_status" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
  # In the corresponding resource, the value of status is ENABLE.
  status = "ACTIVE"
}

locals {
  status_filter_result = [for v in data.huaweicloud_dns_zones.filter_by_status.zones[*].status : v == "ACTIVE"]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by enterprise project ID.
locals {
  enterprise_project_id = huaweicloud_dns_zone.test.0.enterprise_project_id
}

data "huaweicloud_dns_zones" "filter_by_eps_id" {
  zone_type             = huaweicloud_dns_zone.test.0.zone_type
  enterprise_project_id = local.enterprise_project_id
}

locals {
  eps_id_filter_result = [for v in data.huaweicloud_dns_zones.filter_by_eps_id.zones[*].enterprise_project_id :
  v == local.enterprise_project_id]
}

output "is_eps_id_filter_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}

data "huaweicloud_dns_zones" "filter_by_not_found_eps_id" {
  zone_type             = huaweicloud_dns_zone.test.0.zone_type
  enterprise_project_id = "%[2]s"
}

output "eps_id_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_zones.filter_by_not_found_eps_id.zones) == 0
}

# Filter by recordset tags.
data "huaweicloud_dns_zones" "filter_by_tags" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
  tags      = join("|", [for key, value in huaweicloud_dns_zone.test.0.tags : format("%%v,%%v", key, value)])
}

output "is_tags_filter_useful" {
  value = length(data.huaweicloud_dns_zones.filter_by_tags.zones) >= 2
}

data "huaweicloud_dns_zones" "filter_by_not_found_tags" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
  tags      = "not_found_tags,not_found"
}

output "tags_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_zones.filter_by_not_found_tags.zones) == 0
}
`, name, randomId)
}

func testAccDataZones_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_dns_zone" "test" {
  count = 2

  name                  = "${count.index}%[2]s"
  email                 = "email@example.com"
  description           = "Used to query zones"
  zone_type             = "private"
  enterprise_project_id = "%[3]s"

  router {
    router_id = huaweicloud_vpc.test.id
  }

  tags = {
    zone_type = "private"
    owner     = "terraform"
  }
}

%[4]s
`, acceptance.RandomAccResourceName(), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, testAccDataZones_base(name))
}

func TestAccDataZones_public(t *testing.T) {
	var (
		name = fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

		rName = "huaweicloud_dns_zone.test.0"

		all           = "data.huaweicloud_dns_zones.test"
		dcForAllZones = acceptance.InitDataSourceCheck(all)

		byNameFuzzy      = "data.huaweicloud_dns_zones.filter_by_name_fuzzy"
		dcByNameFuzzy    = acceptance.InitDataSourceCheck(byNameFuzzy)
		byNameExact      = "data.huaweicloud_dns_zones.filter_by_name_exact"
		dcByNameExact    = acceptance.InitDataSourceCheck(byNameExact)
		byNotFoundName   = "data.huaweicloud_dns_zones.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byStatus   = "data.huaweicloud_dns_zones.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEpsId           = "data.huaweicloud_dns_zones.filter_by_eps_id"
		dcByEpsId         = acceptance.InitDataSourceCheck(byEpsId)
		byNotFoundEpsId   = "data.huaweicloud_dns_zones.filter_by_not_found_eps_id"
		dcByNotFoundEpsId = acceptance.InitDataSourceCheck(byNotFoundEpsId)

		byTags           = "data.huaweicloud_dns_zones.filter_by_tags"
		dcByTags         = acceptance.InitDataSourceCheck(byTags)
		byNotFoundTags   = "data.huaweicloud_dns_zones.filter_by_not_found_tags"
		dcByNotFoundTags = acceptance.InitDataSourceCheck(byNotFoundTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataZones_public(name),
				Check: resource.ComposeTestCheckFunc(
					dcForAllZones.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "zones.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Fuzzy search by zone name.
					dcByNameFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_name_fuzzy_filter_useful", "true"),
					// Exactly search by zone name, and the name must be unique.
					dcByNameExact.CheckResourceExists(),
					resource.TestCheckOutput("is_exact_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
					// Filter by zone status.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Filter by zone enterprise project ID.
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					dcByNotFoundEpsId.CheckResourceExists(),
					resource.TestCheckOutput("eps_id_not_found_validation_pass", "true"),
					dcByTags.CheckResourceExists(),
					// Filter by zone tags.
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					dcByNotFoundTags.CheckResourceExists(),
					resource.TestCheckOutput("tags_not_found_validation_pass", "true"),
					// Check attributes
					resource.TestCheckResourceAttrPair(byNameExact, "zones.0.id", rName, "id"),
					resource.TestCheckResourceAttrPair(byNameExact, "zones.0.name", rName, "name"),
					resource.TestCheckResourceAttrPair(byNameExact, "zones.0.description", rName, "description"),
					resource.TestCheckResourceAttr(byNameExact, "zones.0.zone_type", "public"),
					resource.TestCheckResourceAttrPair(byNameExact, "zones.0.email", rName, "email"),
					resource.TestCheckResourceAttrPair(byNameExact, "zones.0.ttl", rName, "ttl"),
					// The zone resource does not have 'record_num' parameter.
					resource.TestCheckResourceAttrSet(byNameExact, "zones.0.record_num"),
					// The 'masters' parameter has no corresponding scenario and is an empty list.
					resource.TestCheckResourceAttr(byNameExact, "zones.0.masters.#", "0"),
				),
			},
		},
	})
}

func testAccDataZones_public(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "test" {
  count = 2

  name                  = "${count.index}%[1]s"
  description           = "Used to query zones"
  ttl                   = 300
  status                = "DISABLE"
  enterprise_project_id = "%[2]s"

  tags = {
    zone_type = "public"
    owner     = "terraform"
  }
}

%[3]s
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, testAccDataZones_base(name))
}
