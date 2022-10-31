package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcSubnetsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetsDataSource_Basic(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "subnets.0.vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetsDataSource_Basic(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnets" "test" {
  id = huaweicloud_vpc_subnet.test.id
}
`, testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp))
}

func TestAccVpcSubnetsDataSource_ipv4ByCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetsDataSource_ipv4ByCidr(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "subnets.0.vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetsDataSource_ipv4ByCidr(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnets" "test" {
  cidr = huaweicloud_vpc_subnet.test.cidr
}
`, testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp))
}

func TestAccVpcSubnetsDataSource_ipv4ByName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetsDataSource_ipv4ByName(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "subnets.0.vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetsDataSource_ipv4ByName(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnets" "test" {
  name = huaweicloud_vpc_subnet.test.name
}
`, testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp))
}

func TestAccVpcSubnetsDataSource_ipv4ByVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetsDataSource_ipv4ByVpcId(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "subnets.0.vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetsDataSource_ipv4ByVpcId(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnets" "test" {
  vpc_id = huaweicloud_vpc_subnet.test.vpc_id
}
`, testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp))
}

func TestAccVpcSubnetsDataSource_ipv6Basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetsDataSource_ipv6Basic(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0.dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.ipv4_subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnets.0.ipv6_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "subnets.0.vpc_id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetsDataSource_ipv6Basic(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnets" "test" {
  id = huaweicloud_vpc_subnet.test.id
}
`, testAccVpcSubnetsDataSource_ipv6Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp string) string {
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

func testAccVpcSubnetsDataSource_ipv6Base(rName, cidr, gatewayIp string) string {
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
