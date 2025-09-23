package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKMSQuotasDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_kms_quotas.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testKMSUserQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.resources.0.quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.resources.0.used"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.resources.0.type"),
				),
			},
		},
	})
}

const testKMSUserQuotas_basic = `data "huaweicloud_kms_quotas" "test" {}`
