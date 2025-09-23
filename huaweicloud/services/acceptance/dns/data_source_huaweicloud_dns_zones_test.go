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
		name  = fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))
		rName = "huaweicloud_dns_zone.test.0"

		all           = "data.huaweicloud_dns_zones.test"
		dcForAllZones = acceptance.InitDataSourceCheck(all)

		byZoneId           = "data.huaweicloud_dns_zones.filter_by_zone_id"
		dcByZoneId         = acceptance.InitDataSourceCheck(byZoneId)
		byNotFoundZoneId   = "data.huaweicloud_dns_zones.filter_by_not_found_zone_id"
		dcByNotFoundZoneId = acceptance.InitDataSourceCheck(byNotFoundZoneId)

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

		byRouterId           = "data.huaweicloud_dns_zones.filter_by_router_id"
		dcByRouterId         = acceptance.InitDataSourceCheck(byRouterId)
		byNotFoundRouterId   = "data.huaweicloud_dns_zones.filter_by_not_found_router_id"
		dcByNotFoundRouterId = acceptance.InitDataSourceCheck(byNotFoundRouterId)

		bySortAsc    = "data.huaweicloud_dns_zones.filter_by_sort_asc"
		dcBySortAsc  = acceptance.InitDataSourceCheck(bySortAsc)
		bySortDesc   = "data.huaweicloud_dns_zones.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)
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
					// Filter by zone ID.
					dcByZoneId.CheckResourceExists(),
					resource.TestCheckOutput("is_zone_id_filter_useful", "true"),
					dcByNotFoundZoneId.CheckResourceExists(),
					resource.TestCheckOutput("zone_id_not_found_validation_pass", "true"),
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
					// Filter by the ID of the VPC associated with the private zone.
					dcByRouterId.CheckResourceExists(),
					resource.TestCheckOutput("is_router_id_filter_useful", "true"),
					dcByNotFoundRouterId.CheckResourceExists(),
					resource.TestCheckOutput("router_id_not_found_validation_pass", "true"),
					// Check ascending and descending results.
					dcBySortAsc.CheckResourceExists(),
					dcBySortDesc.CheckResourceExists(),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
					// Check attributes
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.id", rName, "id"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.name", rName, "name"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.description", rName, "description"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.email", rName, "email"),
					resource.TestCheckResourceAttr(byZoneId, "zones.0.zone_type", "private"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.enterprise_project_id", rName, "enterprise_project_id"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.routers.0.router_id", rName, "router.0.router_id"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.routers.0.router_region", rName, "router.0.router_region"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.ttl", rName, "ttl"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.proxy_pattern", rName, "proxy_pattern"),
					resource.TestMatchResourceAttr(byZoneId, "zones.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byZoneId, "zones.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// The zone resource does not have 'record_num' and 'pool_id' parameters.
					resource.TestCheckResourceAttrSet(byZoneId, "zones.0.record_num"),
					resource.TestCheckResourceAttrSet(byZoneId, "zones.0.pool_id"),
					// The 'masters' parameter has no corresponding scenario and is an empty list.
					resource.TestCheckResourceAttr(byZoneId, "zones.0.masters.#", "0"),
					// The 'masters' parameter has no corresponding scenario and is an empty list.
				),
			},
		},
	})
}

func testAccDataZones_base(name, randomId string) string {
	return fmt.Sprintf(`
data "huaweicloud_dns_zones" "test" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
}

# Filter by zone ID.
locals {
  zone_id = huaweicloud_dns_zone.test.0.id
}

data "huaweicloud_dns_zones" "filter_by_zone_id" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
  zone_id   = local.zone_id
}

locals {
  zone_id_filter_result = [for v in data.huaweicloud_dns_zones.filter_by_zone_id.zones[*].id : v == local.zone_id]
}

output "is_zone_id_filter_useful" {
  value = length(local.zone_id_filter_result) > 0 && alltrue(local.zone_id_filter_result)
}

data "huaweicloud_dns_zones" "filter_by_not_found_zone_id" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
  zone_id   = "%[2]s"
}

output "zone_id_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_zones.filter_by_not_found_zone_id.zones) == 0
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

# Filter by zone tags.
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

# Filter in ascending order using 'name' field as key.
data "huaweicloud_dns_zones" "filter_by_sort_asc" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
  sort_key  = "name"
  sort_dir  = "asc"
}

# Filter in descending order using 'name' field as key.
data "huaweicloud_dns_zones" "filter_by_sort_desc" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
  sort_key  = "name"
  sort_dir  = "desc"
}

locals {
  sort_desc_filter_result = data.huaweicloud_dns_zones.filter_by_sort_desc.zones
  sort_asc_first_name     = data.huaweicloud_dns_zones.filter_by_sort_asc.zones[0].name
  sort_desc_last_name     = data.huaweicloud_dns_zones.filter_by_sort_desc.zones[length(local.sort_desc_filter_result) - 1].name
}

output "sort_filter_is_useful" {
  value = local.sort_asc_first_name == local.sort_desc_last_name
}
`, name, randomId)
}

