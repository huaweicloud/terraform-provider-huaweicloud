package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/firewalls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getNetworkACLResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	iecClient, err := conf.IECV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IEC client: %s", err)
	}

	fwGroup, err := firewalls.Get(iecClient, state.Primary.ID).Extract()
	if err != nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return fwGroup, err
}

func TestAccNetworkACLResource_basic(t *testing.T) {
	var fwGroup firewalls.Firewall

	name := fmt.Sprintf("iec-acl-%s", acctest.RandString(5))
	rName := "huaweicloud_iec_network_acl.acl_demo"

	rc := acceptance.InitResourceCheck(
		rName,
		&fwGroup,
		getNetworkACLResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform test acc"),
					testAccCheckNetworkACLNetBlockExists(&fwGroup),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkACL_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform test acc"),
					testAccCheckNetworkACLNetBlockExists(&fwGroup),
				),
			},
		},
	})
}

func TestAccNetworkACLResource_no_subnets(t *testing.T) {
	var fwGroup firewalls.Firewall

	name := fmt.Sprintf("acc-fw-%s", acctest.RandString(5))
	rName := "huaweicloud_iec_network_acl.acl_demo"

	rc := acceptance.InitResourceCheck(
		rName,
		&fwGroup,
		getNetworkACLResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_no_subnets(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-noSubnet"),
					resource.TestCheckResourceAttr(rName, "description", " network acl without subents"),
					resource.TestCheckResourceAttr(rName, "status", "INACTIVE"),
				),
			},
		},
	})
}

func testAccCheckNetworkACLNetBlockExists(r *firewalls.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(r.Subnets) == 0 {
			return fmt.Errorf("the Subnet of IEC network ACL is not set")
		}
		return nil
	}
}

var testAccNetworkACLRules = `
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "vpc_demo"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name        = "subnet_demo"
  cidr        = "192.168.128.0/18"
  vpc_id      = huaweicloud_iec_vpc.vpc_test.id
  site_id     = data.huaweicloud_iec_sites.sites_test.sites[0].id
  gateway_ip  = "192.168.128.3"
}
`

func testAccNetworkACL_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_iec_network_acl" "acl_demo" {
  name        = "%s"
  description = "Created by terraform test acc"
  networks {
    vpc_id    = huaweicloud_iec_vpc.vpc_test.id
    subnet_id = huaweicloud_iec_vpc_subnet.subnet_test.id
  }
}
`, testAccNetworkACLRules, rName)
}

func testAccNetworkACL_basic_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_iec_network_acl" "acl_demo" {
  name        = "%s-update"
  description = "Updated by terraform test acc"
  networks {
    vpc_id    = huaweicloud_iec_vpc.vpc_test.id
    subnet_id = huaweicloud_iec_vpc_subnet.subnet_test.id
  }
}
`, testAccNetworkACLRules, rName)
}

func testAccNetworkACL_no_subnets(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_network_acl" "acl_demo" {
  name        = "%s-noSubnet"
  description = " network acl without subents"
}
`, rName)
}
