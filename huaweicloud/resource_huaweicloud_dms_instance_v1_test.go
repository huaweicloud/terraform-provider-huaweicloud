package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDmsInstancesV1_Rabbitmq(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	var instanceUpdate = fmt.Sprintf("dms_instance_update_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dms_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDms(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccDmsV1Instance_update(instanceUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "instance update description"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
		},
	})
}

func TestAccDmsInstancesV1_Kafka(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dms_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDms(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Instance_KafkaInstance(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func testAccCheckDmsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dms_instance" {
			continue
		}

		_, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("The Dms instance still exists.")
		}
	}
	return nil
}

func testAccCheckDmsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud instance client: %s", err)
		}

		v, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting HuaweiCloud instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("The Dms instance not found.")
		}
		instance = *v
		return nil
	}
}

func testAccDmsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_az" "az_1" {
}
data "huaweicloud_dms_product" "product_1" {
  engine        = "rabbitmq"
  instance_type = "single"
  version       = "3.7.17"
}

resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}
resource "huaweicloud_dms_instance" "instance_1" {
  name              = "%s"
  engine            = "rabbitmq"
  access_user       = "user"
  password          = "Dmstest@123"
  vpc_id            = "%s"
  subnet_id         = "%s"
  security_group_id = huaweicloud_networking_secgroup.secgroup_1.id
  available_zones   = [data.huaweicloud_dms_az.az_1.id]
  product_id        = data.huaweicloud_dms_product.product_1.id
  engine_version    = data.huaweicloud_dms_product.product_1.version
  storage_space     = data.huaweicloud_dms_product.product_1.storage
  storage_spec_code = data.huaweicloud_dms_product.product_1.storage_spec_code

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
	`, instanceName, HW_VPC_ID, HW_NETWORK_ID)
}

func testAccDmsV1Instance_update(instanceUpdate string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_az" "az_1" {
}
data "huaweicloud_dms_product" "product_1" {
  engine        = "rabbitmq"
  instance_type = "single"
  version       = "3.7.17"
}

resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}
resource "huaweicloud_dms_instance" "instance_1" {
  name              = "%s"
  description       = "instance update description"
  engine            = "rabbitmq"
  access_user       = "user"
  password          = "Dmstest@123"
  vpc_id            = "%s"
  subnet_id         = "%s"
  security_group_id = huaweicloud_networking_secgroup.secgroup_1.id
  available_zones   = [data.huaweicloud_dms_az.az_1.id]
  product_id        = data.huaweicloud_dms_product.product_1.id
  engine_version    = data.huaweicloud_dms_product.product_1.version
  storage_space     = data.huaweicloud_dms_product.product_1.storage
  storage_spec_code = data.huaweicloud_dms_product.product_1.storage_spec_code

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}
	`, instanceUpdate, HW_VPC_ID, HW_NETWORK_ID)
}

func testAccDmsV1Instance_KafkaInstance(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_az" "az_1" {
}
data "huaweicloud_dms_product" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "1.1.0"
}

resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}
resource "huaweicloud_dms_instance" "instance_1" {
  name              = "%s"
  engine            = "kafka"
  vpc_id            = "%s"
  subnet_id         = "%s"
  security_group_id = huaweicloud_networking_secgroup.secgroup_1.id
  available_zones   = [data.huaweicloud_dms_az.az_1.id]
  product_id        = data.huaweicloud_dms_product.product_1.id
  engine_version    = data.huaweicloud_dms_product.product_1.version
  specification     = data.huaweicloud_dms_product.product_1.bandwidth
  partition_num     = data.huaweicloud_dms_product.product_1.partition_num
  storage_space     = data.huaweicloud_dms_product.product_1.storage
  storage_spec_code = data.huaweicloud_dms_product.product_1.storage_spec_code

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
	`, instanceName, HW_VPC_ID, HW_NETWORK_ID)
}
