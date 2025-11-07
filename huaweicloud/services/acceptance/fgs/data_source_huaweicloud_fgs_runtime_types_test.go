package fgs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataRuntimeTypes_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_fgs_runtime_types.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataRuntimeTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "types.#", regexp.MustCompile(`[1-9][0-9]*`)),
					resource.TestCheckResourceAttrSet(all, "types.0.type"),
					resource.TestCheckResourceAttrSet(all, "types.0.display_name"),
				),
			},
		},
	})
}

const testAccDataRuntimeTypes_basic = `
data "huaweicloud_fgs_runtime_types" "test" {}
`
