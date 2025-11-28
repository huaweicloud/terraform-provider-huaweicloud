package esw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEswFlavors_basic(t *testing.T) {
	rName := "data.huaweicloud_esw_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceEswFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "flavors.#"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.connections"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.bandwidth"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.pps"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.available_zones.#"),
				),
			},
		},
	})
}

func testAccDatasourceEswFlavors_basic() string {
	return `
data "huaweicloud_esw_flavors" "test" {}
`
}
