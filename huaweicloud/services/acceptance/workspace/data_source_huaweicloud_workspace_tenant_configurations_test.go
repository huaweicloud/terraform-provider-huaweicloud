package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTenantConfigurations_basic(t *testing.T) {
	all := "data.huaweicloud_workspace_tenant_configurations.all"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTenantConfigurations_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "configurations.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "configurations.0.id"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.name"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.status"),
				),
			},
		},
	})
}

const testAccDataTenantConfigurations_basic = `
data "huaweicloud_workspace_tenant_configurations" "all" {}
`
