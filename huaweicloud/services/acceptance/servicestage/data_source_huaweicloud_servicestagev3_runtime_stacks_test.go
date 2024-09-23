package servicestage

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV3RuntimeStacks_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_servicestagev3_runtime_stacks.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV3RuntimeStacks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "runtime_stacks.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "runtime_stacks.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "runtime_stacks.0.deploy_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "runtime_stacks.0.version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "runtime_stacks.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "runtime_stacks.0.status"),
				),
			},
		},
	})
}

const testAccDataV3RuntimeStacks_basic = `data "huaweicloud_servicestagev3_runtime_stacks" "test" {}`
