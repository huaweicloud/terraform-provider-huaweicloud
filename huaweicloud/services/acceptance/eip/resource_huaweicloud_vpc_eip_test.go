package eip

import (
	"fmt"
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
		return nil, fmt.Errorf("error creating HuaweiCloud Network client: %s", err)
	}
	return eips.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcEIP_basic(t *testing.T) {
	var eip eips.PublicIp

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_eip.test"

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
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.share_type", "PER"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.charge_mode", "traffic"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
			{
				Config: testAccVpcEip_tags(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
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

func TestAccVpcEIP_share(t *testing.T) {
	var eip eips.PublicIp

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_eip.test"

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

func TestAccVpcEIP_WithEpsId(t *testing.T) {
	var eip eips.PublicIp

	randName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_vpc_eip.test"

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
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func TestAccVpcEIP_prePaid(t *testing.T) {
	var eip eips.PublicIp

	randName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_vpc_eip.test"

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
				Config: testAccVpcEip_prePaid(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
		},
	})
}

func TestAccVpcEIP_ipv6(t *testing.T) {
	var eip eips.PublicIp

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_eip.test"

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
				Config: testAccVpcEip_ipv6(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.ip_version", "6"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv6_address"),
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
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "%s"
    size        = 5
    charge_mode = "traffic"
  }
}
`, rName)
}

func testAccVpcEip_tags(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "%s"
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

func testAccVpcEip_epsId(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "%s"
    size        = 8
    charge_mode = "traffic"
  }
  enterprise_project_id = "%s"
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID)
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

func testAccVpcEip_prePaid(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1

  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "%s"
    size        = 5
  }
}
`, rName)
}

func testAccVpcEip_ipv6(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  name = "%s"

  publicip {
    type       = "5_bgp"
    ip_version = 6
  }
  bandwidth {
    share_type  = "PER"
    name        = "%s"
    size        = 5
    charge_mode = "traffic"
  }
}
`, rName, rName)
}
