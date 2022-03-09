package dms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsProductDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_product.product1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(dataSourceName, "partition_num", "300"),
					resource.TestCheckResourceAttr(dataSourceName, "storage", "600"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_spec_codes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_spec_code", "dms.physical.storage.high"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.0"),
				),
			},
		},
	})
}

func TestAccDmsProductDataSource_kafkaVmSpec(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_product.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductDataSource_kafkaVmSpec,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(dataSourceName, "vm_specification", "c6.large.2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.0"),
				),
			},
		},
	})
}

func TestAccDmsProductDataSource_rabbitmqSingle(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_product.product1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductDataSource_rabbitmqSingle,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(dataSourceName, "io_type", "high"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_spec_codes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_spec_code", "dms.physical.storage.high"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.0"),
				),
			},
		},
	})
}

func TestAccDmsProductDataSource_rabbitmqCluster(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_product.product1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductDataSource_rabbitmqCluster,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(dataSourceName, "io_type", "high"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_spec_code", "dms.physical.storage.high"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_spec_codes.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "node_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.0"),
				),
			},
		},
	})
}

var testAccDmsProductDataSource_basic = fmt.Sprintf(`
data "huaweicloud_dms_product" "product1" {
  engine            = "kafka"
  version           = "1.1.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}
`)

var testAccDmsProductDataSource_kafkaVmSpec = fmt.Sprintf(`
data "huaweicloud_dms_product" "test" {
  instance_type    = "cluster"
  version          = "2.3.0"
  engine           = "kafka"
  vm_specification = "c6.large.2"
}
`)

var testAccDmsProductDataSource_rabbitmqSingle = fmt.Sprintf(`
data "huaweicloud_dms_product" "product1" {
  engine            = "rabbitmq"
  instance_type     = "single"
  storage_spec_code = "dms.physical.storage.high"
}
`)

var testAccDmsProductDataSource_rabbitmqCluster = fmt.Sprintf(`
data "huaweicloud_dms_product" "product1" {
  engine            = "rabbitmq"
  instance_type     = "cluster"
  storage_spec_code = "dms.physical.storage.high"
}
`)
