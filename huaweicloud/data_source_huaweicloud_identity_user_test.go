package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIdentityUserDataSource_basic(t *testing.T) {
	datasourceName := "data.huaweicloud_identity_user.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityUserDataSourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckIdentityUserDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find identity user data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Identity user dIata source ID not set")
		}

		return nil
	}
}

func testAccIdentityUserDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_identity_user" "test" {
  name = huaweicloud_identity_user.user_1.name
}
`, testAccIdentityV3User_basic(rName))
}
