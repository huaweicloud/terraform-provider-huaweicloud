package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcSubnetDataSource_ipv4Basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnet.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4Basic(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv4ByCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnet.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4ByCidr(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv4ByName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnet.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4ByName(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv4ByVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnet.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4ByVpcId(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv6Basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnet.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv6Basic(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv6_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "%s"
  gateway_ip = "%s"
}`, rName, cidr, rName, cidr, gatewayIp)
}

func testAccVpcSubnetDataSource_ipv4Basic(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnet" "test" {
  id = huaweicloud_vpc_subnet.test.id
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv4ByCidr(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnet" "test" {
  cidr = huaweicloud_vpc_subnet.test.cidr
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv4ByName(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnet" "test" {
  name = huaweicloud_vpc_subnet.test.name
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv4ByVpcId(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc_subnet.test.vpc_id
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv6Base(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%s"
  cidr        = "%s"
  gateway_ip  = "%s"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}`, rName, cidr, rName, cidr, gatewayIp)
}

func testAccVpcSubnetDataSource_ipv6Basic(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnet" "test" {
  id = huaweicloud_vpc_subnet.test.id
}
`, testAccVpcSubnetDataSource_ipv6Base(rName, cidr, gatewayIp))
}
