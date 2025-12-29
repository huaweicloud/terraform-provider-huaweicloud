# Query pending resource share invitations
data "huaweicloud_ram_resource_share_invitations" "test" {
  resource_share_ids = var.resource_share_ids
  status             = "pending"
}

# Process invitations: accept or reject based on configuration
resource "huaweicloud_ram_resource_share_accepter" "test" {
  count = length(data.huaweicloud_ram_resource_share_invitations.test.resource_share_invitations)

  resource_share_invitation_id = data.huaweicloud_ram_resource_share_invitations.test.resource_share_invitations[count.index].id
  action                       = var.action
}
