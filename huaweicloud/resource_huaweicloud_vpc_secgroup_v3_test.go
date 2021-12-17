package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestResourceVpcSecGroupV3(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVpcSecGroup,
				Check:  resource.ComposeTestCheckFunc(
				//resource.TestCheckResourceAttr("huaweicloud_vpc_secgroup_v3.test", "region", "cn-north-4"),
				),
			},
			{
				Config: testVpcSecGroupUpdate,
				Check:  resource.ComposeTestCheckFunc(
				//resource.TestCheckResourceAttr("huaweicloud_vpc_secgroup_v3.test", "region", "cn-north-4"),
				),
			},
		},
	})
}

var testVpcSecGroup = `
resource "huaweicloud_vpc_secgroup_v3" "test"  {
 	region = "cn-north-4"
	dry_run = false
	name = "aaa"
	description = "123"
	enterprise_project_id = "0"
}
`
var testVpcSecGroupUpdate = `
resource "huaweicloud_vpc_secgroup_v3" "test"  {
 	region = "cn-north-4"
	dry_run = false
	name = "a"
	description = "12"
}
`

func testAccCheckVpcSecGroupDestroy(s *terraform.State) error {
	return nil
}
