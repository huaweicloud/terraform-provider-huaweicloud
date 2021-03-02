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
	resourceName := "huaweicloud_dds_instance.instance"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDDSInstanceV3Config_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", "dds-instance"),
					resource.TestCheckResourceAttr(resourceName, "ssl", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"availability_zone",
					"password",
					"flavor",
				},
			},
		},
	})
}

func testAccCheckDDSV3InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.ddsV3Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DDS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dds_instance" {
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
		client, err := config.ddsV3Client(HW_REGION_NAME)
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
			return fmt.Errorf("dds instance not found.")
		}

		return nil
	}
}

var TestAccDDSInstanceV3Config_basic = fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_acc" {
  name = "secgroup_acc"
}

resource "huaweicloud_dds_instance" "instance" {
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
  security_group_id = huaweicloud_networking_secgroup.secgroup_acc.id
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
  tags = {
	foo = "bar"
    owner = "terraform"
  }
}`, HW_REGION_NAME, HW_AVAILABILITY_ZONE, HW_VPC_ID, HW_NETWORK_ID)
