---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_execute_policy"
description: |-
  Manages an AS execute policy resource within HuaweiCloud.
---

# huaweicloud_as_execute_policy

Manages an AS execute policy resource within HuaweiCloud.

-> 1. The current resource is a one-time resource, and destroying this resource will not change the current status.
<br/>2. If the policy belongs the AS group policy, before creating the resource, you need to ensure the AS group status
  and the AS group policy status is **INSERVICE**.
<br/>3. If the policy belongs the AS bandwidth policy, before creating the resource, you need to ensure the AS
  bandwidth policy status is **INSERVICE**.
<br/>4. After creating the resource, you can use the datasource `huaweicloud_as_policy_execute_logs` to querry
  the policy execute result and details.

## Example Usage

```hcl
variable "scaling_policy_id" {}

resource "huaweicloud_as_execute_policy" "test" {
  scaling_policy_id = var.scaling_policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `scaling_policy_id` - (Required, String, NonUpdatable) Specifies the AS group policy ID or AS bandwidth policy ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
