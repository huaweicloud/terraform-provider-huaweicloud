package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDcsInstancesV1_basic(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "ip"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
				),
			},
			{
				Config: testAccDcsV1Instance_updated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_dcs_instance.instance_1", "backup_policy.0.begin_at", "01:00-02:00"),
					resource.TestCheckResourceAttr("huaweicloud_dcs_instance.instance_1", "backup_policy.0.save_days", "2"),
					resource.TestCheckResourceAttr("huaweicloud_dcs_instance.instance_1", "backup_policy.0.backup_at.#", "3"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV1_withEpsId(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_epsId(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccDcsInstancesV1_whitelists(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
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

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccDcsInstancesV1_single(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dcs_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_single(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
				),
			},
		},
	})
}

func testAccCheckDcsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dcsClient, err := config.DcsV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dcs_instance" {
			continue
		}

		_, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("the DCS instance still exists")
		}
	}
	return nil
}

func testAccCheckDcsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dcsClient, err := config.DcsV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud instance client: %s", err)
		}

		v, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmtp.Errorf("Error getting Huaweicloud instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmtp.Errorf("the DCS instance not found")
		}
		instance = *v
		return nil
	}
}

func testAccDcsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
	data "huaweicloud_availability_zones" "test" {}

	data "huaweicloud_vpc" "test" {
	  name = "vpc-default"
	}

	data "huaweicloud_vpc_subnet" "test" {
	  name = "subnet-default"
	}

	data "huaweicloud_dcs_az" "az_1" {
	  code = data.huaweicloud_availability_zones.test.names[0]
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 0.125
	  vpc_id            = data.huaweicloud_vpc.test.id
	  subnet_id         = data.huaweicloud_vpc_subnet.test.id
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "redis.ha.xu1.tiny.r2.128-h"
      backup_policy {
        backup_type = "manual"
        begin_at    = "00:00-01:00"
        period_type = "weekly"
        backup_at = [4]
        save_days = 1
      }

	  tags = {
	    key = "value"
	    owner = "terraform"
	  }
	}
	`, instanceName)
}

func testAccDcsV1Instance_updated(instanceName string) string {
	return fmt.Sprintf(`
	data "huaweicloud_availability_zones" "test" {}

	data "huaweicloud_vpc" "test" {
	  name = "vpc-default"
	}

	data "huaweicloud_vpc_subnet" "test" {
	  name = "subnet-default"
	}

	data "huaweicloud_dcs_az" "az_1" {
	  code = data.huaweicloud_availability_zones.test.names[0]
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 0.125
	  vpc_id            = data.huaweicloud_vpc.test.id
	  subnet_id         = data.huaweicloud_vpc_subnet.test.id
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "redis.ha.xu1.tiny.r2.128-h"
      backup_policy {
        backup_type = "manual"
        begin_at    = "01:00-02:00"
        period_type = "weekly"
        backup_at = [1, 2, 4]
        save_days = 2
      }

	  tags = {
	    key = "value"
	    owner = "terraform"
	  }
	}
	`, instanceName)
}

func testAccDcsV1Instance_epsId(instanceName string) string {
	return fmt.Sprintf(`
	data "huaweicloud_availability_zones" "test" {}

	data "huaweicloud_vpc" "test" {
	  name = "vpc-default"
	}

	data "huaweicloud_vpc_subnet" "test" {
	  name = "subnet-default"
	}

	data "huaweicloud_dcs_az" "az_1" {
		code = data.huaweicloud_availability_zones.test.names[0]
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 0.125
	  vpc_id            = data.huaweicloud_vpc.test.id
	  subnet_id         = data.huaweicloud_vpc_subnet.test.id
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "redis.ha.xu1.tiny.r2.128-h"
      backup_policy {
        backup_type = "manual"
        begin_at    = "00:00-01:00"
        period_type = "weekly"
        backup_at = [1]
        save_days = 1
      }
	  enterprise_project_id = "%s"
	}
	`, instanceName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDcsV1Instance_tiny(instanceName string) string {
	return fmt.Sprintf(`
	data "huaweicloud_availability_zones" "test" {}

	data "huaweicloud_vpc" "test" {
	  name = "vpc-default"
	}

	data "huaweicloud_vpc_subnet" "test" {
	  name = "subnet-default"
	}

	data "huaweicloud_dcs_az" "az_1" {
	  code = data.huaweicloud_availability_zones.test.names[0]
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 0.125
	  vpc_id            = data.huaweicloud_vpc.test.id
	  subnet_id         = data.huaweicloud_vpc_subnet.test.id
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "redis.ha.xu1.tiny.r2.128-h"
      backup_policy {
        backup_type = "manual"
        begin_at    = "00:00-01:00"
        period_type = "weekly"
        backup_at = [1]
        save_days = 1
      }
	}
	`, instanceName)
}

func testAccDcsV1Instance_single(instanceName string) string {
	return fmt.Sprintf(`
	data "huaweicloud_availability_zones" "test" {}

	data "huaweicloud_vpc" "test" {
	  name = "vpc-default"
	}

	data "huaweicloud_vpc_subnet" "test" {
	  name = "subnet-default"
	}

	data "huaweicloud_dcs_az" "az_1" {
	  code = data.huaweicloud_availability_zones.test.names[0]
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 2
	  vpc_id            = data.huaweicloud_vpc.test.id
	  subnet_id         = data.huaweicloud_vpc_subnet.test.id
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "redis.single.xu1.large.2-h"
	}
	`, instanceName)
}

func testAccDcsV1Instance_whitelists(instanceName string) string {
	return fmt.Sprintf(`
	data "huaweicloud_availability_zones" "test" {}

	data "huaweicloud_vpc" "test" {
	  name = "vpc-default"
	}

	data "huaweicloud_vpc_subnet" "test" {
	  name = "subnet-default"
	}

	data "huaweicloud_dcs_az" "az_1" {
	  code = data.huaweicloud_availability_zones.test.names[0]
	}

	resource "huaweicloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "Huawei_test"
	  engine            = "Redis"
	  capacity          = 2
	  vpc_id            = data.huaweicloud_vpc.test.id
	  subnet_id         = data.huaweicloud_vpc_subnet.test.id
	  available_zones   = [data.huaweicloud_dcs_az.az_1.id]
	  product_id        = "redis.ha.xu1.large.r2.2-h"
      backup_policy {
        backup_type = "manual"
        begin_at    = "00:00-01:00"
        period_type = "weekly"
        backup_at = [1]
        save_days = 1
      }

	  whitelists {
		group_name = "test-group1"
		ip_address = ["192.168.10.100", "192.168.0.0/24"]
	  }
	  whitelists {
		group_name = "test-group2"
		ip_address = ["172.16.10.100", "172.16.0.0/24"]
	  }
	}
	`, instanceName)
}
