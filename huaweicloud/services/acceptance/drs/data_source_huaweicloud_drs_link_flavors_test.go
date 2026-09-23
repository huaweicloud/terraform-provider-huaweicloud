package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsLinkFlavors_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_link_flavors.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsLinkFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "vm_flavor.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flow_flavor.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDrsLinkFlavors_filter(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_link_flavors.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsLinkFlavors_filter(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "vm_flavor.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flow_flavor.#"),
				),
			},
		},
	})
}

func testAccDataSourceDrsLinkFlavors_basic() string {
	return `
data "huaweicloud_drs_link_flavors" "test" {}
`
}

func testAccDataSourceDrsLinkFlavors_filter() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_link_flavors" "test" {
  engine_type   = "mysql"
  job_type      = "sync"
  job_direction = "up"
  net_type      = "vpc"
  node_type     = "high"
}
`)
}
