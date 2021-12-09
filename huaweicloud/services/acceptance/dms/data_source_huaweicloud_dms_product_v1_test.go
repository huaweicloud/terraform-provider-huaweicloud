package dms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsProductV1DataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_product.product1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(dataSourceName, "partition_num", "300"),
					resource.TestCheckResourceAttr(dataSourceName, "storage", "600"),
				),
			},
		},
	})
}

func TestAccDmsProductV1DataSource_rabbitmqSingle(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_product.product1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSource_rabbitmqSingle,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(dataSourceName, "io_type", "high"),
				),
			},
		},
	})
}

func TestAccDmsProductV1DataSource_rabbitmqCluster(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_product.product1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSource_rabbitmqCluster,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(dataSourceName, "io_type", "high"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_spec_code", "dms.physical.storage.high"),
				),
			},
		},
	})
}

var testAccDmsProductV1DataSource_basic = fmt.Sprintf(`
data "huaweicloud_dms_product" "product1" {
  engine            = "kafka"
  version           = "1.1.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}
`)

var testAccDmsProductV1DataSource_rabbitmqSingle = fmt.Sprintf(`
data "huaweicloud_dms_product" "product1" {
  engine            = "rabbitmq"
  instance_type     = "single"
  storage_spec_code = "dms.physical.storage.high"
}
`)

var testAccDmsProductV1DataSource_rabbitmqCluster = fmt.Sprintf(`
data "huaweicloud_dms_product" "product1" {
  engine            = "rabbitmq"
  instance_type     = "cluster"
  storage_spec_code = "dms.physical.storage.high"
}
`)
