package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataHistoryTransactions_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_history_transactions.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByOrderBy   = "data.huaweicloud_das_history_transactions.filter_by_order_by"
		dcFilterByOrderBy = acceptance.InitDataSourceCheck(filterByOrderBy)

		filterByOrderField   = "data.huaweicloud_das_history_transactions.filter_by_order_field"
		dcFilterByOrderField = acceptance.InitDataSourceCheck(filterByOrderField)

		filterByLastSecMin   = "data.huaweicloud_das_history_transactions.filter_by_last_sec_min"
		dcFilterByLastSecMin = acceptance.InitDataSourceCheck(filterByLastSecMin)

		filterByLastSecMax   = "data.huaweicloud_das_history_transactions.filter_by_last_sec_max"
		dcFilterByLastSecMax = acceptance.InitDataSourceCheck(filterByLastSecMax)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataHistoryTransactions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "transactions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "transactions.0.last_sec"),
					resource.TestCheckResourceAttrSet(all, "transactions.0.wait_locks"),
					resource.TestCheckResourceAttrSet(all, "transactions.0.hold_locks"),
					resource.TestCheckResourceAttrSet(all, "transactions.0.detail"),
					resource.TestMatchResourceAttr(all, "transactions.0.occurrence_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					// filter by order_by
					dcFilterByOrderBy.CheckResourceExists(),
					resource.TestCheckOutput("is_order_by_filter_useful", "true"),

					// filter by order_field
					dcFilterByOrderField.CheckResourceExists(),
					resource.TestCheckOutput("is_order_field_filter_useful", "true"),

					// filter by last_sec_min
					dcFilterByLastSecMin.CheckResourceExists(),
					resource.TestCheckOutput("is_last_sec_min_filter_useful", "true"),

					// filter by last_sec_max
					dcFilterByLastSecMax.CheckResourceExists(),
					resource.TestCheckOutput("is_last_sec_max_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataHistoryTransactions_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccDataHistoryTransactions_basic() string {
	return fmt.Sprintf(`
%[1]s

// Without any filter parameters.
data "huaweicloud_das_history_transactions" "all" {
  instance_id    = local.instance_ids[0]
  datastore_type = "MySQL"
  start_time     = "2000-01-01T00:00:00+08:00"
  end_time       = "2099-01-01T00:00:00+08:00"
}

# Filter by order_by
locals {
  order_by = "desc"
}

data "huaweicloud_das_history_transactions" "filter_by_order_by" {
  instance_id    = local.instance_ids[0]
  datastore_type = "MySQL"
  start_time     = "2000-01-01T00:00:00+08:00"
  end_time       = "2099-01-01T00:00:00+08:00"
  order_by       = local.order_by
}

locals {
  order_by_filter_result = [
    for i, v in data.huaweicloud_das_history_transactions.filter_by_order_by.transactions : true
  ]
}

output "is_order_by_filter_useful" {
  value = length(local.order_by_filter_result) > 0 && alltrue(local.order_by_filter_result)
}

# Filter by order_field
locals {
  order_field = "lastSec"
}

data "huaweicloud_das_history_transactions" "filter_by_order_field" {
  instance_id    = local.instance_ids[0]
  datastore_type = "MySQL"
  start_time     = "2000-01-01T00:00:00+08:00"
  end_time       = "2099-01-01T00:00:00+08:00"
  order_field    = local.order_field
}

locals {
  order_field_filter_result = [
    for i, v in data.huaweicloud_das_history_transactions.filter_by_order_field.transactions : true
  ]
}

output "is_order_field_filter_useful" {
  value = length(local.order_field_filter_result) > 0 && alltrue(local.order_field_filter_result)
}

# Filter by last_sec_min
locals {
  last_sec_min = 1
}

data "huaweicloud_das_history_transactions" "filter_by_last_sec_min" {
  instance_id    = local.instance_ids[0]
  datastore_type = "MySQL"
  start_time     = "2000-01-01T00:00:00+08:00"
  end_time       = "2099-01-01T00:00:00+08:00"
  last_sec_min   = local.last_sec_min
}

locals {
  last_sec_min_filter_result = [
    for v in data.huaweicloud_das_history_transactions.filter_by_last_sec_min.transactions : v.last_sec >= local.last_sec_min
  ]
}

output "is_last_sec_min_filter_useful" {
  value = length(local.last_sec_min_filter_result) > 0 && alltrue(local.last_sec_min_filter_result)
}

# Filter by last_sec_max
locals {
  last_sec_max = 3600
}

data "huaweicloud_das_history_transactions" "filter_by_last_sec_max" {
  instance_id    = local.instance_ids[0]
  datastore_type = "MySQL"
  start_time     = "2000-01-01T00:00:00+08:00"
  end_time       = "2099-01-01T00:00:00+08:00"
  last_sec_max   = local.last_sec_max
}

locals {
  last_sec_max_filter_result = [
    for v in data.huaweicloud_das_history_transactions.filter_by_last_sec_max.transactions : v.last_sec <= local.last_sec_max
  ]
}

output "is_last_sec_max_filter_useful" {
  value = length(local.last_sec_max_filter_result) > 0 && alltrue(local.last_sec_max_filter_result)
}
`, testAccDataHistoryTransactions_base())
}
