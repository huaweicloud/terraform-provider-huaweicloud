package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBNoSQLFlavors_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_default(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.engine", "cassandra"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
				),
			},
			{
				Config: testAccGaussDBNoSQLFlavors_cassandra(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "cassandra"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_mongodb(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_mongodb(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "mongodb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_influxdb(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_influxdb(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "influxdb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_redis(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_redis(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "redis"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_vcpus(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_vcpus(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "4"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_memory(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_memory(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "8"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_az(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_az(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dataSourceName, "flavors.0.availability_zones.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
				),
			},
		},
	})
}

func testAccGaussDBNoSQLFlavors_default() string {
	return (`
data "huaweicloud_gaussdb_nosql_flavors" "test" {}
`)
}

func testAccGaussDBNoSQLFlavors_cassandra() string {
	return (`
data "huaweicloud_gaussdb_nosql_flavors" "test" {
  engine = "cassandra"
}
`)
}

func testAccGaussDBNoSQLFlavors_mongodb() string {
	return (`
data "huaweicloud_gaussdb_nosql_flavors" "test" {
  engine = "mongodb"
}
`)
}

func testAccGaussDBNoSQLFlavors_influxdb() string {
	return (`
data "huaweicloud_gaussdb_nosql_flavors" "test" {
  engine = "influxdb"
}
`)
}

func testAccGaussDBNoSQLFlavors_redis() string {
	return (`
data "huaweicloud_gaussdb_nosql_flavors" "test" {
  engine = "redis"
}
`)
}

func testAccGaussDBNoSQLFlavors_vcpus() string {
	return (`
data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus = 4
}
`)
}

func testAccGaussDBNoSQLFlavors_memory() string {
	return (`
data "huaweicloud_gaussdb_nosql_flavors" "test" {
  memory = 8
}
`)
}

func testAccGaussDBNoSQLFlavors_az() string {
	return (`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`)
}
