---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_image_signature_policy_execute"
description: |-
  Manages a SWR enterprise image signature policy execute resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_image_signature_policy_execute

Manages a SWR enterprise image signature policy execute resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "policy_id" {}

resource "huaweicloud_swr_enterprise_image_signature_policy_execute" "test" {
  instance_id    = var.instance_id
  namespace_name = var.namespace_name
  policy_id      = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String, NonUpdatable) Specifies the namespace name.

* `policy_id` - (Required, String, NonUpdatable) Specifies the policy ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `execution_id` - Specifies the execution ID.
