package swr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrFeatureGates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_feature_gates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrFeatureGates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "enable_experience"),
					resource.TestCheckResourceAttrSet(dataSource, "enable_hss_service"),
					resource.TestCheckResourceAttrSet(dataSource, "enable_image_scan"),
					resource.TestCheckResourceAttrSet(dataSource, "enable_sm3"),
					resource.TestCheckResourceAttrSet(dataSource, "enable_image_sync"),
					resource.TestCheckResourceAttrSet(dataSource, "enable_cci_service"),
					resource.TestCheckResourceAttrSet(dataSource, "enable_image_label"),
					resource.TestCheckResourceAttrSet(dataSource, "enable_pipeline"),
				),
			},
		},
	})
}

const testDataSourceSwrFeatureGates_basic = `data "huaweicloud_swr_feature_gates" "test" {}`
