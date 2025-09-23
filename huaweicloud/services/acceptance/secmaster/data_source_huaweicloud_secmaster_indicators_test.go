package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIndicators_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_indicators.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterIndicatorTypeID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceIndicators_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "indicators.#"),
					resource.TestCheckResourceAttrSet(dataSource, "indicators.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "indicators.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "indicators.0.threat_degree"),
					resource.TestCheckResourceAttrSet(dataSource, "indicators.0.type.#"),
					resource.TestCheckResourceAttrSet(dataSource, "indicators.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "indicators.0.confidence"),
					resource.TestCheckResourceAttrSet(dataSource, "indicators.0.data_class_id"),

					resource.TestCheckOutput("is_ids_filter_useful", "true"),
					resource.TestCheckOutput("is_data_class_id_filter_useful", "true"),
					resource.TestCheckOutput("is_condition_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceIndicators_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_secmaster_indicators" "test" {
  workspace_id = "%[2]s"

  depends_on = [
    huaweicloud_secmaster_indicator.test1,
    huaweicloud_secmaster_indicator.test2,
  ]
}

locals {
  id            = data.huaweicloud_secmaster_indicators.test.indicators[0].id
  name          = data.huaweicloud_secmaster_indicators.test.indicators[0].name
  status        = data.huaweicloud_secmaster_indicators.test.indicators[0].status
  data_class_id = data.huaweicloud_secmaster_indicators.test.indicators[0].data_class_id
}

data "huaweicloud_secmaster_indicators" "filter_by_id" {
  workspace_id = "%[2]s"
  ids          = [local.id]
}

data "huaweicloud_secmaster_indicators" "filter_by_data_class_id" {
  workspace_id  = "%[2]s"
  data_class_id = local.data_class_id
}

data "huaweicloud_secmaster_indicators" "filter_by_condition" {
  workspace_id = "%[2]s"

  condition {
    conditions {
      name = "name"
      data = ["name", "=", local.name]
    }

    conditions {
      name = "status"
      data = ["status", "=", local.status]
    }

    logics = ["name", "and", "status"]
  }
}

output "is_ids_filter_useful" {
  value = length(data.huaweicloud_secmaster_indicators.filter_by_id.indicators) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_indicators.filter_by_id.indicators[*].id : v == local.id]
  )
}

output "is_data_class_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_indicators.filter_by_data_class_id.indicators) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_indicators.filter_by_data_class_id.indicators[*].data_class_id : v == local.data_class_id]
  )
}

output "is_condition_filter_useful" {
  value = length(data.huaweicloud_secmaster_indicators.filter_by_condition.indicators) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_indicators.filter_by_condition.indicators : v.name == local.name && v.status == local.status]
  )
}
`, testIndicators_buildData(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testIndicators_buildData(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_indicator" "test1" {
  workspace_id = "%[1]s"
  name         = "%[2]s"

  type {
    category       = "Domain"
    indicator_type = "Domain"
    id             = "%[3]s"
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  status                = "Open"
  confidence            = "80"
  first_occurrence_time = "2024-08-21T17:23:55.000+08:00"
  last_occurrence_time  = "2024-08-22T11:15:30.000+08:00"
  threat_degree         = "Black"
  granularity           = "1"
  value                 = "test.terraform.com"
}

resource "huaweicloud_secmaster_indicator" "test2" {
  workspace_id = "%[1]s"
  name         = "%[2]s"

  type {
    category       = "Domain"
    indicator_type = "Domain"
    id             = "%[3]s"
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  status                = "Closed"
  confidence            = "90"
  first_occurrence_time = "2024-08-19T09:33:55.000+08:00"
  last_occurrence_time  = "2024-08-22T21:15:30.000+08:00"
  threat_degree         = "Gray"
  granularity           = "1"
  value                 = "test.terraform.com"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name, acceptance.HW_SECMASTER_INDICATOR_TYPE_ID)
}
