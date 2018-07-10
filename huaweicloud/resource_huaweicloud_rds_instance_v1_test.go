package huaweicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
)

func TestAccRDSV1Instance_basic(t *testing.T) {
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRDSV1InstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccSInstanceV1Config_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists("huaweicloud_rds_instance_v1.instance", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_rds_instance_v1.instance", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckRDSV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	rdsClient, err := config.RdsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud rds: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_rds_instance_v1" {
			continue
		}

		_, err := instances.Get(rdsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Instance still exists. ")
		}
	}

	return nil
}

func testAccCheckRDSV1InstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		rdsClient, err := config.RdsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud rds client: %s ", err)
		}

		found, err := instances.Get(rdsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found. ")
		}

		*instance = *found
		time.Sleep(30 * time.Second)

		return nil
	}
}

var TestAccSInstanceV1Config_basic = fmt.Sprintf(`
data "huaweicloud_rds_flavors_v1" "flavor" {
    region = "cn-north-1"
    datastore_name = "PostgreSQL"
    datastore_version = "9.5.5"
    speccode = "rds.pg.s1.medium.ha"
}

resource "huaweicloud_compute_secgroup_v2" "secgrp_rds" {
  name        = "secgrp-rds-instance"
  description = "Rds Security Group"
}

resource "huaweicloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type = "PostgreSQL"
    version = "9.5.5"
  }
  flavorref = "${data.huaweicloud_rds_flavors_v1.flavor.id}"
  volume {
    type = "COMMON"
    size = 100
  }
  region = "cn-north-1"
  availabilityzone = "cn-north-1a"
  vpc = "%s"
  nics {
    subnetid = "%s"
  }
  securitygroup {
    id = "${huaweicloud_compute_secgroup_v2.secgrp_rds.id}"
  }
  dbport = "8635"
  backupstrategy = {
    starttime = "01:00:00"
    keepdays = 1
  }
  dbrtpd = "Huangwei!120521"
  ha = {
    enable = true
    replicationmode = "async"
  }
  depends_on = ["huaweicloud_compute_secgroup_v2.secgrp_rds"]
}`, OS_VPC_ID, OS_NETWORK_ID)
