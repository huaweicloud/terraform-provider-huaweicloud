package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceQuotas_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_geminidb_quotas.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.used"),

					resource.TestCheckOutput("quotas_exist", "true"),
					resource.TestCheckOutput("datastore_type_and_mode_filter_useful", "true"),
					resource.TestCheckOutput("product_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceQuotas_basic = `
data "huaweicloud_geminidb_quotas" "test" {}

data "huaweicloud_geminidb_quotas" "datastore_type_and_mode_filter" {
  datastore_type = "cassandra"
  mode           = "Cluster"
}

data "huaweicloud_geminidb_quotas" "product_type_filter" {
  datastore_type = "redis"
  mode           = "CloudNativeCluster"
  product_type   = "Capacity"
}

output "quotas_exist" {
  value = length(data.huaweicloud_geminidb_quotas.test.quotas) > 0
}

output "datastore_type_and_mode_filter_useful" {
  value = length(data.huaweicloud_geminidb_quotas.datastore_type_and_mode_filter.quotas) > 0
}

output "product_type_filter_useful" {
  value = length(data.huaweicloud_geminidb_quotas.product_type_filter.quotas) > 0
}
`
