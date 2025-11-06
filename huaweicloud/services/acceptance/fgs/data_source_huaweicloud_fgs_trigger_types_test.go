package fgs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTriggerTypes_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_fgs_trigger_types.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTriggerTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "types.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

const testAccDataTriggerTypes_basic = `
data "huaweicloud_fgs_trigger_types" "all" {}
`
