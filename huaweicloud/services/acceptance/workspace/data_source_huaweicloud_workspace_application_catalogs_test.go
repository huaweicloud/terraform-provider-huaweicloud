package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplicationCatalogs_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_application_catalogs.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplicationCatalogs_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "catalogs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.id"),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.zh"),
					resource.TestCheckResourceAttrSet(all, "catalogs.0.en"),
				),
			},
		},
	})
}

const testAccDataApplicationCatalogs_basic = `
data "huaweicloud_workspace_application_catalogs" "all" {}
`
