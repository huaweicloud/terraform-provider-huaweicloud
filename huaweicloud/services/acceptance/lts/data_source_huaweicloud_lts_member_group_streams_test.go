package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMemberGroupStreams_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_lts_member_group_streams.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSMemberAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMemberGroupStreams_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "groups.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.log_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.log_group_id"),
					resource.TestMatchResourceAttr(dataSourceName, "groups.0.log_streams.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.log_streams.0.log_stream_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.log_streams.0.log_stream_id"),
				),
			},
		},
	})
}

func testAccDataSourceMemberGroupStreams_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_lts_member_group_streams" "test" {
  member_account_id = "%s"
}
`, acceptance.HW_LTS_LOG_CONVERGE_MEMBER_ACCOUNT_ID)
}
