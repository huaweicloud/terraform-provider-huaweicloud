package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/instances"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccDmsInstancesV1_basic(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	var instanceUpdate = fmt.Sprintf("dms_instance_update_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDmsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists("huaweicloud_dms_instance_v1.instance_1", instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_instance_v1.instance_1", "name", instanceName),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_instance_v1.instance_1", "engine", "rabbitmq"),
				),
			},
			resource.TestStep{
				Config: testAccDmsV1Instance_update(instanceUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists("huaweicloud_dms_instance_v1.instance_1", instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_instance_v1.instance_1", "name", instanceUpdate),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_instance_v1.instance_1", "description", "instance update description"),
				),
			},
		},
	})
}

func TestAccDmsInstancesV1_KafkaInstance(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDmsV1Instance_KafkaInstance(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists("huaweicloud_dms_instance_v1.instance_1", instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_instance_v1.instance_1", "name", instanceName),
				),
			},
		},
	})
}

func testAccCheckDmsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dmsClient, err := config.dmsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dms_instance_v1" {
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

		config := testAccProvider.Meta().(*Config)
		dmsClient, err := config.dmsV1Client(OS_REGION_NAME)
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
       resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
         name = "secgroup_1"
         description = "secgroup_1"
       }
       data "huaweicloud_dms_az_v1" "az_1" {
		}
       data "huaweicloud_dms_product_v1" "product_1" {
          engine = "rabbitmq"
          instance_type = "single"
          version = "3.7.0"
		}
		resource "huaweicloud_dms_instance_v1" "instance_1" {
			name  = "%s"
          engine = "rabbitmq"
          storage_space = "${data.huaweicloud_dms_product_v1.product_1.storage}"
          access_user = "user"
          password = "Dmstest@123"
          vpc_id = "%s"
          security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup_1.id}"
          subnet_id = "%s"
          available_zones = ["${data.huaweicloud_dms_az_v1.az_1.id}"]
          product_id = "${data.huaweicloud_dms_product_v1.product_1.id}"
          engine_version = "${data.huaweicloud_dms_product_v1.product_1.version}"
          depends_on      = ["data.huaweicloud_dms_product_v1.product_1", "huaweicloud_networking_secgroup_v2.secgroup_1"]
		}
	`, instanceName, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccDmsV1Instance_update(instanceUpdate string) string {
	return fmt.Sprintf(`
       resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
         name = "secgroup_1"
         description = "secgroup_1"
       }
       data "huaweicloud_dms_az_v1" "az_1" {
		}
       data "huaweicloud_dms_product_v1" "product_1" {
          engine = "rabbitmq"
          instance_type = "single"
          version = "3.7.0"
		}
		resource "huaweicloud_dms_instance_v1" "instance_1" {
			name  = "%s"
          description = "instance update description"
          engine = "rabbitmq"
          storage_space = "${data.huaweicloud_dms_product_v1.product_1.storage}"
          access_user = "user"
          password = "Dmstest@123"
          vpc_id = "%s"
          security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup_1.id}"
          subnet_id = "%s"
          available_zones = ["${data.huaweicloud_dms_az_v1.az_1.id}"]
          product_id = "${data.huaweicloud_dms_product_v1.product_1.id}"
          engine_version = "${data.huaweicloud_dms_product_v1.product_1.version}"
          depends_on      = ["data.huaweicloud_dms_product_v1.product_1", "huaweicloud_networking_secgroup_v2.secgroup_1"]
		}
	`, instanceUpdate, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccDmsV1Instance_KafkaInstance(instanceName string) string {
	return fmt.Sprintf(`
       resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
         name = "secgroup_1"
         description = "secgroup_1"
       }
       data "huaweicloud_dms_az_v1" "az_1" {
		}
       data "huaweicloud_dms_product_v1" "product_1" {
          engine = "kafka"
          instance_type = "cluster"
          version = "1.1.0"
		}
		resource "huaweicloud_dms_instance_v1" "instance_1" {
			name  = "%s"
          engine = "kafka"
          partition_num = "${data.huaweicloud_dms_product_v1.product_1.partition_num}"
          storage_space = "${data.huaweicloud_dms_product_v1.product_1.storage}"
          specification = "${data.huaweicloud_dms_product_v1.product_1.bandwidth}"
          vpc_id = "%s"
          security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup_1.id}"
          subnet_id = "%s"
          available_zones = ["${data.huaweicloud_dms_az_v1.az_1.id}"]
          product_id = "${data.huaweicloud_dms_product_v1.product_1.id}"
          engine_version = "${data.huaweicloud_dms_product_v1.product_1.version}"
          storage_spec_code = "${data.huaweicloud_dms_product_v1.product_1.storage_spec_code}"
          depends_on      = ["data.huaweicloud_dms_product_v1.product_1", "huaweicloud_networking_secgroup_v2.secgroup_1"]
		}
	`, instanceName, OS_VPC_ID, OS_NETWORK_ID)
}
