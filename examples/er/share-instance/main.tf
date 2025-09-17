data "huaweicloud_er_availability_zones" "test" {
  provider = huaweicloud.owner
}

# Owner creates an ER instance to share.
resource "huaweicloud_er_instance" "test" {
  provider = huaweicloud.owner

  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name        = var.instance_name
  asn         = var.instance_asn
  description = var.instance_description

  enable_default_propagation     = var.instance_enable_default_propagation
  enable_default_association     = var.instance_enable_default_association
  auto_accept_shared_attachments = var.instance_auto_accept_shared_attachments
}

data "huaweicloud_ram_resource_permissions" "test" {
  provider = huaweicloud.owner

  resource_type = "er:instances"

  depends_on = [huaweicloud_er_instance.test]
}

# Owner creates a RAM shared resource to initiate a shared ER request.
resource "huaweicloud_ram_resource_share" "test" {
  provider = huaweicloud.owner

  name          = var.resource_share_name
  principals    = [var.principal_account_id]
  resource_urns = ["er:${var.region_name}:${var.owner_account_id}:instances:${huaweicloud_er_instance.test.id}"]

  permission_ids = data.huaweicloud_ram_resource_permissions.test.permissions[*].id
}

# Principal queries the shared ER requests that need to be accepted.
data "huaweicloud_ram_resource_share_invitations" "test" {
  provider = huaweicloud.principal

  status = "pending"

  depends_on = [huaweicloud_ram_resource_share.test]
}

# Principal (ER instance acceptor) to accept request from owner shared ER.
resource "huaweicloud_ram_resource_share_accepter" "test" {
  provider = huaweicloud.principal

  resource_share_invitation_id = try([for v in data.huaweicloud_ram_resource_share_invitations.test.resource_share_invitations : v.id if v.resource_share_id == huaweicloud_ram_resource_share.test.id][0], "")
  action                       = "accept"

  # After accepting the request, querying data.huaweicloud_ram_resource_share_invitations again will be empty.
  # This resource is a one-time resource. Add ignore_changes to prevent resource changes when executing terraform plan.
  lifecycle {
    ignore_changes = [
      resource_share_invitation_id,
    ]
  }
}

resource "huaweicloud_vpc" "test" {
  provider = huaweicloud.principal

  name = var.principal_vpc_name
  cidr = var.principal_vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  provider = huaweicloud.principal

  vpc_id     = huaweicloud_vpc.test.id
  name       = var.principal_subnet_name
  cidr       = var.principal_subnet_cidr != "" ? var.principal_subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.principal_subnet_gateway_ip != "" ? var.principal_subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

# Principal creates a VPC attachment.
resource "huaweicloud_er_vpc_attachment" "test" {
  provider = huaweicloud.principal

  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  name        = var.attachment_name

  depends_on = [huaweicloud_ram_resource_share_accepter.test]
}

# The owner accepts attachment from principals.
resource "huaweicloud_er_attachment_accepter" "test" {
  provider = huaweicloud.owner

  instance_id   = huaweicloud_er_instance.test.id
  attachment_id = huaweicloud_er_vpc_attachment.test.id
  action        = "accept"
}
