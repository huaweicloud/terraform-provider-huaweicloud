---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_secret_rotate"
description: |-
  Manages a resource to rotate secret within HuaweiCloud.
---

# huaweicloud_csms_secret_rotate

Manages a resource to rotate secret within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not recover the rotated secret,
  but will only remove the resource information from the tfstate file.

-> This resource does not support rotation of common secrets. Running this resource will immediately perform credential
  rotation. Within the specified credentials, a new credential version will be created for encrypted storage of randomly
  generated credential values ​​in the background. Simultaneously, the newly created credential version will be marked as
  **SYSCURRENT**.

## Example Usage

```hcl
resource "huaweicloud_csms_secret_rotate" "test" {
  secret_name = "my_secret"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `secret_name` - (Required, String, NonUpdatable) Specifies the name of the secret to be rotated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as the task ID).

* `task_id` - The ID of the rotation task.

* `rotation_func_urn` - The URN of the rotation function.

* `operate_type` - The operation type of the rotation task.

* `task_time` - The time when the rotation task was created.

* `attempt_nums` - The number of attempts to rotate the secret.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
