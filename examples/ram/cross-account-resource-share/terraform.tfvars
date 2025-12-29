# Resource Share Configuration
resource_share_name        = "cross-account-vpc-share"
description                = "Share VPC resources with other accounts in the organization"

# Principals: Account IDs or Organization IDs to share resources with
# Should been replace the real account IDs
principals = [
  "01234567890123456789012345678901",
  "98765432109876543210987654321098"
]

# The list of URNs of one or more resources to be shared.
# Should been replace the real URNs
resource_urns = [
  "vpc:cn-north-4:8f06724e5c6f41f59d3e2f3ad897bb4d:subnet:5de72eeb-7977-4602-8186-8766982d9bcc",
]

# The list of RAM permissions associated with the resource share
# Should been replace the real permission IDs
permission_ids = [
  "f5153698-ca8b-4b3c-a839-13ff71f67885"
]
