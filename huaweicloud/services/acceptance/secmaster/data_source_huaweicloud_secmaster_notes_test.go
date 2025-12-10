package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNotes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_notes.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a collector channel group.
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNotes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.data.0.content"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_deleted"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.marked_note"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.note_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.user.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.user.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.war_room_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.workspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.content"),

					resource.TestCheckOutput("is_war_room_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceNotes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_notes" "test" {
  workspace_id = "%[1]s"
}

# Filter by war_room_id
locals {
  war_room_id = data.huaweicloud_secmaster_notes.test.data[0].war_room_id
}

data "huaweicloud_secmaster_notes" "war_room_id_filter" {
  workspace_id = "%[1]s"
  war_room_id  = local.war_room_id
}

output "is_war_room_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_notes.war_room_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_notes.war_room_id_filter.data[*].war_room_id : v == local.war_room_id]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
