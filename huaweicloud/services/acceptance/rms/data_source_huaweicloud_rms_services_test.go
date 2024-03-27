package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsServices_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_services.basic"
	dataSource2 := "data.huaweicloud_rms_services.filter_by_name"
	dataSource3 := "data.huaweicloud_rms_services.filter_by_track"
	serviceName := "ecs"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsServices_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_track_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsServices_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_rms_services" "basic" {}

data "huaweicloud_rms_services" "filter_by_name" {
  name = "%[1]s"
}

data "huaweicloud_rms_services" "filter_by_track" {
  track = "tracked"
}

locals {
  name_filter_result  = [for v in data.huaweicloud_rms_services.filter_by_name.services[*].name : v == "%[1]s"]
  track_filter_result = [
    for v in data.huaweicloud_rms_services.filter_by_track.services: alltrue([for vv in v.resource_types: vv.track == "tracked"])]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_services.basic.services) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "is_track_filter_useful" {
  value = alltrue(local.track_filter_result) && length(local.track_filter_result) > 0
}
`, name)
}
