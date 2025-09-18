package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTagAntileakageMap_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_waf_tag_antileakage_map.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTagAntileakageMap_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "leakagemap.#"),
					resource.TestCheckResourceAttrSet(dataSource, "leakagemap.0.sensitive.#"),
					resource.TestCheckResourceAttrSet(dataSource, "leakagemap.0.code.#"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.#"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.id_card"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.sensitive"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.phone"),
					resource.TestCheckResourceAttrSet(dataSource, "locale.0.email"),
				),
			},
		},
	})
}

const testAccDataSourceTagAntileakageMap_basic = `
data "huaweicloud_waf_tag_antileakage_map" "test" {
  lang = "en"
}`
