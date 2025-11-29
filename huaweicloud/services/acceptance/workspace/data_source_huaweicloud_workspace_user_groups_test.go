package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataUserGroups_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_workspace_user_groups.all"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "groups.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dcName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dcName, "groups.0.platform_type"),
				),
			},
		},
	})
}

func testAccDataUserGroups_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_user_group" "test" {
  name        = "%[1]s"
  type        = "LOCAL"
  description = "Created by terraform script"
}

data "huaweicloud_workspace_user_groups" "all" {
  depends_on = [
    huaweicloud_workspace_user_group.test,
  ]
}
`, name)
}
