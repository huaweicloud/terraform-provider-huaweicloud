package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIamIdentityV5AccountSummary_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identityv5_account_summary.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceIamIdentityV5AccountSummary_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "agencies_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "attached_policies_per_agency_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "attached_policies_per_group_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "attached_policies_per_user_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "groups_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "policies"),
					resource.TestCheckResourceAttrSet(dataSource, "policies_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "groups"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_size_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "root_user_mfa_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "users"),
					resource.TestCheckResourceAttrSet(dataSource, "users_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "versions_per_policy_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "agencies"),
				),
			},
		},
	})
}

const testDataSourceDataSourceIamIdentityV5AccountSummary_basic = `
data "huaweicloud_identityv5_account_summary" "test" {}
`
