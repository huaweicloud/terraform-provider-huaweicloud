package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceResourceShareInvitations_basic(t *testing.T) {
	var (
		dName = "data.huaweicloud_ram_resource_share_invitations.test"
		dc    = acceptance.InitDataSourceCheck(dName)

		byResourceShareIDs   = "data.huaweicloud_ram_resource_share_invitations.filter_by_resource_share_ids"
		dcByResourceShareIDs = acceptance.InitDataSourceCheck(byResourceShareIDs)

		byResourceShareInvitationIDs   = "data.huaweicloud_ram_resource_share_invitations.filter_by_resource_share_invitation_ids"
		dcByResourceShareInvitationIDs = acceptance.InitDataSourceCheck(byResourceShareInvitationIDs)

		byStatus   = "data.huaweicloud_ram_resource_share_invitations.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please first ensure that there is a sharing invitation under the current account.
			acceptance.TestAccPreCheckRAMShareInvitationId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceResourceShareInvitations_basic,
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

					dcByResourceShareIDs.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_share_ids_filter_useful", "true"),

					dcByResourceShareInvitationIDs.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_share_invitation_ids_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

var testAccDatasourceResourceShareInvitations_basic = `
data "huaweicloud_ram_resource_share_invitations" "test" {}

# Filter by resource_share_ids
locals {
  resource_share_ids = split(",", data.huaweicloud_ram_resource_share_invitations.test.resource_share_invitations[0].resource_share_id)
}

data "huaweicloud_ram_resource_share_invitations" "filter_by_resource_share_ids" {
  resource_share_ids = local.resource_share_ids
}

locals {
  resource_share_ids_filter_result = [
    for v in data.huaweicloud_ram_resource_share_invitations.filter_by_resource_share_ids.resource_share_invitations[*].
    resource_share_id : contains(local.resource_share_ids, v)
  ]
}

output "is_resource_share_ids_filter_useful" {
  value = length(local.resource_share_ids_filter_result) > 0 && alltrue(local.resource_share_ids_filter_result)
}

# Filter by resource_share_invitation_ids
locals {
  resource_share_invitation_ids = split(",", data.huaweicloud_ram_resource_share_invitations.test.resource_share_invitations[0].id)
}

data "huaweicloud_ram_resource_share_invitations" "filter_by_resource_share_invitation_ids" {
  resource_share_invitation_ids = local.resource_share_invitation_ids
}

locals {
  resource_share_invitation_ids_filter_result = [
    for v in data.huaweicloud_ram_resource_share_invitations.filter_by_resource_share_invitation_ids.
    resource_share_invitations[*].id : contains(local.resource_share_invitation_ids, v)
  ]
}

output "is_resource_share_invitation_ids_filter_useful" {
  value = length(local.resource_share_invitation_ids_filter_result) > 0 && alltrue(local.resource_share_invitation_ids_filter_result)
}

# Filter by status
locals {
  status = data.huaweicloud_ram_resource_share_invitations.test.resource_share_invitations[0].status
}

data "huaweicloud_ram_resource_share_invitations" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ram_resource_share_invitations.filter_by_status.resource_share_invitations[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`
