package cbh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCbhFlavors_basic(t *testing.T) {
	rName := "data.huaweicloud_cbh_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCbhFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "flavors.#", "2"),
					resource.TestCheckResourceAttr(rName, "flavors.0.resource_spec", "cbh.basic.10"),
					resource.TestCheckResourceAttr(rName, "flavors.0.region", "cn-north-4"),
					resource.TestCheckResourceAttr(rName, "flavors.0.period_unit", "year"),
					resource.TestCheckResourceAttr(rName, "flavors.0.period", "1"),
					resource.TestCheckResourceAttr(rName, "flavors.0.subscription_num", "1"),
					resource.TestCheckResourceAttr(rName, "flavors.0.flavor_id", "OFFI812504551075774477"),
					resource.TestCheckResourceAttr(rName, "flavors.1.resource_spec", "cbh.basic.50"),
					resource.TestCheckResourceAttr(rName, "flavors.1.region", "cn-north-4"),
					resource.TestCheckResourceAttr(rName, "flavors.1.period_unit", "month"),
					resource.TestCheckResourceAttr(rName, "flavors.1.period", "1"),
					resource.TestCheckResourceAttr(rName, "flavors.1.subscription_num", "1"),
					resource.TestCheckResourceAttr(rName, "flavors.1.flavor_id", "OFFI740586375358963717"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.resource_spec"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.region"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.period_unit"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.period"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.subscription_num"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.1.resource_spec"),
					resource.TestCheckResourceAttrSet(rName, "flavors.1.region"),
					resource.TestCheckResourceAttrSet(rName, "flavors.1.period_unit"),
					resource.TestCheckResourceAttrSet(rName, "flavors.1.period"),
					resource.TestCheckResourceAttrSet(rName, "flavors.1.subscription_num"),
					resource.TestCheckResourceAttrSet(rName, "flavors.1.flavor_id"),
				),
			},
		},
	})
}

func testAccDatasourceCbhFlavors_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cbh_flavors" "test" {
  project_id = "%s"

  flavors {
    resource_spec    = "cbh.basic.10"
    region           = "cn-north-4"
    period_unit      = "year"
    period           = 1
    subscription_num = 1
  }
  flavors {
    resource_spec    = "cbh.basic.50"
    region           = "cn-north-4"
    period_unit      = "month"
    period           = 1
    subscription_num = 1
  }
}
`, acceptance.HW_PROJECT_ID)
}
