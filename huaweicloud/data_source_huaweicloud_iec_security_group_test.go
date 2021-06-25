package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccIECSecurityGroupDataSource_basic(t *testing.T) {
	rName := fmtp.Sprintf("iec-secgroup-%s", acctest.RandString(5))
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
	return fmtp.Sprintf(`
resource "huaweicloud_iec_security_group" "my_group" {
  name        = "%s"
  description = "%s"
}

data "huaweicloud_iec_security_group" "by_name" {
  name = huaweicloud_iec_security_group.my_group.name
}
`, rName, description)
}
