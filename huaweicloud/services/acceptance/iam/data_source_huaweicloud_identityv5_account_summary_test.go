package iam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV5AccountSummary_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identityv5_account_summary.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5AccountSummary_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "agencies_quota"),
					resource.TestCheckResourceAttrSet(dcName, "attached_policies_per_agency_quota"),
					resource.TestCheckResourceAttrSet(dcName, "attached_policies_per_group_quota"),
					resource.TestCheckResourceAttrSet(dcName, "attached_policies_per_user_quota"),
					resource.TestCheckResourceAttrSet(dcName, "groups_quota"),
					resource.TestMatchResourceAttr(dcName, "policies", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "policies_quota"),
					resource.TestMatchResourceAttr(dcName, "groups", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "policy_size_quota"),
					resource.TestCheckResourceAttrSet(dcName, "root_user_mfa_enabled"),
					resource.TestMatchResourceAttr(dcName, "users", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "users_quota"),
					resource.TestCheckResourceAttrSet(dcName, "versions_per_policy_quota"),
					resource.TestCheckResourceAttrSet(dcName, "agencies"),
				),
			},
		},
	})
}

const testAccDataV5AccountSummary_basic = `
data "huaweicloud_identityv5_account_summary" "test" {}
`
