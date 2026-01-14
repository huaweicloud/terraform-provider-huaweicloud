package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAadDomainGlobalConfig_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aad_domain_global_config.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceAadDomainGlobalConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "support_tls.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cipher.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "cipher.0.algo"),
					resource.TestCheckResourceAttrSet(dataSource, "cipher.0.desc_cn"),
					resource.TestCheckResourceAttrSet(dataSource, "cipher.0.desc_en"),
				),
			},
		},
	})
}

const testDataSourceDataSourceAadDomainGlobalConfig_basic = `
data "huaweicloud_aad_domain_global_config" "test" {
}
`
