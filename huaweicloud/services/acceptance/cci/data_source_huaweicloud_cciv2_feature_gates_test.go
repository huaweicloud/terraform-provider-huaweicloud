package cci

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2FeatureGates_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_feature_gates.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2FeatureGates_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "feature_gates.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "feature_gates.0.feature"),
					resource.TestCheckResourceAttrSet(dataSourceName, "feature_gates.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "feature_gates.0.value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "feature_gates.0.description"),
				),
			},
		},
	})
}

const testAccDataSourceV2FeatureGates_basic = `data "huaweicloud_cciv2_feature_gates" "test" {}`
