package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCloudConnectionQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_cloud_connection_quotas.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcCloudConnectionQuotas_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_number"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_used"),
					resource.TestCheckOutput("cloud_connection_id_filter_is_useful", "true"),
					resource.TestCheckOutput("region_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcCloudConnectionQuotas_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name = "%s"
}
`, name)
}

func testDataSourceCcCloudConnectionQuotas_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cc_cloud_connection_quotas" "test" {
  quota_type = "cloud_connection" 
}

data "huaweicloud_cc_cloud_connection_quotas" "cloud_connection_id_filter" {
  quota_type          = "cloud_connection_region"
  cloud_connection_id = huaweicloud_cc_connection.test.id
}
locals {
  cloud_connection_id = huaweicloud_cc_connection.test.id
}
output "cloud_connection_id_filter_is_useful" {
  value = length(data.huaweicloud_cc_cloud_connection_quotas.cloud_connection_id_filter) > 0 && alltrue(
  [for v in data.huaweicloud_cc_cloud_connection_quotas.cloud_connection_id_filter[*].cloud_connection_id :
  v == local.cloud_connection_id]
  )
}

data "huaweicloud_cc_cloud_connection_quotas" "region_id_filter" {
  quota_type          = "region_network_instance"
  cloud_connection_id = huaweicloud_cc_connection.test.id
  region_id           = huaweicloud_cc_connection.test.region
}
locals {
  region_id = huaweicloud_cc_connection.test.region
}
output "region_id_filter_is_useful" {
  value = length(data.huaweicloud_cc_cloud_connection_quotas.region_id_filter) > 0 && alltrue(
  [for v in data.huaweicloud_cc_cloud_connection_quotas.region_id_filter[*].region_id : v == local.region_id]
  )
}
`, testDataSourceCcCloudConnectionQuotas_base(name))
}
