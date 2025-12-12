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

func testDataSourceIdentityv5Tags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identityv5_resource_tags" "test" {
  resource_type = "user"
  resource_id   = huaweicloud_identityv5_user.user_1.id
  
  depends_on = [huaweicloud_identityv5_resource_tag.test]
}
`, testAccV5ResourceTag_basic(name))
}
