package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestResourceVpcAddressGroupV3(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcAddressGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVpcAdressGroupV3Config,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testVpcAdressGroupV3Update,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

var testVpcAdressGroupV3Config = `
resource "huaweicloud_vpc_address_group" "test" {
	dry_run = false
	name = "test1"
	ip_version = 4
	description  =  "vpc test"
	ip_set	=	[
		"192.168.5.0/24",
		"192.168.9.0/24"
	]
}
`

var testVpcAdressGroupV3Update = `
resource "huaweicloud_vpc_address_group" "test" {
	dry_run = false
	name = "test02"
	ip_version = 4
	description  =  "vpc test02"
	ip_set	=	[
		"192.168.5.0/24",
		"192.168.9.0/24"
	]
}
`

func testAccCheckVpcAddressGroupDestroy(s *terraform.State) error {
	return nil
}
