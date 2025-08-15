package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppHdaUpgradeRecords_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_app_hda_upgrade_records.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppHdaUpgradeRecords_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					// The server_id, machine_name, server_name, server_group_name are empty, unknown these fields which
					// scenarios will be returned currently.
					resource.TestMatchResourceAttr(all, "records.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "records.0.sid"),
					resource.TestCheckResourceAttrSet(all, "records.0.current_version"),
					resource.TestCheckResourceAttrSet(all, "records.0.target_version"),
					resource.TestCheckResourceAttrSet(all, "records.0.upgrade_status"),
					resource.TestCheckResourceAttrSet(all, "records.0.upgrade_time"),
				),
			},
		},
	})
}

const testAccDataAppHdaUpgradeRecords_basic = `
data "huaweicloud_workspace_app_hda_upgrade_records" "test" {}
`
