package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDmsKafkaInstances_basic(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "huaweicloud_dms_kafka_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"manager_password",
					"used_storage_space",
				},
			},
			{
				Config: testAccDmsKafkaInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "kafka test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
		},
	})
}

func TestAccDmsKafkaInstances_withEpsId(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_dms_kafka_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaInstance_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckDmsKafkaInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud dms instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dms_kafka_instance" {
			continue
		}

		_, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("The Dms kafka instance still exists.")
		}
	}
	return nil
}

func testAccCheckDmsKafkaInstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud dms instance client: %s", err)
		}

		v, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting HuaweiCloud dms kafka instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("The Dms kafka instance not found.")
		}
		instance = *v
		return nil
	}
}

func testAccDmsKafkaInstance_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_az" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dms_product" "test" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "2.3.0"
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
  description = "secgroup for kafka"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  description       = "kafka test"
  access_user       = "user"
  password          = "Kafkatest@123"
  vpc_id            = data.huaweicloud_vpc.test.id
  network_id        = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  available_zones   = [data.huaweicloud_dms_az.test.id]
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  bandwidth         = data.huaweicloud_dms_product.test.bandwidth
  storage_space     = data.huaweicloud_dms_product.test.storage
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code
  manager_user      = "kafka-user"
  manager_password  = "Kafkatest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, rName)
}

func testAccDmsKafkaInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_az" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dms_product" "test" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "2.3.0"
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
  description = "secgroup for kafka"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  description       = "kafka test update"
  access_user       = "user"
  password          = "Kafkatest@123"
  vpc_id            = data.huaweicloud_vpc.test.id
  network_id        = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  available_zones   = [data.huaweicloud_dms_az.test.id]
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  bandwidth         = data.huaweicloud_dms_product.test.bandwidth
  storage_space     = data.huaweicloud_dms_product.test.storage
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code
  manager_user      = "kafka-user"
  manager_password  = "Kafkatest@123"

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}
`, rName, updateName)
}

func testAccDmsKafkaInstance_withEpsId(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_az" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dms_product" "test" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "2.3.0"
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
  description = "secgroup for kafka"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name                  = "%s"
  description           = "kafka test"
  access_user           = "user"
  password              = "Kafkatest@123"
  vpc_id                = data.huaweicloud_vpc.test.id
  network_id            = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  available_zones       = [data.huaweicloud_dms_az.test.id]
  product_id            = data.huaweicloud_dms_product.test.id
  engine_version        = data.huaweicloud_dms_product.test.version
  bandwidth             = data.huaweicloud_dms_product.test.bandwidth
  storage_space         = data.huaweicloud_dms_product.test.storage
  storage_spec_code     = data.huaweicloud_dms_product.test.storage_spec_code
  manager_user          = "kafka-user"
  manager_password      = "Kafkatest@123"
  enterprise_project_id = "%s"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}
