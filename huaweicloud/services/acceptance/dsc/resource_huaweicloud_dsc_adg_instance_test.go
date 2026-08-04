package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
)

func getAdgInstanceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	return dsc.GetAdgInstanceById(client, state.Primary.ID)
}

func TestAccResourceDscAdgInstance_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_dsc_adg_instance.test"
		name         = acceptance.RandomAccResourceName()

		obj interface{}
		rc  = acceptance.InitResourceCheck(
			resourceName,
			&obj,
			getAdgInstanceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDscAdgInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "ADG"),
					resource.TestCheckResourceAttr(resourceName, "deploy_mode", "CLOUD"),
					resource.TestCheckResourceAttr(resourceName, "mode", "ha"),
					resource.TestCheckResourceAttr(resourceName, "admin_name", "sysadmin"),
					resource.TestCheckResourceAttr(resourceName, "charge_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "publicip_id",
						"huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip",
						"huaweicloud_vpc_eip.test1", "address"),
				),
			},
			{
				Config: testAccResourceDscAdgInstance_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "admin_name", "secadmin"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "publicip_id",
						"huaweicloud_vpc_eip.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip",
						"huaweicloud_vpc_eip.test2", "address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"admin_name",
					"outside_ins_config",
					"period_unit",
					"period",
					"auto_renew",
				},
			},
		},
	})
}

func testAccResourceDscAdgInstance_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s-adg"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s-adg"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s-adg"
}

resource "huaweicloud_vpc_eip" "test1" {
  name = "%[1]s-eip"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s-eip-bw1"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip" "test2" {
  name = "%[1]s-eip-up"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s-eip-bw2"
    size        = 5
    charge_mode = "traffic"
  }
}
`, name)
}

func testAccResourceDscAdgInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_adg_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  type              = "ADG"
  specification     = "dsc.adg.basic"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  admin_name        = "sysadmin"
  password          = "Test@12345678"
  deploy_mode       = "CLOUD"
  mode              = "ha"
  publicip_id       = huaweicloud_vpc_eip.test1.id

  charge_mode   = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`, testAccResourceDscAdgInstance_base(name), name)
}

func testAccResourceDscAdgInstance_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_adg_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  type              = "ADG"
  specification     = "dsc.adg.basic"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  admin_name        = "secadmin"
  password          = "Test@123456789"
  deploy_mode       = "CLOUD"
  mode              = "ha"
  publicip_id       = huaweicloud_vpc_eip.test2.id

  charge_mode   = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}
`, testAccResourceDscAdgInstance_base(name), name)
}
