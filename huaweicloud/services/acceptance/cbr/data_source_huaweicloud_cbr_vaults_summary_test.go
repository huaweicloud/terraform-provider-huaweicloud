package cbr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVaultsSummary_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cbr_vaults_summary.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVaultsSummary_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "used_size"),
				),
			},
		},
	})
}

const testAccDataVaultsSummary_basic = `data "huaweicloud_cbr_vaults_summary" "test" {}`
