package iec

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFlavorsDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_iec_flavors.flavors_test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFlavorsConfig(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "flavors.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttr(dataSourceName, "region", acceptance.HW_REGION_NAME),
				),
			},
		},
	})
}

func TestAccFlavorsDataSource_FilterName(t *testing.T) {
	dataSourceName := "data.huaweicloud_iec_flavors.flavors_test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFlavorsWithName(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "region", acceptance.HW_REGION_NAME),
				),
			},
		},
	})
}

func testAccFlavorsConfig() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_flavors" "flavors_test" {
  region = "%s"
}
	`, acceptance.HW_REGION_NAME)
}

func testAccFlavorsWithName() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_flavors" "flavors_test" {
  region = "%s"
  name   = "c6.large.2"
}
	`, acceptance.HW_REGION_NAME)
}
