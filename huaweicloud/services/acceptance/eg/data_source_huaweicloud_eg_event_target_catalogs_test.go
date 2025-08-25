package eg

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEventTargetCatalogs_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_eg_event_target_catalogs.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byFuzzyLabel   = "data.huaweicloud_eg_event_target_catalogs.filter_by_fuzzy_label"
		dcByFuzzyLabel = acceptance.InitDataSourceCheck(byFuzzyLabel)

		bySupportTypes   = "data.huaweicloud_eg_event_target_catalogs.filter_by_support_types"
		dcBySupportTypes = acceptance.InitDataSourceCheck(bySupportTypes)

		bySortDesc   = "data.huaweicloud_eg_event_target_catalogs.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)

		bySortAsc   = "data.huaweicloud_eg_event_target_catalogs.filter_by_sort_asc"
		dcBySortAsc = acceptance.InitDataSourceCheck(bySortAsc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEventTargetCatalogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "catalogs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.id"),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.name"),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.label"),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.provider_type"),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.support_types.#"),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.created_time"),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.updated_time"),
					dcByFuzzyLabel.CheckResourceExists(),
					resource.TestCheckOutput("is_fuzzy_label_filter_useful", "true"),
					dcBySupportTypes.CheckResourceExists(),
					resource.TestCheckOutput("is_support_types_filter_useful", "true"),
					dcBySortDesc.CheckResourceExists(),
					dcBySortAsc.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_filter_useful", "true"),
					// description and parameters may be empty, so we don't check them.
				),
			},
		},
	})
}

func testAccDataEventTargetCatalogs_basic() string {
	return `
data "huaweicloud_eg_event_target_catalogs" "test" {}

locals {
  label         = try(lookup(data.huaweicloud_eg_event_target_catalogs.test.catalogs[0], "label", null), null)
  support_types = try([for v in data.huaweicloud_eg_event_target_catalogs.test.catalogs[*].support_types : v if length(v) >= 2][0], [])
}

# Filter by fuzzy label.
data "huaweicloud_eg_event_target_catalogs" "filter_by_fuzzy_label" {
  fuzzy_label = local.label
}

locals {
  fuzzy_label_filter_result = [for v in data.huaweicloud_eg_event_target_catalogs.filter_by_fuzzy_label.catalogs[*].label :
  strcontains(v, local.label)]
}

output "is_fuzzy_label_filter_useful" {
  value = length(local.fuzzy_label_filter_result) > 0 && alltrue(local.fuzzy_label_filter_result)
}

# Filter by support types.
data "huaweicloud_eg_event_target_catalogs" "filter_by_support_types" {
  support_types = local.support_types
}

locals {
  support_types_filter_result = [for v in data.huaweicloud_eg_event_target_catalogs.filter_by_support_types.catalogs[*].support_types :
  alltrue([for item in local.support_types : contains(v, item)])]
}

output "is_support_types_filter_useful" {
  value = length(local.support_types_filter_result) > 0 && alltrue(local.support_types_filter_result)
}

# Filter by sort.
data "huaweicloud_eg_event_target_catalogs" "filter_by_sort_desc" {
  sort = "created_time:DESC"
}

data "huaweicloud_eg_event_target_catalogs" "filter_by_sort_asc" {
  sort = "created_time:ASC"
}

locals {
  sort_desc_filter_result = data.huaweicloud_eg_event_target_catalogs.filter_by_sort_desc.catalogs[*].created_time
  sort_asc_filter_result  = data.huaweicloud_eg_event_target_catalogs.filter_by_sort_asc.catalogs[*].created_time
}

output "is_sort_filter_useful" {
  value = (
    length(local.sort_desc_filter_result) == length(local.sort_asc_filter_result) &&
    try(local.sort_desc_filter_result[0], "") == try(element(local.sort_asc_filter_result, -1), "")
  )
}
`
}
