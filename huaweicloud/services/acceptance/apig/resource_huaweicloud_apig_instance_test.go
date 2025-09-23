package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getInstanceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	return apig.QueryInstanceDetail(client, state.Primary.ID)
}

func TestAccInstance_basic(t *testing.T) {
	var (
		instance instances.Instance

		resourceName = "huaweicloud_apig_instance.test"
		rName        = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInstance_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "loadbalancer_provider", "elb"),
					resource.TestCheckResourceAttr(resourceName, "vpcep_service_name", "apig"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcep_service_address"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_charging_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "egress_address", ""),
					resource.TestCheckResourceAttr(resourceName, "ingress_address", ""),
					resource.TestCheckResourceAttr(resourceName, "custom_ingress_ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_ingress_ports.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "custom_ingress_ports.0.port", "3662"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_ingress_ports.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_ingress_ports.0.status"),
				),
			},
			{
				Config: testAccInstance_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "edition", "PROFESSIONAL"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_charging_mode", "bandwidth"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(resourceName, "tags.newKey", "value"),
					resource.TestCheckResourceAttr(resourceName, "vpcep_service_name", "new_custom_apig"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcep_service_address"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttrSet(resourceName, "egress_address"),
					resource.TestCheckResourceAttrSet(resourceName, "ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "custom_ingress_ports.#", "2"),
				),
			},
			{
				Config: testAccInstance_basic_step3(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "6"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_charging_mode", "traffic"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_size", "6"),
					resource.TestCheckResourceAttrSet(resourceName, "egress_address"),
					resource.TestCheckResourceAttrSet(resourceName, "ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "custom_ingress_ports.#", "0"),
				),
			},
			{
				Config: testAccInstance_basic_step4(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_charging_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "egress_address", ""),
					resource.TestCheckResourceAttr(resourceName, "ingress_address", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInstance_basic_step1(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_apig_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

  edition               = "BASIC"
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
  maintain_begin        = "14:00:00"
  description           = "created by acc test"

  tags = {
    foo = "bar"
    key = "value"
  }

  custom_ingress_ports {
    protocol = "HTTP"
    port     = 3662
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccInstance_basic_step2(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_networking_secgroup" "new" {
  name = "%[2]s_new"
}

resource "huaweicloud_apig_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.new.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  vpcep_service_name = "new_custom_apig"

  edition               = "PROFESSIONAL"
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
  maintain_begin        = "18:00:00"

  // Network configuration
  bandwidth_size                  = 5 # The bandwidth value must be greater than or equal to 5
  ingress_bandwidth_charging_mode = "bandwidth"
  ingress_bandwidth_size          = 5

  tags = {
    foo    = "baar"
    newKey = "value"
  }

  custom_ingress_ports {
    protocol = "HTTP"
    port     = 3662
  }

  custom_ingress_ports {
    protocol = "HTTPS"
    port     = 3665
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccInstance_basic_step3(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_networking_secgroup" "new" {
  name = "%[2]s_new"
}

resource "huaweicloud_apig_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.new.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  vpcep_service_name = "new_custom_apig"

  edition               = "PROFESSIONAL"
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
  maintain_begin        = "18:00:00"

  // Network configuration
  bandwidth_size                  = 6
  ingress_bandwidth_charging_mode = "traffic"
  ingress_bandwidth_size          = 6

  tags = {
    foo    = "baar"
    newKey = "value"
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccInstance_basic_step4(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_networking_secgroup" "new" {
  name = "%[2]s_new"
}

resource "huaweicloud_apig_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.new.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  vpcep_service_name = "new_custom_apig"

  edition               = "PROFESSIONAL"
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
  maintain_begin        = "18:00:00"

  tags = {
    foo    = "baar"
    newKey = "value"
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}
