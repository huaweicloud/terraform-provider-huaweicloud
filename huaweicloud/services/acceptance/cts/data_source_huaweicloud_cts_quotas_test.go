package cts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCtsQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cts_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCtsQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.quota"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.used"),
				),
			},
		},
	})
}

const testDataSourceCtsQuotas_basic = `data "huaweicloud_cts_quotas" "test" {}`
