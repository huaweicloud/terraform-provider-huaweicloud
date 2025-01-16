package swr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_key"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.used"),
				),
			},
		},
	})
}

const testDataSourceSwrQuotas_basic = `data "huaweicloud_swr_quotas" "test" {}`
