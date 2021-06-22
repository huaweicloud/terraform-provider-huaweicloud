package huaweicloud

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccEnterpriseProjectDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_enterprise_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProjectDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnterpriseProjectDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "default"),
					resource.TestCheckResourceAttr(resourceName, "id", "0"),
				),
			},
		},
	})
}

func testAccCheckEnterpriseProjectDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find enterprise project data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("enterprise project data source ID not set ")
		}

		return nil
	}
}

const testAccEnterpriseProjectDataSource_basic = `
data "huaweicloud_enterprise_project" "test" {
  name = "default"
}
`
