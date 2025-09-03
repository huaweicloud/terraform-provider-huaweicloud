package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocApplicationCapacities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_application_capacities.test"
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
				Config: testDataSourceDataSourceCocApplicationCapacities_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sum_cpu"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sum_mem"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.cloud_service_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckOutput("group_id_filter_is_useful", "true"),
					resource.TestCheckOutput("component_id_filter_is_useful", "true"),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocApplicationCapacities_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_application_capacities" "test" {
  group_id = huaweicloud_coc_group.test.id
  provider_obj {
    cloud_service_name = "ecs"
    type               = "cloudservers"
  }
  depends_on = [huaweicloud_coc_group_resource_relation.test]
}

output "group_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_capacities.test.data) > 0
}

data "huaweicloud_coc_application_capacities" "component_id_filter" {
  component_id = huaweicloud_coc_component.test.id
  provider_obj {
    cloud_service_name = "ecs"
    type               = "cloudservers"
  }
  depends_on = [huaweicloud_coc_group_resource_relation.test]
}

output "component_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_capacities.component_id_filter.data) > 0
}

data "huaweicloud_coc_application_capacities" "application_id_filter" {
  application_id = huaweicloud_coc_application.test.id
  provider_obj {
    cloud_service_name = "ecs"
    type               = "cloudservers"
  }
  depends_on = [huaweicloud_coc_group_resource_relation.test]
}

output "application_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_capacities.application_id_filter.data) > 0
}
`, testAccGroupResourceRelation_basic(name))
}
