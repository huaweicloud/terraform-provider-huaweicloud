# Share (owner).
provider "huaweicloud" {
  alias = "owner"

  region     = var.region_name
  access_key = var.owner_ak
  secret_key = var.owner_sk
}

# Other account (principal).
provider "huaweicloud" {
  alias = "principal"

  region     = var.region_name
  access_key = var.principal_ak
  secret_key = var.principal_sk
}

data "huaweicloud_er_availability_zones" "test" {
  provider = huaweicloud.owner
}

# Owner creates an ER instance to share.
resource "huaweicloud_er_instance" "test" {
  provider = huaweicloud.owner

  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name        = var.er_instance_name
  asn         = "64512"
  description = "Create an ER instace to share"

  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = false
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

  name = var.vpc_name
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  provider = huaweicloud.principal

  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
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

