package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcSubnetsDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	allSubnets := "data.huaweicloud_iec_vpc_subnets.all"
	siteSubnets := "data.huaweicloud_iec_vpc_subnets.site"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig_base(rName),
			},
			{
				Config: testAccSubnetsDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccSubnetsDataSourceID(allSubnets),
					testAccSubnetsDataSourceID(siteSubnets),
					resource.TestCheckResourceAttr(allSubnets, "subnets.#", "2"),
					resource.TestCheckResourceAttr(siteSubnets, "subnets.#", "1"),
					resource.TestCheckResourceAttrSet(siteSubnets, "subnets.0.id"),
					resource.TestCheckResourceAttrSet(siteSubnets, "subnets.0.site_info"),
					resource.TestCheckResourceAttrSet(siteSubnets, "subnets.0.status"),
				),
			},
		},
	})
}

func testAccSubnetsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find IEC VPC subnets data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("IEC VPC subnets data source ID not set")
		}

		return nil
	}
}

func testAccNetworkConfig_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "%s"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_1" {
  name       = "%s-1"
  cidr       = "192.168.0.0/24"
  vpc_id     = huaweicloud_iec_vpc.vpc_test.id
  site_id    = data.huaweicloud_iec_sites.sites_test.sites[0].id
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_2" {
  name       = "%s-2"
  cidr       = "192.168.1.0/24"
  vpc_id     = huaweicloud_iec_vpc.vpc_test.id
  site_id    = data.huaweicloud_iec_sites.sites_test.sites[1].id
  gateway_ip = "192.168.1.1"
}
`, rName, rName, rName)
}

func testAccSubnetsDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iec_vpc_subnets" "all" {
  vpc_id = huaweicloud_iec_vpc.vpc_test.id
}

data "huaweicloud_iec_vpc_subnets" "site" {
  vpc_id  = huaweicloud_iec_vpc.vpc_test.id
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`, testAccNetworkConfig_base(rName))
}
