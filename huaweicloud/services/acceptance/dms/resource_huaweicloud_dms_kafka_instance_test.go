package dms

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDmsKafkaInstances_basic(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "huaweicloud_dms_kafka_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsKafkaInstanceDestroy,
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
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaInstance_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckDmsKafkaInstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dms instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dms_kafka_instance" {
			continue
		}

		_, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("The Dms kafka instance still exists.")
		}
	}
	return nil
}

func testAccCheckDmsKafkaInstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud dms instance client: %s", err)
		}

		v, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmtp.Errorf("Error getting HuaweiCloud dms kafka instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmtp.Errorf("The Dms kafka instance not found.")
		}
		instance = *v
		return nil
	}
}

func testAccDmsKafkaInstance_Base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_az" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/24"
  description = "test for kafka"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
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
`, rName, rName, rName)
}

func testAccDmsKafkaInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  description       = "kafka test"
  access_user       = "user"
  password          = "Kafkatest@123"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
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
`, testAccDmsKafkaInstance_Base(rName), rName)
}

func testAccDmsKafkaInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  description       = "kafka test update"
  access_user       = "user"
  password          = "Kafkatest@123"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
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
`, testAccDmsKafkaInstance_Base(rName), updateName)
}

func testAccDmsKafkaInstance_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_instance" "test" {
  name                  = "%s"
  description           = "kafka test"
  access_user           = "user"
  password              = "Kafkatest@123"
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id
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
`, testAccDmsKafkaInstance_Base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
