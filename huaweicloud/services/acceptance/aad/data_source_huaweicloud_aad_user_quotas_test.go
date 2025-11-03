package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceUserQuotas_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_user_quotas.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "instance"),
				),
			},
		},
	})
}

func TestAccDataSourceUserQuotas_waf(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_user_quotas.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserQuotas_waf(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "custom"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cc_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "geo_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "white_ip"),
				),
			},
		},
	})
}

func TestAccDataSourceUserQuotas_domainPort(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_user_quotas.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserQuotas_domainPort(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "domain_port_quota"),
				),
			},
		},
	})
}

func testAccDataSourceUserQuotas_basic() string {
	return `
data "huaweicloud_aad_user_quotas" "test" {
  type = "instance"
}
`
}

func testAccDataSourceUserQuotas_waf() string {
	return `
data "huaweicloud_aad_user_quotas" "test" {
  type = "waf"
}
`
}

func testAccDataSourceUserQuotas_domainPort() string {
	return `
data "huaweicloud_aad_user_quotas" "test" {
  type = "domain_port"
}
`
}
