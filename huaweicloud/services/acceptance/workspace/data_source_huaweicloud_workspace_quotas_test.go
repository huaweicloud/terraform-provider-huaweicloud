package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataQuotas_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_quotas.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "quotas.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestMatchResourceAttr(all, "quotas.0.resources.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.resources.0.quota"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.resources.0.used"),
					resource.TestMatchResourceAttr(all, "site_quotas.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestMatchResourceAttr(all, "site_quotas.0.resources.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "site_quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(all, "site_quotas.0.resources.0.quota"),
					resource.TestCheckResourceAttrSet(all, "site_quotas.0.resources.0.used"),
				),
			},
		},
	})
}

const testAccDataQuotas_basic = `
data "huaweicloud_workspace_quotas" "all" {}
`
