package cph

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceServerFlavors_basic(t *testing.T) {
	rName := "data.huaweicloud_cph_server_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceServerFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "flavors.0.vcpus", "64"),
					resource.TestCheckResourceAttr(rName, "flavors.0.type", "0"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.memory"),
					resource.TestCheckResourceAttr(rName, "flavors.0.extend_spec.#", "1"),
				),
			},
		},
	})
}

func testAccDatasourceServerFlavors_basic() string {
	return `
data "huaweicloud_cph_server_flavors" "test" {
  type  = "0"
  vcpus = 64
}
`
}
