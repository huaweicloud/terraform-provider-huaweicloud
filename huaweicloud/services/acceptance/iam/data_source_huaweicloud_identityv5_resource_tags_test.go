package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV5ResourceTags_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identityv5_resource_tags.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5Tags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "resource_id"),
					resource.TestCheckResourceAttr(all, "resource_type", "user"),
					resource.TestCheckResourceAttr(all, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(all, "tags.key", "value"),
				),
			},
		},
	})
}

func testAccDataV5Tags_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identityv5_resource_tag" "test" {
  resource_type = "user"
  resource_id   = huaweicloud_identityv5_user.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testAccDataV5Tags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identityv5_resource_tags" "test" {
  resource_type = "user"
  resource_id   = huaweicloud_identityv5_user.test.id
  
  depends_on = [huaweicloud_identityv5_resource_tag.test]
}
`, testAccDataV5Tags_base(name))
}
