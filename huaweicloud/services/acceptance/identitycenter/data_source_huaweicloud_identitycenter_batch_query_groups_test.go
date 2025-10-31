package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenterBatchQueryGroups_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_identitycenter_batch_query_groups.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterBatchQueryGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.name"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.description"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.updated_at"),
				),
			},
		},
	})
}

func testAccDatasourceIdentityCenterBatchQueryGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_identitycenter_groups" "test"{
  depends_on        = [huaweicloud_identitycenter_group.test]
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}

data "huaweicloud_identitycenter_batch_query_groups" "test"{
  depends_on        = [huaweicloud_identitycenter_group.test]
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  group_ids         = [for group in data.huaweicloud_identitycenter_groups.test.groups: group.id]
}
`, testIdentityCenterGroup_basic(name))
}
