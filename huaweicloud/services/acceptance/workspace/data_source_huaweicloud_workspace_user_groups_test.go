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

		all = "data.huaweicloud_workspace_user_groups.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "groups.0.id"),
					resource.TestCheckResourceAttrSet(all, "groups.0.name"),
					resource.TestCheckResourceAttrSet(all, "groups.0.platform_type"),
					resource.TestMatchResourceAttr(all, "groups.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "groups.0.user_quantity"),
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

# Without any filter parameter.
data "huaweicloud_workspace_user_groups" "all" {
  depends_on = [
    huaweicloud_workspace_user_group.test,
  ]
}
`, name)
}
