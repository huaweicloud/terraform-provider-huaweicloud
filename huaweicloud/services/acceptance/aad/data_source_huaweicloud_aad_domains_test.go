package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAadDomains_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aad_domains.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare AAD protected domains before running this test cases.
			acceptance.TestAccPrecheckAadDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAadDomains_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.cname"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.protocol.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.real_server_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.real_servers"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.waf_status"),
				),
			},
		},
	})
}

const testDataSourceAadDomains_basic = `
data "huaweicloud_aad_domains" "test" {
}
`
