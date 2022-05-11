package servicestage

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataComponentRuntimes_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_servicestage_component_runtimes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataComponentRuntimes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "runtimes.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func TestAccDataComponentRuntimes_byName(t *testing.T) {
	dataSourceName := "data.huaweicloud_servicestage_component_runtimes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataComponentRuntimes_byName,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "runtimes.#", "1"),
				),
			},
		},
	})
}

func TestAccDataComponentRuntimes_byPort(t *testing.T) {
	dataSourceName := "data.huaweicloud_servicestage_component_runtimes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataComponentRuntimes_byPort,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "runtimes.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccDataComponentRuntimes_basic = `
data "huaweicloud_servicestage_component_runtimes" "test" {}
`

const testAccDataComponentRuntimes_byName = `
data "huaweicloud_servicestage_component_runtimes" "test" {
  name = "Nodejs14"
}
`

const testAccDataComponentRuntimes_byPort = `
data "huaweicloud_servicestage_component_runtimes" "test" {
  default_port = 80
}
`
