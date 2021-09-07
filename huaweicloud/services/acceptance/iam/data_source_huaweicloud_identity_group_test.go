package iam

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIdentityGroupDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_identity_group.test"
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupDataSource_by_name(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityGroupDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckIdentityGroupDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Group data source ID not set")
		}

		return nil
	}
}

func testAccIdentityGroupDataSource_by_name(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name        = "%s"
  description = "A ACC test group"
}

data "huaweicloud_identity_group" "test" {
  name = huaweicloud_identity_group.group_1.name
  
  depends_on = [
    huaweicloud_identity_group.group_1
  ]
}
`, rName)
}
