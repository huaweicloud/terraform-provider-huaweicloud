package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityv5ResourceTags_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_resource_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceIdentityv5Tags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_id"),
					resource.TestCheckResourceAttr(dataSourceName, "resource_type", "user"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.key", "value"),
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

func testDataSourceIdentityv5Tags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identityv5_resource_tags" "test" {
  resource_type = "user"
  resource_id   = huaweicloud_identityv5_user.test.id
  
  depends_on = [huaweicloud_identityv5_resource_tag.test]
}
`, testAccDataV5Tags_base(name))
}
