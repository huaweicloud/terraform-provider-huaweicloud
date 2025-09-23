package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocApplicationCapacityOrders_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_application_capacity_orders.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocApplicationCapacityOrders_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.rank_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.rank_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.rank_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.rank_list.0.value"),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
					resource.TestCheckOutput("component_id_filter_is_useful", "true"),
					resource.TestCheckOutput("group_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocApplicationCapacityOrders_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_application_capacity_orders" "test" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  application_id     = huaweicloud_coc_application.test.id
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "application_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_capacity_orders.test.data) > 0
}

data "huaweicloud_coc_application_capacity_orders" "component_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  component_id       = huaweicloud_coc_component.test.id
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "component_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_capacity_orders.component_id_filter.data) > 0
}

data "huaweicloud_coc_application_capacity_orders" "group_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  group_id           = huaweicloud_coc_group.test.id
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "group_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_capacity_orders.group_id_filter.data) > 0
}
`, testAccGroupResourceRelation_basic(name))
}
