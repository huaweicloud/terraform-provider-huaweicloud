package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureBusinessMetrics_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_business_metrics.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByName   = "data.huaweicloud_dataarts_architecture_business_metrics.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByCreateBy   = "data.huaweicloud_dataarts_architecture_business_metrics.filter_by_create_by"
		dcFilterByCreateBy = acceptance.InitDataSourceCheck(filterByCreateBy)

		filterByOwner   = "data.huaweicloud_dataarts_architecture_business_metrics.filter_by_owner"
		dcFilterByOwner = acceptance.InitDataSourceCheck(filterByOwner)

		filterByStatus   = "data.huaweicloud_dataarts_architecture_business_metrics.filter_by_status"
		dcFilterByStatus = acceptance.InitDataSourceCheck(filterByStatus)

		filterByBizCatalogId   = "data.huaweicloud_dataarts_architecture_business_metrics.filter_by_biz_catalog_id"
		dcFilterByBizCatalogId = acceptance.InitDataSourceCheck(filterByBizCatalogId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsBizCatalogID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataArchitectureBusinessMetrics_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Architecture business metrics"),
			},
			{
				Config: testAccDataSourceArchitectureBusinessMetrics_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "metrics.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// filter by name
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(filterByName, "metrics.0.id"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.name",
						"huaweicloud_dataarts_architecture_business_metric.test", "name"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.biz_catalog_id",
						"huaweicloud_dataarts_architecture_business_metric.test", "biz_catalog_id"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.time_filters",
						"huaweicloud_dataarts_architecture_business_metric.test", "time_filters"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.interval_type",
						"huaweicloud_dataarts_architecture_business_metric.test", "interval_type"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.owner",
						"huaweicloud_dataarts_architecture_business_metric.test", "owner"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.owner_department",
						"huaweicloud_dataarts_architecture_business_metric.test", "owner_department"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.destination",
						"huaweicloud_dataarts_architecture_business_metric.test", "destination"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.definition",
						"huaweicloud_dataarts_architecture_business_metric.test", "definition"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.expression",
						"huaweicloud_dataarts_architecture_business_metric.test", "expression"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.description",
						"huaweicloud_dataarts_architecture_business_metric.test", "description"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.apply_scenario",
						"huaweicloud_dataarts_architecture_business_metric.test", "apply_scenario"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.technical_metric",
						"huaweicloud_dataarts_architecture_business_metric.test", "technical_metric"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.measure",
						"huaweicloud_dataarts_architecture_business_metric.test", "measure"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.general_filters",
						"huaweicloud_dataarts_architecture_business_metric.test", "general_filters"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.data_origin",
						"huaweicloud_dataarts_architecture_business_metric.test", "data_origin"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.unit",
						"huaweicloud_dataarts_architecture_business_metric.test", "unit"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.name_alias",
						"huaweicloud_dataarts_architecture_business_metric.test", "name_alias"),
					resource.TestCheckResourceAttrPair(filterByName, "metrics.0.code",
						"huaweicloud_dataarts_architecture_business_metric.test", "code"),
					resource.TestCheckResourceAttrSet(filterByName, "metrics.0.status"),
					resource.TestCheckResourceAttrSet(filterByName, "metrics.0.biz_catalog_path"),
					resource.TestCheckResourceAttrSet(filterByName, "metrics.0.created_by"),
					resource.TestCheckResourceAttrSet(filterByName, "metrics.0.updated_by"),
					resource.TestMatchResourceAttr(filterByName, "metrics.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$`)),
					resource.TestMatchResourceAttr(filterByName, "metrics.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$`)),
					resource.TestCheckResourceAttrSet(filterByName, "metrics.0.l1"),

					// filter by create by
					dcFilterByCreateBy.CheckResourceExists(),
					resource.TestCheckOutput("is_create_by_filter_useful", "true"),

					// filter by owner
					dcFilterByOwner.CheckResourceExists(),
					resource.TestCheckOutput("is_owner_filter_useful", "true"),

					// filter by status
					dcFilterByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					// filter by biz catalog ID
					dcFilterByBizCatalogId.CheckResourceExists(),
					resource.TestCheckOutput("is_biz_catalog_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureBusinessMetrics_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_business_metrics" "test" {
  workspace_id = "%[1]s"
}
`, randUUID.String())
}

func testAccDataSourceArchitectureBusinessMetrics_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_business_metric" "test" {
  workspace_id     = "%[1]s"
  biz_catalog_id   = "%[2]s"
  owner            = "test-owner"
  owner_department = "test-department"
  name             = "%[3]s"
  name_alias       = "addition_of_three_numbers"
  time_filters     = "双周"
  interval_type    = "HOUR"
  destination      = "test destination"
  definition       = "test definition"
  expression       = "a+b+c"
  description      = "Addition of three numbers"
  apply_scenario   = "test apply scenario"
  technical_metric = 100
  measure          = "test measure"
  general_filters  = "test general filters"
  data_origin      = "test data origin"
  unit             = "test unit"
  code             = "testCode"
}

data "huaweicloud_dataarts_architecture_business_metrics" "all" {
  depends_on = [
    huaweicloud_dataarts_architecture_business_metric.test,
  ]

  workspace_id = "%[1]s"
}

# Filter by name
locals {
  metric_name = huaweicloud_dataarts_architecture_business_metric.test.name
}

data "huaweicloud_dataarts_architecture_business_metrics" "filter_by_name" {
  depends_on = [
    huaweicloud_dataarts_architecture_business_metric.test,
  ]

  workspace_id = "%[1]s"
  name         = local.metric_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_business_metrics.filter_by_name.metrics : v.name == local.metric_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by create by
locals {
  create_by = huaweicloud_dataarts_architecture_business_metric.test.created_by
}

data "huaweicloud_dataarts_architecture_business_metrics" "filter_by_create_by" {
  depends_on = [
    huaweicloud_dataarts_architecture_business_metric.test,
  ]

  workspace_id = "%[1]s"
  create_by    = local.create_by
}

locals {
  create_by_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_business_metrics.filter_by_create_by.metrics : v.created_by == local.create_by
  ]
}

output "is_create_by_filter_useful" {
  value = length(local.create_by_filter_result) > 0 && alltrue(local.create_by_filter_result)
}

# Filter by owner
locals {
  metric_owner = huaweicloud_dataarts_architecture_business_metric.test.owner
}

data "huaweicloud_dataarts_architecture_business_metrics" "filter_by_owner" {
  depends_on = [
    huaweicloud_dataarts_architecture_business_metric.test,
  ]

  workspace_id = "%[1]s"
  owner        = local.metric_owner
}

locals {
  owner_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_business_metrics.filter_by_owner.metrics : v.owner == local.metric_owner
  ]
}

output "is_owner_filter_useful" {
  value = length(local.owner_filter_result) > 0 && alltrue(local.owner_filter_result)
}

# Filter by status
locals {
  metric_status = huaweicloud_dataarts_architecture_business_metric.test.status
}

data "huaweicloud_dataarts_architecture_business_metrics" "filter_by_status" {
  depends_on = [
    huaweicloud_dataarts_architecture_business_metric.test,
  ]

  workspace_id = "%[1]s"
  status       = local.metric_status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_business_metrics.filter_by_status.metrics : v.status == local.metric_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by biz catalog ID
locals {
  biz_catalog_id = huaweicloud_dataarts_architecture_business_metric.test.biz_catalog_id
}

data "huaweicloud_dataarts_architecture_business_metrics" "filter_by_biz_catalog_id" {
  depends_on = [
    huaweicloud_dataarts_architecture_business_metric.test,
  ]

  workspace_id   = "%[1]s"
  biz_catalog_id = local.biz_catalog_id
}

locals {
  biz_catalog_id_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_business_metrics.filter_by_biz_catalog_id.metrics : v.biz_catalog_id == local.biz_catalog_id
  ]
}

output "is_biz_catalog_id_filter_useful" {
  value = length(local.biz_catalog_id_filter_result) > 0 && alltrue(local.biz_catalog_id_filter_result)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_BIZ_CATALOG_ID, name)
}
