package ces

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesEvents_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_events.test"
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
				Config: testDataSourceCesEvents_event(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.event_name"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.event_type"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.latest_event_source"),
					resource.TestMatchResourceAttr(dataSource,
						"events.0.latest_occur_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_timeRange_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesEvents_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name              = "%[1]s"
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
`, name)
}

func testDataSourceCesEvents_update(name string) string {
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
`, name)
}

func testDataSourceCesEvents_event(name string) string {
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

data "huaweicloud_ces_events" "test" {
  depends_on = [huaweicloud_vpc_subnet.test]
}

locals {
  name       = "modifySubnet"
  type       = "EVENT.SYS"
  start_time = "%[2]s"
  end_time   = "%[3]s"
}

data "huaweicloud_ces_events" "filter_by_name" {
  name = local.name

  depends_on = [huaweicloud_vpc_subnet.test]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_ces_events.filter_by_name.events) >= 1 && alltrue(
    [for e in data.huaweicloud_ces_events.filter_by_name.events : e.event_name == local.name]
  )
}

data "huaweicloud_ces_events" "filter_by_type" {
  type = local.type
  
  depends_on = [huaweicloud_vpc_subnet.test]
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_ces_events.filter_by_type.events) >= 1 && alltrue(
    [for e in data.huaweicloud_ces_events.filter_by_type.events : e.event_type == local.type]
  )
}

data "huaweicloud_ces_events" "filter_by_timeRange" {
  from = local.start_time
  to   = local.end_time

  depends_on = [huaweicloud_vpc_subnet.test]
}

output "is_timeRange_filter_useful" {
  value = length(data.huaweicloud_ces_events.filter_by_timeRange.events) >= 1 
}
`, name, acceptance.HW_CES_START_TIME, acceptance.HW_CES_END_TIME)
}
