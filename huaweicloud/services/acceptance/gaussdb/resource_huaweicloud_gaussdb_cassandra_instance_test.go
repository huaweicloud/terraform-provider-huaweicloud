package gaussdb

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccGeminiDBInstance_basic(t *testing.T) {
	var instance instances.GeminiDBInstance

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_gaussdb_cassandra_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGeminiDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstanceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
		},
	})
}

func TestAccGeminiDBInstance_prePaid(t *testing.T) {
	var instance instances.GeminiDBInstance

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_gaussdb_cassandra_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGeminiDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstanceConfig_prePaid(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccGeminiDBInstanceConfig_prePaid(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBInstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.GeminiDBV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GeminiDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_gaussdb_cassandra_instance" {
			continue
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
		if found.Id != "" {
			return fmtp.Errorf("Instance <%s> still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckGeminiDBInstanceExists(n string, instance *instances.GeminiDBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s.", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set.")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.GeminiDBV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud GeminiDB client: %s", err)
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("Instance <%s> not found.", rs.Primary.ID)
			}
			return err
		}
		if found.Id == "" {
			return fmtp.Errorf("Instance <%s> not found.", rs.Primary.ID)
		}
		instance = &found

		return nil
	}
}

func testAccVpcConfig_Base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"

  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = huaweicloud_vpc.test.id

  timeouts {
    delete = "20m"
  }
}
`, rName, rName)
}

func testAccGeminiDBInstanceConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_gaussdb_cassandra_instance" "test" {
  name        = "%s"
  password    = "Test@12345678"
  flavor      = "geminidb.cassandra.xlarge.4"
  volume_size = 100
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  ssl         = true
  node_num    = 4

  security_group_id = data.huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccVpcConfig_Base(rName), rName)
}

func testAccGeminiDBInstanceConfig_prePaid(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_gaussdb_cassandra_instance" "test" {
  name        = "%[2]s"
  password    = "Test@12345678"
  flavor      = "geminidb.cassandra.xlarge.4"
  volume_size = 100
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  ssl         = true
  node_num    = 4

  security_group_id = data.huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[3]v"
}
`, testAccVpcConfig_Base(rName), rName, isAutoRenew)
}
