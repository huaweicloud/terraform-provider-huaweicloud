package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEipVpcv3Bandwidths_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpcv3_bandwidths.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEipVpcv3Bandwidths_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.#"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.publicip_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.publicip_info.0.publicip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.publicip_info.0.publicip_id"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.publicip_info.0.publicip_type"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.publicip_info.0.ip_version"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.billing_info"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.admin_state"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.enable_bandwidth_rules"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.bandwidth_type"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.ingress_size"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.rule_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.ratio_95peak_plus"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.bandwidth_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidths.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("name_like_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_type_filter_is_useful", "true"),
					resource.TestCheckOutput("ingress_size_filter_is_useful", "true"),
					resource.TestCheckOutput("admin_state_filter_is_useful", "true"),
					resource.TestCheckOutput("billing_info_filter_is_useful", "true"),
					resource.TestCheckOutput("enable_bandwidth_rules_filter_is_useful", "true"),
					resource.TestCheckOutput("rule_quota_filter_is_useful", "true"),
					resource.TestCheckOutput("public_border_group_filter_is_useful", "true"),
					resource.TestCheckOutput("charge_mode_filter_is_useful", "true"),
					resource.TestCheckOutput("size_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("fields_filterr_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEipVpcv3Bandwidths_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%s"
  size = 5

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type = "WHOLE"
    id         = huaweicloud_vpc_bandwidth.test.id
  }
}
`, name)
}

func testDataSourceEipVpcv3Bandwidths_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpcv3_bandwidths" "test" {
  depends_on = [huaweicloud_vpc_eip.test]
}

locals {
  name = "%[2]s"
}

data "huaweicloud_vpcv3_bandwidths" "name_filter" {
  depends_on = [huaweicloud_vpc_eip.test]

  name = "%[2]s"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.name_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.name_filter.eip_bandwidths[*].name : v == local.name]
  )
}

locals {
  name_like = split("_", "%[2]s")[0]
}

data "huaweicloud_vpcv3_bandwidths" "name_like_filter" {
  depends_on = [huaweicloud_vpc_eip.test]

  name_like = split("_", "%[2]s")[0]
}

output "name_like_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.name_like_filter.eip_bandwidths) > 0
}

locals {
  bandwidth_type = "share"
}

data "huaweicloud_vpcv3_bandwidths" "bandwidth_type_filter" {
  depends_on = [huaweicloud_vpc_eip.test]

  bandwidth_type = "share"
}

output "bandwidth_type_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.bandwidth_type_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.bandwidth_type_filter.eip_bandwidths[*].bandwidth_type : v == local.bandwidth_type]
  )
}

locals {
  ingress_size = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].ingress_size
}

data "huaweicloud_vpcv3_bandwidths" "ingress_size_filter" {
  depends_on = [data.huaweicloud_vpcv3_bandwidths.test]

  ingress_size = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].ingress_size
}

output "ingress_size_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.ingress_size_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.ingress_size_filter.eip_bandwidths[*].ingress_size : v == local.ingress_size]
  )
}

locals {
  admin_state = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].admin_state
}

data "huaweicloud_vpcv3_bandwidths" "admin_state_filter" {
  depends_on = [data.huaweicloud_vpcv3_bandwidths.test]

  admin_state = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].admin_state
}

output "admin_state_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.admin_state_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.admin_state_filter.eip_bandwidths[*].admin_state : v == local.admin_state]
  )
}

locals {
  billing_info = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].billing_info
}

data "huaweicloud_vpcv3_bandwidths" "billing_info_filter" {
  depends_on = [data.huaweicloud_vpcv3_bandwidths.test]

  billing_info = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].billing_info
}

output "billing_info_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.billing_info_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.billing_info_filter.eip_bandwidths[*].billing_info : v == local.billing_info]
  )
}

locals {
  enable_bandwidth_rules = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].enable_bandwidth_rules
}

data "huaweicloud_vpcv3_bandwidths" "enable_bandwidth_rules_filter" {
  depends_on = [data.huaweicloud_vpcv3_bandwidths.test]

  enable_bandwidth_rules = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].enable_bandwidth_rules
}

output "enable_bandwidth_rules_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.enable_bandwidth_rules_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.enable_bandwidth_rules_filter.eip_bandwidths[*].enable_bandwidth_rules :
  v == local.enable_bandwidth_rules]
  )
}

locals {
  rule_quota = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].rule_quota
}

data "huaweicloud_vpcv3_bandwidths" "rule_quota_filter" {
  depends_on = [data.huaweicloud_vpcv3_bandwidths.test]

  rule_quota = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].rule_quota
}

output "rule_quota_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.rule_quota_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.rule_quota_filter.eip_bandwidths[*].rule_quota : v == local.rule_quota]
  )
}

locals {
  public_border_group = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].public_border_group
}

data "huaweicloud_vpcv3_bandwidths" "public_border_group_filter" {
  depends_on = [data.huaweicloud_vpcv3_bandwidths.test]

  public_border_group = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].public_border_group
}

output "public_border_group_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.public_border_group_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.public_border_group_filter.eip_bandwidths[*].public_border_group :
  v == local.public_border_group]
  )
}

data "huaweicloud_vpcv3_bandwidths" "charge_mode_filter" {
  depends_on = [huaweicloud_vpc_eip.test, data.huaweicloud_vpcv3_bandwidths.test]

  charge_mode = "bandwidth"
}

output "charge_mode_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.charge_mode_filter.eip_bandwidths) > 0
}

locals {
  size = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].size
}

data "huaweicloud_vpcv3_bandwidths" "size_filter" {
  depends_on = [data.huaweicloud_vpcv3_bandwidths.test]

  size = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].size
}

output "size_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.size_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.size_filter.eip_bandwidths[*].size : v == local.size]
  )
}

locals {
  type = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].type
}

data "huaweicloud_vpcv3_bandwidths" "type_filter" {
  depends_on = [data.huaweicloud_vpcv3_bandwidths.test]

  type = data.huaweicloud_vpcv3_bandwidths.test.eip_bandwidths[0].type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.type_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.type_filter.eip_bandwidths[*].type : v == local.type]
  )
}

locals {
  field = "name"
}

data "huaweicloud_vpcv3_bandwidths" "fields_filter" {
  depends_on = [huaweicloud_vpc_eip.test]

  fields = ["name"]
}

output "fields_filterr_is_useful" {
  value = length(data.huaweicloud_vpcv3_bandwidths.fields_filter.eip_bandwidths) > 0 && alltrue(
  [for v in data.huaweicloud_vpcv3_bandwidths.fields_filter.eip_bandwidths[*].name : v == local.name]
  )
}
`, testDataSourceEipVpcv3Bandwidths_base(name), name)
}
