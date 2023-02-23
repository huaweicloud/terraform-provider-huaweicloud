package ddm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmFlavors_basic(t *testing.T) {
	rName := "data.huaweicloud_ddm_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "flavors.#", "1"),
					resource.TestCheckResourceAttr(rName, "flavors.0.cpu_arch", "X86"),
					resource.TestCheckResourceAttr(rName, "flavors.0.code", "ddm.c6.2xlarge.2"),
					resource.TestCheckResourceAttr(rName, "flavors.0.vcpus", "8"),
					resource.TestCheckResourceAttr(rName, "flavors.0.memory", "16"),
				),
			},
		},
	})
}

func testAccDatasourceDdmFlavors_basic() string {
	return `
data "huaweicloud_ddm_engines" test {
  version = "3.0.8.5"
}

data "huaweicloud_ddm_flavors" "test" {
  engine_id = data.huaweicloud_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
  code      = "ddm.c6.2xlarge.2"
  vcpus     = 8
  memory    = 16
}
`
}
