package deprecated

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dms/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDmsInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud DMS client(V1): %s", err)
	}

	return instances.Get(client, state.Primary.ID).Extract()
}

func TestAccDmsInstancesV1_Rabbitmq(t *testing.T) {
	var instance instances.Instance
	var instanceName = acceptance.RandomAccResourceName()
	var instanceUpdate = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dms_instance.instance_1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccDmsV1Instance_update(instanceName, instanceUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dms_instance.instance_1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Instance_KafkaInstance(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func testAccDmsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "vpc_subnet_1" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.vpc_1.id
}

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
  vpc_id            = huaweicloud_vpc.vpc_1.id
  subnet_id         = huaweicloud_vpc_subnet.vpc_subnet_1.id
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
}`, instanceName, instanceName, instanceName)
}

func testAccDmsV1Instance_update(instanceName, instanceUpdate string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "vpc_subnet_1" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.vpc_1.id
}

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
  vpc_id            = huaweicloud_vpc.vpc_1.id
  subnet_id         = huaweicloud_vpc_subnet.vpc_subnet_1.id
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
}`, instanceName, instanceName, instanceUpdate)
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
}`, instanceName, acceptance.HW_VPC_ID, acceptance.HW_NETWORK_ID)
}
