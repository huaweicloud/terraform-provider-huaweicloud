package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenterGroupMemberships_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_identitycenter_group_memberships.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterGroupMemberships_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "group_memberships.0.group_id"),
					resource.TestCheckResourceAttrSet(rName, "group_memberships.0.member_id.0.user_id"),
					resource.TestCheckResourceAttrSet(rName, "group_memberships.0.identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "group_memberships.0.membership_id"),
				),
			},
			{
				Config: testAccDatasourceIdentityCenterGroupMembershipForMember_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "group_memberships.0.group_id"),
					resource.TestCheckResourceAttrSet(rName, "group_memberships.0.member_id.0.user_id"),
					resource.TestCheckResourceAttrSet(rName, "group_memberships.0.identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "group_memberships.0.membership_id"),
				),
			},
		},
	})
}

func testAccDatasourceIdentityCenterGroupMemberships_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_identitycenter_group_memberships" "test"{
  depends_on        = [huaweicloud_identitycenter_group_membership.test]
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  group_id          = huaweicloud_identitycenter_group.test.id
}
`, testGroupMembership_basic(name))
}

func testAccDatasourceIdentityCenterGroupMembershipForMember_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_identitycenter_group_memberships" "test"{
  depends_on        = [huaweicloud_identitycenter_group_membership.test]
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  user_id           = huaweicloud_identitycenter_user.test.id
}
`, testGroupMembership_basic(name))
}
