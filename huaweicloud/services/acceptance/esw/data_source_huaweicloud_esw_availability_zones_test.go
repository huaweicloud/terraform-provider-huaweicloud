package esw

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEswAvailabilityZones_basic(t *testing.T) {
	rName := "data.huaweicloud_esw_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceEswAvailabilityZones_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "available_zones.#"),
				),
			},
		},
	})
}

func testAccDatasourceEswAvailabilityZones_basic() string {
	return `
data "huaweicloud_esw_availability_zones" "test" {}
`
}
