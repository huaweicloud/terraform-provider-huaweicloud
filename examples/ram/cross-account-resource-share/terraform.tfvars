# Resource Share Configuration
resource_share_name        = "cross-account-vpc-share"
description                = "Share VPC resources with other accounts in the organization"

# Principals: Account IDs or Organization IDs to share resources with
# Should been replace the real account IDs
principals = [
  "01234567890123456789012345678901",
  "98765432109876543210987654321098"
]
