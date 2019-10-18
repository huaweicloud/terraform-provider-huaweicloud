package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccDcsInstancesV1_basic(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists("huaweicloud_dcs_instance_v1.instance_1", instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_dcs_instance_v1.instance_1", "name", instanceName),
					resource.TestCheckResourceAttr(
						"huaweicloud_dcs_instance_v1.instance_1", "engine", "Redis"),
				),
			},
		},
	})
}

func testAccCheckDcsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dcsClient, err := config.dcsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dcs_instance_v1" {
			continue
		}

		_, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("The Dcs instance still exists.")
		}
	}
	return nil
}

func testAccCheckDcsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dcsClient, err := config.dcsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud instance client: %s", err)
		}

		v, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting Huaweicloud instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("The Dcs instance not found.")
		}
		instance = *v
		return nil
	}
}

func testAccDcsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
       resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
         name = "secgroup_1"
         description = "secgroup_1"
       }
       data "huaweicloud_dcs_az_v1" "az_1" {
         port = "8002"
         code = "%s"
		}
       data "huaweicloud_dcs_product_v1" "product_1" {
          spec_code = "dcs.master_standby"
		}
		resource "huaweicloud_dcs_instance_v1" "instance_1" {
			name  = "%s"
          engine_version = "3.0.7"
          password = "Huawei_test"
          engine = "Redis"
          capacity = 2
          vpc_id = "%s"
          security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup_1.id}"
          subnet_id = "%s"
          available_zones = ["${data.huaweicloud_dcs_az_v1.az_1.id}"]
          product_id = "${data.huaweicloud_dcs_product_v1.product_1.id}"
          save_days = 1
          backup_type = "manual"
          begin_at = "00:00-01:00"
          period_type = "weekly"
          backup_at = [1]
          depends_on      = ["data.huaweicloud_dcs_product_v1.product_1", "huaweicloud_networking_secgroup_v2.secgroup_1"]
		}
	`, OS_AVAILABILITY_ZONE, instanceName, OS_VPC_ID, OS_NETWORK_ID)
}
