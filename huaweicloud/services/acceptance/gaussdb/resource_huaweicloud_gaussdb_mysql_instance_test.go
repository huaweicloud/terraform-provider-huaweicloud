package gaussdb

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccGaussDBInstance_basic(t *testing.T) {
	var instance instances.TaurusDBInstance

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_gaussdb_mysql_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGaussDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBInstanceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "audit_log_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccGaussDBInstanceConfig_basicUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "audit_log_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccGaussDBInstance_prePaid(t *testing.T) {
	var (
		instance instances.TaurusDBInstance

		resourceName = "huaweicloud_gaussdb_mysql_instance.test"
		password     = acceptance.RandomPassword()
		rName        = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckGaussDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBInstanceConfig_prePaid(rName, password, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccGaussDBInstanceConfig_prePaid(rName, password, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func testAccCheckGaussDBInstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.GaussdbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud GaussDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_gaussdb_mysql_instance" {
			continue
		}

		v, err := instances.Get(client, rs.Primary.ID).Extract()
		if err == nil && v.Id == rs.Primary.ID {
			return fmtp.Errorf("Instance <%s> still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckGaussDBInstanceExists(n string, instance *instances.TaurusDBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s.", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set.")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.GaussdbV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("error creating HuaweiCloud GaussDB client: %s", err)
		}

		found, err := instances.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmtp.Errorf("instance <%s> not found.", rs.Primary.ID)
		}
		instance = found

		return nil
	}
}

func testAccGaussDBInstanceConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name                  = "%s"
  password              = "Test@12345678"
  flavor                = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = data.huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccVpcConfig_Base(rName), rName)
}

func testAccGaussDBInstanceConfig_basicUpdate(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name                  = "%s"
  password              = "Test@12345678"
  flavor                = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = data.huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"
  audit_log_enabled     = true

  tags = {
    foo_update = "bar"
    key        = "value_update"
  }
}
`, testAccVpcConfig_Base(rName), rName)
}

func testAccGaussDBInstanceConfig_prePaid(rName, password string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id

  flavor   = "gaussdb.mysql.4xlarge.x86.4"
  name     = "%[2]s"
  password = "%[3]s"

  enterprise_project_id = "0"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[4]v"
}
`, testAccVpcConfig_Base(rName), rName, password, isAutoRenew)
}
