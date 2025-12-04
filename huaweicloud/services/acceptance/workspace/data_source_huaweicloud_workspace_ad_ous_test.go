package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, please ensure the Workspace service is connected to the AD.
func TestAccDataAdOus_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_ad_ous.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAdOus_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					// Without any filter parameter.
					resource.TestMatchResourceAttr(all, "ous.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "ous.0.name"),
					resource.TestCheckResourceAttrSet(all, "ous.0.ou_dn"),
				),
			},
		},
	})
}

const testAccDataAdOus_basic = `
# Without any filter parameter.
data "huaweicloud_workspace_ad_ous" "all" {}
`
