package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceAggregatorResourceDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_aggregator_resource_detail.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourceAggregatorResourceDetail_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dataSource, "aggregator_id",
						"huaweicloud_rms_resource_aggregator.test", "id"),
					resource.TestCheckResourceAttrPair(dataSource, "resource_id",
						"data.huaweicloud_rms_resource_aggregator_discovered_resources.test", "resources.0.resource_id"),
					resource.TestCheckResourceAttrPair(dataSource, "service_type",
						"data.huaweicloud_rms_resource_aggregator_discovered_resources.test", "resources.0.service"),
					resource.TestCheckResourceAttrPair(dataSource, "type",
						"data.huaweicloud_rms_resource_aggregator_discovered_resources.test", "resources.0.type"),
					resource.TestCheckResourceAttrPair(dataSource, "region_id",
						"data.huaweicloud_rms_resource_aggregator_discovered_resources.test", "resources.0.region_id"),
					resource.TestCheckResourceAttrPair(dataSource, "resource_name",
						"data.huaweicloud_rms_resource_aggregator_discovered_resources.test", "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "aggregator_domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "ep_id"),
					resource.TestCheckResourceAttrSet(dataSource, "created"),
					resource.TestCheckResourceAttrSet(dataSource, "updated"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.%"),
					resource.TestCheckResourceAttrSet(dataSource, "properties.%"),
				),
			},
		},
	})
}

func testDataSourceResourceAggregatorResourceDetail_base(name, accountId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = "%[1]s"
  type        = "ACCOUNT"
  account_ids = ["%[2]s"]
}

data "huaweicloud_rms_resource_aggregator_discovered_resources" "test" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
}
`, name, accountId)
}

func testDataSourceResourceAggregatorResourceDetail_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rms_resource_aggregator_resource_detail" "test" {
  aggregator_id      = huaweicloud_rms_resource_aggregator.test.id
  resource_id        = data.huaweicloud_rms_resource_aggregator_discovered_resources.test.resources[0].resource_id
  service_type       = data.huaweicloud_rms_resource_aggregator_discovered_resources.test.resources[0].service
  type               = data.huaweicloud_rms_resource_aggregator_discovered_resources.test.resources[0].type
  source_account_id  = data.huaweicloud_rms_resource_aggregator_discovered_resources.test.resources[0].source_account_id
  region_id          = data.huaweicloud_rms_resource_aggregator_discovered_resources.test.resources[0].region_id
  resource_name      = data.huaweicloud_rms_resource_aggregator_discovered_resources.test.resources[0].resource_name
}
`, testDataSourceResourceAggregatorResourceDetail_base(name, acceptance.HW_DOMAIN_ID))
}
