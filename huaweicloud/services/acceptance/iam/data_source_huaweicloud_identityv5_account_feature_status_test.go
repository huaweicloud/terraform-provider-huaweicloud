package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAccountFeatureStatus_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identityv5_account_feature_status.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5AccountFeatureStatus_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "feature_status"),
				),
			},
		},
	})
}

const testAccDataV5AccountFeatureStatus_basic = `
data "huaweicloud_identityv5_account_feature_status" "test" {
  feature_name = "v5_console"
}
`
