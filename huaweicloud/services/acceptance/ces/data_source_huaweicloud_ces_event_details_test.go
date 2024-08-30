package ces

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesEventDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_event_details.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)
	resourceName := "huaweicloud_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCesTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesEvents_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/24"),
				),
			},
			{
				Config: testDataSourceCesEvents_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/24"),
				),
			},
			{
				Config: testDataSourceCesEventDetails_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "event_info.0.event_id"),
					resource.TestCheckResourceAttrSet(dataSource, "event_info.0.event_name"),
					resource.TestCheckResourceAttrSet(dataSource, "event_info.0.event_source"),
					resource.TestCheckResourceAttrSet(dataSource, "event_info.0.detail.0.content"),
					resource.TestCheckResourceAttrSet(dataSource, "event_info.0.detail.0.event_state"),
					resource.TestCheckResourceAttrSet(dataSource, "event_info.0.detail.0.event_level"),
					resource.TestCheckResourceAttrSet(dataSource, "event_info.0.detail.0.event_type"),
					resource.TestMatchResourceAttr(dataSource,
						"event_info.0.time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_source_filter_useful", "true"),
					resource.TestCheckOutput("is_timeRange_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesEventDetails_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s-update"
  cidr = "192.168.0.0/16"
}
	
resource "huaweicloud_vpc_subnet" "test" {
  name              = "%[1]s-update"
  cidr              = "192.168.0.0/24"
  gateway_ip        = "192.168.0.1"
  vpc_id            = huaweicloud_vpc.test.id
  description       = "created by acc test"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
	  
  tags = {
    foo = "bar"
    key = "value"
  }
}
	
locals {
  name       = "modifySubnet"
  type       = "EVENT.SYS"
  source     = "SYS.VPC"
  start_time = "%[2]s"
  end_time   = "%[3]s"
}
	
data "huaweicloud_ces_event_details" "test" {
  name = local.name
  type = local.type 

  depends_on = [huaweicloud_vpc_subnet.test]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_ces_event_details.test.event_info) >= 1 && alltrue(
    [for item in data.huaweicloud_ces_event_details.test.event_info[*]: item.event_name == local.name]
  )
}

data "huaweicloud_ces_event_details" "filter_by_source" {
  name   = local.name
  type   = local.type
  source = local.source
  
  depends_on = [huaweicloud_vpc_subnet.test]
}

output "is_source_filter_useful" {
  value = length(data.huaweicloud_ces_event_details.filter_by_source.event_info) >= 1 && alltrue(
    [for item in data.huaweicloud_ces_event_details.filter_by_source.event_info[*]: item.event_source == local.source]
  )
}

data "huaweicloud_ces_event_details" "filter_by_timeRange" {
  name = local.name
  type = local.type
  from = local.start_time
  to   = local.end_time

  depends_on = [huaweicloud_vpc_subnet.test]
}

output "is_timeRange_filter_useful" {
  value = length(data.huaweicloud_ces_event_details.filter_by_timeRange.event_info) >= 1 && alltrue(
    [for item in data.huaweicloud_ces_event_details.filter_by_timeRange.event_info[*]: item.event_name == local.name]
  )
}
`, name, acceptance.HW_CES_START_TIME, acceptance.HW_CES_END_TIME)
}
