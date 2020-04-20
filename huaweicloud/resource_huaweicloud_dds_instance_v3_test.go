package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dds/v3/instances"
)

func TestAccDDSV3Instance_basic(t *testing.T) {
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDDSInstanceV3Config_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists("huaweicloud_dds_instance_v3.instance", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_dds_instance_v3.instance", "name", "dds-instance"),
				),
			},
		},
	})
}

func testAccCheckDDSV3InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.ddsV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DDS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dds_instance_v3" {
			continue
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		instances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}

		if instances.TotalCount > 0 {
			return fmt.Errorf("Instance still exists. ")
		}
	}

	return nil
}

func testAccCheckDDSV3InstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.ddsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud DDS client: %s ", err)
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		instances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}
		if instances.TotalCount == 0 {
			return fmt.Errorf("Instance not found. ")
		}

		return nil
	}
}

var TestAccDDSInstanceV3Config_basic = fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup_acc" {
  name = "secgroup_acc"
}

resource "huaweicloud_dds_instance_v3" "instance" {
  name = "dds-instance"
  datastore {
    type = "DDS-Community"
    version = "3.4"
    storage_engine = "wiredTiger"
  }
  region = "%s"
  availability_zone = "%s"
  vpc_id = "%s"
  subnet_id = "%s"
  security_group_id = huaweicloud_networking_secgroup_v2.secgroup_acc.id
  password = "Test@123"
  mode = "Sharding"
  flavor {
    type = "mongos"
    num = 2
    spec_code = "dds.mongodb.c3.medium.4.mongos"
  }
  flavor {
    type = "shard"
    num = 2
    storage = "ULTRAHIGH"
    size = 20
    spec_code = "dds.mongodb.c3.medium.4.shard"
  }
  flavor {
    type = "config"
    num = 1
    storage = "ULTRAHIGH"
    size = 20
    spec_code = "dds.mongodb.c3.large.2.config"
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days = "8"
  }
}`, OS_REGION_NAME, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID)
