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
					resource.TestCheckResourceAttr(rName, "flavor_groups.#", "1"),
					resource.TestCheckResourceAttr(rName, "flavor_groups.0.groupType", "X86"),
					resource.TestCheckResourceAttr(rName, "flavor_groups.0.flavors.0.type_code", "hws.resource.type.ddm"),
					resource.TestCheckResourceAttr(rName, "flavor_groups.0.flavors.0.code", "ddm.c6.2xlarge.2"),
					resource.TestCheckResourceAttr(rName, "flavor_groups.0.flavors.0.iaas_code", "c6s.2xlarge.2"),
				),
			},
		},
	})
}

func testAccDatasourceDdmFlavors_basic() string {
	return `
data "huaweicloud_ddm_flavors" "test" {
  engine_id  = "8de8fef1-d188-34ea-b7e2-26bbf958476b"
  group_type = "X86"
  type_code  = "hws.resource.type.ddm"
}
`
}
