package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIECSecurityGroupDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("iec-secgroup-%s", acctest.RandString(5))
	description := "This is a test of iec security group"
	resourceName := "data.huaweicloud_iec_security_group.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecSecurityGroupV1Destory,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIECSecurityGroup_basic(rName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
		},
	})
}

func testAccDataSourceIECSecurityGroup_basic(rName, description string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_security_group" "my_group" {
  name        = "%s"
  description = "%s"
}

data "huaweicloud_iec_security_group" "by_name" {
  name = huaweicloud_iec_security_group.my_group.name
}
`, rName, description)
}
