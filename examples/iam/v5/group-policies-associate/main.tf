data "huaweicloud_identityv5_policies" "test" {
  policy_type = var.policy_type
}

locals {
  filtered_policies = [for policy in data.huaweicloud_identityv5_policies.test.policies : policy if contains(var.policy_names, policy.policy_name)]
}

resource "huaweicloud_identityv5_group" "test" {
  group_name  = var.group_name
  description = var.group_description
}

resource "huaweicloud_identityv5_policy_group_attach" "test" {
  count = length(local.filtered_policies)

  policy_id = try(local.filtered_policies[count.index].policy_id, null)
  group_id  = huaweicloud_identityv5_group.test.id
}
