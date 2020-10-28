package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
)

func TestAccDcsInstancesV1_basic(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "3.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "ip"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV1_withEpsId(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_epsId(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", OS_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func TestAccDcsInstancesV1_whitelists(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_whitelists(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.#", "2"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV1_tiny(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_tiny(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
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
		if rs.Type != "huaweicloud_dcs_instance" {
			continue
		}

		_, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("the DCS instance still exists")
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
			return fmt.Errorf("the DCS instance not found")
		}
		instance = *v
		return nil
	}
}

func testAccDcsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
	resource "huaweicloud_networking_secgroup" "secgroup_1" {
	  name        = "secgroup_1"
	  description = "secgroup_1"
	}
	data "huaweicloud_dcs_az" "az_1" {
	  code = "%s"
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "3.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 2
	  vpc_id            = "%s"
	  subnet_id         = "%s"
	  security_group_id = huaweicloud_networking_secgroup.secgroup_1.id
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "dcs.master_standby-h"
	  save_days         = 1
	  backup_type       = "manual"
	  begin_at          = "00:00-01:00"
	  period_type       = "weekly"
	  backup_at         = [1]

	  tags = {
	    key = "value"
	    owner = "terraform"
	  }
	}
	`, OS_AVAILABILITY_ZONE, instanceName, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccDcsV1Instance_epsId(instanceName string) string {
	return fmt.Sprintf(`
	resource "huaweicloud_networking_secgroup" "secgroup_1" {
	  name        = "secgroup_1"
	  description = "secgroup_1"
	}
	data "huaweicloud_dcs_az" "az_1" {
	  code = "%s"
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "3.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 2
	  vpc_id            = "%s"
	  security_group_id = huaweicloud_networking_secgroup.secgroup_1.id
	  subnet_id         = "%s"
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "dcs.master_standby-h"
	  save_days         = 1
	  backup_type       = "manual"
	  begin_at          = "00:00-01:00"
	  period_type       = "weekly"
	  backup_at         = [1]
	  enterprise_project_id = "%s"
	}
	`, OS_AVAILABILITY_ZONE, instanceName, OS_VPC_ID, OS_NETWORK_ID, OS_ENTERPRISE_PROJECT_ID)
}

func testAccDcsV1Instance_tiny(instanceName string) string {
	return fmt.Sprintf(`
	data "huaweicloud_dcs_az" "az_1" {
	  code = "%s"
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 0.125
	  vpc_id            = "%s"
	  subnet_id         = "%s"
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "redis.ha.au1.tiny.128-h"
	  save_days         = 1
	  backup_type       = "manual"
	  begin_at          = "00:00-01:00"
	  period_type       = "weekly"
	  backup_at         = [1]
	}
	`, OS_AVAILABILITY_ZONE, instanceName, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccDcsV1Instance_whitelists(instanceName string) string {
	return fmt.Sprintf(`
	data "huaweicloud_dcs_az" "az_1" {
	  code = "%s"
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 2
	  vpc_id            = "%s"
	  subnet_id         = "%s"
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "redis.ha.au1.large.r2.2-h"
	  save_days         = 1
	  backup_type       = "manual"
	  begin_at          = "00:00-01:00"
	  period_type       = "weekly"
	  backup_at         = [1]

	  whitelists {
		group_name = "test-group1"
		ip_address = ["192.168.10.100", "192.168.0.0/24"]
	  }
	  whitelists {
		group_name = "test-group2"
		ip_address = ["172.16.10.100", "172.16.0.0/24"]
	  }
	}
	`, OS_AVAILABILITY_ZONE, instanceName, OS_VPC_ID, OS_NETWORK_ID)
}
