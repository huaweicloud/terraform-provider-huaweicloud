package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplicationCatalogs_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_application_catalogs.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplicationCatalogs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "catalogs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "catalogs.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "catalogs.0.zh"),
					resource.TestCheckResourceAttrSet(dcName, "catalogs.0.en"),
				),
			},
		},
	})
}

const testAccDataApplicationCatalogs_basic = `
data "huaweicloud_workspace_application_catalogs" "test" {}
`
