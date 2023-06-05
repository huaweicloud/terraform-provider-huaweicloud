package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceOrganizations_basic(t *testing.T) {
	rName := "data.huaweicloud_organizations.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceOrganizations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "root_tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "root_tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccDatasourceOrganizations_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_organizations" "test" {
  depends_on = [huaweicloud_organizations.test]
}
`, testOrganizations_basic())
}
