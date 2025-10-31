package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccountAssignments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_account_assignments.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAccountAssignments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "account_assignments.#"),
					resource.TestCheckResourceAttr(dataSource, "account_assignments.0.account_id", acceptance.HW_IDENTITY_CENTER_ACCOUNT_ID),
					resource.TestCheckResourceAttrPair(dataSource, "account_assignments.0.permission_set_id",
						"huaweicloud_identitycenter_permission_set.test", "id"),
					resource.TestCheckResourceAttrPair(dataSource, "account_assignments.0.principal_id",
						"huaweicloud_identitycenter_user.test", "id"),
					resource.TestCheckResourceAttr(dataSource, "account_assignments.0.principal_type", "USER"),
				),
			},
		},
	})
}

func testDataSourceAccountAssignments_basic(name string) string {
	return fmt.Sprintf(`
%[1]s
 
data "huaweicloud_identitycenter_account_assignments" "test" {
  depends_on     = [huaweicloud_identitycenter_account_assignment.test]
  instance_id    = data.huaweicloud_identitycenter_instance.test.id
  principal_id   = huaweicloud_identitycenter_user.test.id
  principal_type = "USER"
}
`, testAccountAssignment_basic(name))
}
