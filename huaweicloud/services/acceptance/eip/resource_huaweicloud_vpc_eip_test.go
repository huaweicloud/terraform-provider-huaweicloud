package eip

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getEipResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
	}
	return eips.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcEip_basic(t *testing.T) {
	var (
		eip eips.PublicIp

		resourceName = "huaweicloud_vpc_eip.test"
		randName     = acceptance.RandomAccResourceName()
		udpateName   = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.size", "5"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.share_type", "PER"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.charge_mode", "traffic"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
			{
				Config: testAccVpcEip_update(udpateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", udpateName),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.ip_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", udpateName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.size", "8"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.charge_mode", "bandwidth"),
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

func TestAccVpcEip_share(t *testing.T) {
	var (
		eip eips.PublicIp

		randName     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_vpc_eip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_share(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
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

func TestAccVpcEip_WithEpsId(t *testing.T) {
	var (
		eip eips.PublicIp

		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_vpc_eip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_epsId(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccVpcEip_epsId_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccVpcEip_prePaid(t *testing.T) {
	var (
		eip eips.PublicIp

		randName     = acceptance.RandomAccResourceNameWithDash()
		updateName   = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_vpc_eip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_prePaid(randName, 5, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.size", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
			{
				Config: testAccVpcEip_prePaid(updateName, 5, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", updateName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.size", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
			{
				Config: testAccVpcEip_prePaid(updateName, 6, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", updateName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.size", "6"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"charging_mode", "period", "period_unit", "auto_renew"},
			},
		},
	})
}

func TestAccVpcEip_ChangeToPeriod(t *testing.T) {
	var (
		eip          eips.PublicIp
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_vpc_eip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
				),
			},
			{
				Config: testAccVpcEip_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.charge_mode", "bandwidth"),
				),
			},
			{
				Config: testAccVpcEip_prePaid(randName, 8, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
			{
				Config:      testAccVpcEip_prePaidChangeToPostPaid(randName, 8, true),
				ExpectError: regexp.MustCompile(`error updating the charging mode of the EIP`),
			},
			{
				Config: testAccVpcEip_prePaid(randName, 8, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"charging_mode", "period", "period_unit", "auto_renew"},
			},
		},
	})
}

func TestAccVpcEip_deprecated(t *testing.T) {
	var (
		eip eips.PublicIp

		randName        = acceptance.RandomAccResourceName()
		resourceName    = "huaweicloud_vpc_eip.test"
		vipResourceName = "huaweicloud_networking_vip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_deprecated(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "BOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.ip_version", "4"),
					resource.TestCheckResourceAttrPair(resourceName, "private_ip", vipResourceName, "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "port_id", vipResourceName, "id"),
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

func testAccVpcEip_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  name = "%[1]s"

  publicip {
    type       = "5_bgp"
    ip_version = 4
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = 5
    charge_mode = "traffic"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccVpcEip_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  name = "%[1]s"

  publicip {
    type       = "5_bgp"
    ip_version = 6
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = 8
    charge_mode = "bandwidth"
  }

  tags = {
    foo  = "bar1"
    key1 = "value"
  }
}
`, rName)
}

func testAccVpcEip_epsId(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  enterprise_project_id = "0"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = 5
    charge_mode = "traffic"
  }
}
`, rName)
}

func testAccVpcEip_epsId_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  enterprise_project_id = "%[1]s"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[2]s"
    size        = 5
    charge_mode = "traffic"
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, rName)
}

func testAccVpcEip_share(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%s"
  size = 5
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type = "WHOLE"
    id         = huaweicloud_vpc_bandwidth.test.id
  }
}
`, rName)
}

func testAccVpcEip_prePaid(rName string, size int, isAutoRenew bool) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  name = "%[1]s"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = %[2]d
    charge_mode = "bandwidth"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[3]v"
}
`, rName, size, isAutoRenew)
}

func testAccVpcEip_deprecated(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_vip" "test" {
  name       = "%[1]s"
  network_id = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type    = "5_bgp"
    port_id = huaweicloud_networking_vip.test.id
  }

  bandwidth {
    name        = "%[1]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, rName)
}

func testAccVpcEip_prePaidChangeToPostPaid(rName string, size int, _ bool) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  name = "%[1]s"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = %[2]d
    charge_mode = "bandwidth"
  }

  charging_mode = "postPaid"
}
`, rName, size)
}
