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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckGaussDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBInstanceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckGaussDBInstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.GaussdbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
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
			return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
		}

		found, err := instances.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmtp.Errorf("Instance <%s> not found.", rs.Primary.ID)
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
  name        = "%s"
  password    = "Test@123"
  flavor      = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  security_group_id = data.huaweicloud_networking_secgroup.test.id

  enterprise_project_id = "0"
}
`, testAccVpcConfig_Base(rName), rName)
}
