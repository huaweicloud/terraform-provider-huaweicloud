package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcGlobalConnectionBandwidthSpecCodes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_global_connection_bandwidth_spec_codes.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcGlobalConnectionBandwidthSpecCodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.local_area"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.remote_area"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.sku"),
				),
			},
			{
				Config: testDataSourceCcGlobalConnectionBandwidthSpecCodes_localArea(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "spec_codes.0.local_area", "cn-north-beijing4"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.remote_area"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.sku"),
				),
			},
			{
				Config: testDataSourceCcGlobalConnectionBandwidthSpecCodes_level(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "spec_codes.0.level", "Au"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.local_area"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.remote_area"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "spec_codes.0.sku"),
				),
			},
		},
	})
}

func testDataSourceCcGlobalConnectionBandwidthSpecCodes_basic() string {
	return `data "huaweicloud_cc_global_connection_bandwidth_spec_codes" "test" {}`
}

func testDataSourceCcGlobalConnectionBandwidthSpecCodes_localArea() string {
	return `
data "huaweicloud_cc_global_connection_bandwidth_spec_codes" "test" {
  local_area = "cn-north-beijing4"	
}`
}

func testDataSourceCcGlobalConnectionBandwidthSpecCodes_level() string {
	return `
data "huaweicloud_cc_global_connection_bandwidth_spec_codes" "test" {
  level = "Au"	
}`
}
