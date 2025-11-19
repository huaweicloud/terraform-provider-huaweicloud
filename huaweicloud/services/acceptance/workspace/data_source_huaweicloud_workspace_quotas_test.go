package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataQuotas_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_workspace_quotas.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "quotas.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "quotas.0.resources.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.resources.0.quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.resources.0.used"),
					resource.TestMatchResourceAttr(dataSourceName, "site_quotas.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "site_quotas.0.resources.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "site_quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "site_quotas.0.resources.0.quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "site_quotas.0.resources.0.used"),
				),
			},
		},
	})
}

const testAccDataQuotas_basic = `
data "huaweicloud_workspace_quotas" "test" {}
`
