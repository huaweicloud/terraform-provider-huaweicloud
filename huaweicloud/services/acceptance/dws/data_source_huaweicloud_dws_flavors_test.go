package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDwsFlavorsDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_dws_flavors.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDwsFlavorDataSourceID(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.vcpus", "8"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_memory(t *testing.T) {
	resourceName := "data.huaweicloud_dws_flavors.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_memory,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDwsFlavorDataSourceID(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.ram", "64"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_all(t *testing.T) {
	resourceName := "data.huaweicloud_dws_flavors.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDwsFlavorDataSourceID(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.vcpus", "8"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.ram", "64"),
					resource.TestCheckResourceAttrPair(resourceName, "flavors.0.availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_az(t *testing.T) {
	resourceName := "data.huaweicloud_dws_flavors.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_az,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDwsFlavorDataSourceID(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrPair(resourceName, "flavors.0.availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
				),
			},
		},
	})
}

func testAccCheckDwsFlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find dws flavors data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DWS flavors data source ID not set")
		}

		return nil
	}
}

const testAccDwsFlavorsDataSource_basic = `
data "huaweicloud_dws_flavors" "test" {
  vcpus = 8
}
`

const testAccDwsFlavorsDataSource_memory = `
data "huaweicloud_dws_flavors" "test" {
  ram = 64
}
`
const testAccDwsFlavorsDataSource_all = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dws_flavors" "test" {
  vcpus             = 8
  ram               = 64
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`

const testAccDwsFlavorsDataSource_az = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dws_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`
