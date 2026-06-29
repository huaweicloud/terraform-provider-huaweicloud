package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBProjectTags_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_project_tags.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBProjectTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.values.#"),
				),
			},
		},
	})
}

const testAccGeminiDBProjectTags_basic = `
data "huaweicloud_geminidb_project_tags" "test" {}
`
