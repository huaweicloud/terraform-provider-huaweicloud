package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceResourceShareInvitations_basic(t *testing.T) {
	dName := "data.huaweicloud_ram_resource_share_invitations.test"
	dc := acceptance.InitDataSourceCheck(dName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMShareInvitationId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceResourceShareInvitations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dName, "resource_share_invitations.0.id"),
					resource.TestCheckResourceAttrSet(dName, "resource_share_invitations.0.resource_share_id"),
					resource.TestCheckResourceAttrSet(dName, "resource_share_invitations.0.resource_share_name"),
					resource.TestCheckResourceAttrSet(dName, "resource_share_invitations.0.receiver_account_id"),
					resource.TestCheckResourceAttrSet(dName, "resource_share_invitations.0.sender_account_id"),
					resource.TestCheckResourceAttrSet(dName, "resource_share_invitations.0.status"),
					resource.TestCheckResourceAttrSet(dName, "resource_share_invitations.0.created_at"),
					resource.TestCheckResourceAttrSet(dName, "resource_share_invitations.0.updated_at"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceResourceShareInvitations_basic() string {
	return (`
data "huaweicloud_ram_resource_share_invitations" "test" {}

locals {
  status = data.huaweicloud_ram_resource_share_invitations.test.resource_share_invitations[0].status
}

data "huaweicloud_ram_resource_share_invitations" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ram_resource_share_invitations.filter_by_status.resource_share_invitations[*].status : 
    v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`)
}
