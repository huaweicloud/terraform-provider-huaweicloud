package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwServiceGroupMembers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_service_group_members.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwServiceGroupMembers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.item_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.source_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dest_port"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_source_port_filter_useful", "true"),
					resource.TestCheckOutput("is_dest_port_filter_useful", "true"),
					resource.TestCheckOutput("is_protocol_filter_useful", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceCfwServiceGroupMembers_predefinedGroupMembers(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_service_group_members.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwPredefinedServiceGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwServiceGroupMembers_predefinedGroupMembers(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.item_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.source_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dest_port"),
				),
			},
		},
	})
}

func testDataSourceCfwServiceGroupMembers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id          = huaweicloud_cfw_service_group_member.m2.id
  source_port = huaweicloud_cfw_service_group_member.m2.source_port
  dest_port   = huaweicloud_cfw_service_group_member.m1.dest_port
  protocol    = huaweicloud_cfw_service_group_member.m1.protocol
}

data "huaweicloud_cfw_service_group_members" "test" {
  group_id = huaweicloud_cfw_service_group.test.id

  depends_on = [
    huaweicloud_cfw_service_group_member.m1,
    huaweicloud_cfw_service_group_member.m2,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_cfw_service_group_members.test.records) >= 2
}

data "huaweicloud_cfw_service_group_members" "filter_by_id" {
  group_id = huaweicloud_cfw_service_group.test.id
  item_id  = local.id

  depends_on = [
    huaweicloud_cfw_service_group_member.m1,
    huaweicloud_cfw_service_group_member.m2,
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cfw_service_group_members.filter_by_id.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_service_group_members.filter_by_id.records[*] : v.item_id == local.id]
  )
}

data "huaweicloud_cfw_service_group_members" "filter_by_source_port" {
  group_id    = huaweicloud_cfw_service_group.test.id
  source_port = local.source_port
  
  depends_on = [
    huaweicloud_cfw_service_group_member.m1,
    huaweicloud_cfw_service_group_member.m2,
  ]
}

output "is_source_port_filter_useful" {
  value = length(data.huaweicloud_cfw_service_group_members.filter_by_source_port.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_service_group_members.filter_by_source_port.records[*] : v.source_port == local.source_port]
  )
}

data "huaweicloud_cfw_service_group_members" "filter_by_dest_port" {
  group_id  = huaweicloud_cfw_service_group.test.id
  dest_port = local.dest_port
  
  depends_on = [
    huaweicloud_cfw_service_group_member.m1,
    huaweicloud_cfw_service_group_member.m2,
  ]
}

output "is_dest_port_filter_useful" {
  value = length(data.huaweicloud_cfw_service_group_members.filter_by_dest_port.records) >= 1 && alltrue(
    [for v in data.huaweicloud_cfw_service_group_members.filter_by_dest_port.records[*] : v.dest_port == local.dest_port]
  )
}

data "huaweicloud_cfw_service_group_members" "filter_by_protocol" {
  group_id = huaweicloud_cfw_service_group.test.id
  protocol = local.protocol

  depends_on = [
    huaweicloud_cfw_service_group_member.m1,
    huaweicloud_cfw_service_group_member.m2,
  ]
}

output "is_protocol_filter_useful" {
  value = length(data.huaweicloud_cfw_service_group_members.filter_by_protocol.records) >= 1 && alltrue([
    for v in data.huaweicloud_cfw_service_group_members.filter_by_protocol.records[*] : 
      v.protocol == local.protocol
  ])
}
`, testDataSourceCfwServiceGroupMembers_base(name))
}

func testDataSourceCfwServiceGroupMembers_predefinedGroupMembers() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_service_group_members" "test" {
  group_id   = "%[1]s"
  group_type = "1"
}
`, acceptance.HW_CFW_PREDEFINED_SERVICE_GROUP1)
}

func testDataSourceCfwServiceGroupMembers_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_service_group_member" "m1" {
  group_id    = huaweicloud_cfw_service_group.test.id
  protocol    = 6
  source_port = "80"
  dest_port   = "8080"
}

resource "huaweicloud_cfw_service_group_member" "m2" {
  group_id    = huaweicloud_cfw_service_group.test.id
  protocol    = 17
  source_port = "80"
  dest_port   = "81"
}
`, testServiceGroup_basic(name))
}
