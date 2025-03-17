package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalEipQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_global_eip_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGlobalEipQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.min"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.quota"),
				),
			},
		},
	})
}

const testDataSourceGlobalEipQuotas_basic = `data "huaweicloud_global_eip_quotas" "test" {}`
