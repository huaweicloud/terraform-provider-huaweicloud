---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_replication_policy_execute"
description: |-
  Manages a SWR enterprise replication policy execute resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_replication_policy_execute

Manages a SWR enterprise replication policy execute resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "policy_id" {}

resource "huaweicloud_swr_enterprise_replication_policy_execute" "test" {
  instance_id = var.instance_id
  policy_id   = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `policy_id` - (Required, Int, NonUpdatable) Specifies the policy ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `execution_id` - Indicates the execution ID.
