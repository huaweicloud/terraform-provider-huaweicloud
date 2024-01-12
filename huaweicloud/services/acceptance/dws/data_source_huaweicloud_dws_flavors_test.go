package dws

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDwsFlavorsDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_dws_flavors.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.datastore_version"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zones.#"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.vcpus", "8"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_memory(t *testing.T) {
	resourceName := "data.huaweicloud_dws_flavors.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_memory,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zones.#"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.memory", "64"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_all(t *testing.T) {
	resourceName := "data.huaweicloud_dws_flavors.test"
	dc := acceptance.InitDataSourceCheck(resourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_all,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zones.#"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.vcpus", "8"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.memory", "64"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_az(t *testing.T) {
	resourceName := "data.huaweicloud_dws_flavors.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_az,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
				),
			},
		},
	})
}

const testAccDwsFlavorsDataSource_basic = `
data "huaweicloud_dws_flavors" "test" {
  vcpus = 8
}
`

const testAccDwsFlavorsDataSource_memory = `
data "huaweicloud_dws_flavors" "test" {
  memory = 64
}
`

const testAccDwsFlavorsDataSource_all = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dws_flavors" "test" {
  vcpus             = 8
  memory            = 64
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`

const testAccDwsFlavorsDataSource_az = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dws_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`
