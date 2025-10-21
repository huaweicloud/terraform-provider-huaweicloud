package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseFeatureGates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_feature_gates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseFeatureGates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "enable_enterprise"),
					resource.TestCheckResourceAttrSet(dataSource, "cer_available"),
					resource.TestCheckResourceAttrSet(dataSource, "enable_user_def_obs"),
				),
			},
		},
	})
}

const testDataSourceSwrEnterpriseFeatureGates_basic = `data "huaweicloud_swr_enterprise_feature_gates" "test" {}`
