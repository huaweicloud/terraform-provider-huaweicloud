package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDmsProductV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDms(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductV1DataSourceID("data.huaweicloud_dms_product_v1.product1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "engine", "kafka"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "partition_num", "300"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "storage", "600"),
				),
			},
		},
	})
}

func TestAccDmsProductV1DataSource_rabbitmqSingle(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDms(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSource_rabbitmqSingle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductV1DataSourceID("data.huaweicloud_dms_product_v1.product1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "node_num", "3"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "io_type", "normal"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "storage", "100"),
				),
			},
		},
	})
}

func TestAccDmsProductV1DataSource_rabbitmqCluster(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDms(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSource_rabbitmqCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductV1DataSourceID("data.huaweicloud_dms_product_v1.product1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "node_num", "5"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "storage", "500"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "io_type", "high"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_product_v1.product1", "storage_spec_code", "dms.physical.storage.high"),
				),
			},
		},
	})
}

func testAccCheckDmsProductV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find Dms product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Dms product data source ID not set")
		}

		return nil
	}
}

var testAccDmsProductV1DataSource_basic = fmt.Sprintf(`
data "huaweicloud_dms_product_v1" "product1" {
engine = "kafka"
version = "1.1.0"
instance_type = "cluster"
partition_num = 300
storage = 600
storage_spec_code = "dms.physical.storage.high"
}
`)

var testAccDmsProductV1DataSource_rabbitmqSingle = fmt.Sprintf(`
data "huaweicloud_dms_product_v1" "product1" {
engine = "rabbitmq"
version = "3.7.0"
instance_type = "single"
node_num = 3
storage = 100
storage_spec_code = "dms.physical.storage.normal"
}
`)

var testAccDmsProductV1DataSource_rabbitmqCluster = fmt.Sprintf(`
data "huaweicloud_dms_product_v1" "product1" {
engine = "rabbitmq"
version = "3.7.0"
instance_type = "cluster"
node_num = 5
storage = 500
storage_spec_code = "dms.physical.storage.high"
}
`)
