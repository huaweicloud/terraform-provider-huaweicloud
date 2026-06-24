---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_create_retry_policy"
description: |-
  Manages a create or retry policy resource within HuaweiCloud.
---

# huaweicloud_secmaster_create_retry_policy

Manages a create or retry policy resource within HuaweiCloud.

-> This resource is a one-time action resource used to create or retry SecMaster policy. Deleting this resource will not
  change the status of the current SecMaster policy, but will only remove the resource information from the
  tfstate file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "version" {}
variable "account_scope" {}
variable "eps_scope" {}
variable "region_scope" {}
variable "defense_connection_id" {}

resource "huaweicloud_secmaster_create_retry_policy" "test" {
  workspace_id     = var.workspace_id
  action_type      = "create"
  version          = var.version
  block_target     = "192.168.0.0"
  policy_category  = "BLOCK"
  policy_direction = "INGRESS,EGRESS"
  account_scope    = var.account_scope
  eps_scope        = var.eps_scope
  region_scope     = var.region_scope

  block_age {
    is_block_ageing = false
  }

  defense_policy_list {
    defense_connection_id = var.defense_connection_id
  }

  policy_type {
    policy_type = "Source Ip"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the policy belongs.

* `action_type` - (Required, String, NonUpdatable) Specifies the operation type.
  The valid values are as follows:
  + **create**: Create a policy.
  + **retry**: Retry a policy.

* `version` - (Required, String, NonUpdatable) Specifies the service version, e.g. **25.5.0**.

* `block_age` - (Required, List, NonUpdatable) Specifies the block aging configuration.
  The [block_age](#create_retry_policy_block_age) structure is documented below.

* `block_target` - (Required, String, NonUpdatable) Specifies the policy object.

* `defense_policy_list` - (Required, List, NonUpdatable) Specifies the list of defense policies corresponding to the
  operation connection.
  The [defense_policy_list](#create_retry_policy_defense_policy_list) structure is documented below.

* `policy_category` - (Required, String, NonUpdatable) Specifies the policy category.
  The valid values are as follows:
  + **WHITE**: Whitelist (add objects such as IP to the whitelist).
  + **BLOCK**: Blocklist (add objects such as IP to the blocklist).

* `policy_type` - (Required, List, NonUpdatable) Specifies the block type.
  The [policy_type](#create_retry_policy_policy_type) structure is documented below.

* `retry_list` - (Optional, List, NonUpdatable) Specifies the list of policy IDs to retry.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the policy.

* `labels` - (Optional, String, NonUpdatable) Specifies the labels of the policy.

* `policy_direction` - (Optional, String, NonUpdatable) Specifies the policy direction.

* `account_scope` - (Optional, String, NonUpdatable) Specifies the account scope.

* `eps_scope` - (Optional, String, NonUpdatable) Specifies the enterprise project scope.

* `region_scope` - (Optional, String, NonUpdatable) Specifies the region scope.

<a name="create_retry_policy_block_age"></a>
The `block_age` block supports:

* `is_block_ageing` - (Required, Bool, NonUpdatable) Specifies whether to enable block aging.

* `block_ageing` - (Optional, String, NonUpdatable) Specifies the aging time in milliseconds.

<a name="create_retry_policy_defense_policy_list"></a>
The `defense_policy_list` block supports:

* `defense_connection_id` - (Required, String, NonUpdatable) Specifies the operation connection ID.

* `defense_connection_name` - (Optional, String, NonUpdatable) Specifies the operation connection name.

* `defense_connection_region_id` - (Optional, String, NonUpdatable) Specifies the region ID of the defense policy.

* `defense_connection_region_name` - (Optional, String, NonUpdatable) Specifies the region name of the defense policy.

* `defense_type` - (Optional, String, NonUpdatable) Specifies the defense service type.

* `target_enterprise_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

* `target_enterprise_name` - (Optional, String, NonUpdatable) Specifies the enterprise project name.

* `target_project_id` - (Optional, String, NonUpdatable) Specifies the project ID of the defense policy.

* `target_project_name` - (Optional, String, NonUpdatable) Specifies the project name of the defense policy.

<a name="create_retry_policy_policy_type"></a>
The `policy_type` block supports:

* `policy_type` - (Required, String, NonUpdatable) Specifies the block type.
  The valid values are as follows:
  + **User Name**
  + **Source Ip**
  + **Domain Name**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the task ID).