func testAccDataZones_basic(name string) string {
	randomId, _ := uuid.GenerateUUID()
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

# Filter by router Id.
locals {
  router_id = try(tolist(huaweicloud_dns_zone.test.0.router)[0].router_id, "")
}

data "huaweicloud_dns_zones" "filter_by_router_id" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
  router_id = local.router_id
}

locals {
  router_id_filter_result = [for v in flatten(data.huaweicloud_dns_zones.filter_by_router_id.zones[*].routers[*].router_id) :
  v == local.router_id]
}

output "is_router_id_filter_useful" {
  value = alltrue(local.router_id_filter_result) && length(local.router_id_filter_result) > 0
}

data "huaweicloud_dns_zones" "filter_by_not_found_router_id" {
  zone_type = huaweicloud_dns_zone.test.0.zone_type
  router_id = "%[5]s"
}

output "router_id_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_zones.filter_by_not_found_router_id.zones) == 0
}
`, acceptance.RandomAccResourceName(), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
		testAccDataZones_base(name, randomId), randomId)
}

func TestAccDataZones_public(t *testing.T) {
	var (
		name = fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

		rName = "huaweicloud_dns_zone.test.0"

		all           = "data.huaweicloud_dns_zones.test"
		dcForAllZones = acceptance.InitDataSourceCheck(all)

		byZoneId           = "data.huaweicloud_dns_zones.filter_by_zone_id"
		dcByZoneId         = acceptance.InitDataSourceCheck(byZoneId)
		byNotFoundZoneId   = "data.huaweicloud_dns_zones.filter_by_not_found_zone_id"
		dcByNotFoundZoneId = acceptance.InitDataSourceCheck(byNotFoundZoneId)

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

		bySortAsc    = "data.huaweicloud_dns_zones.filter_by_sort_asc"
		dcBySortAsc  = acceptance.InitDataSourceCheck(bySortAsc)
		bySortDesc   = "data.huaweicloud_dns_zones.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)
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
					// Filter by zone ID.
					dcByZoneId.CheckResourceExists(),
					resource.TestCheckOutput("is_zone_id_filter_useful", "true"),
					dcByNotFoundZoneId.CheckResourceExists(),
					resource.TestCheckOutput("zone_id_not_found_validation_pass", "true"),
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
					// Check ascending and descending results.
					dcBySortAsc.CheckResourceExists(),
					dcBySortDesc.CheckResourceExists(),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
					// Check attributes
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.id", rName, "id"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.name", rName, "name"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.description", rName, "description"),
					resource.TestCheckResourceAttr(byZoneId, "zones.0.zone_type", "public"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.email", rName, "email"),
					resource.TestCheckResourceAttrPair(byZoneId, "zones.0.ttl", rName, "ttl"),
					resource.TestMatchResourceAttr(byZoneId, "zones.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byZoneId, "zones.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// The zone resource does not have 'record_num' and 'pool_id' parameters.
					resource.TestCheckResourceAttrSet(byZoneId, "zones.0.record_num"),
					resource.TestCheckResourceAttrSet(byZoneId, "zones.0.pool_id"),
					// The 'masters' parameter has no corresponding scenario and is an empty list.
					resource.TestCheckResourceAttr(byZoneId, "zones.0.masters.#", "0"),
				),
			},
		},
	})
}

func testAccDataZones_public(name string) string {
	randomId, _ := uuid.GenerateUUID()
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
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, testAccDataZones_base(name, randomId))
}
