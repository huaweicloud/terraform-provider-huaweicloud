---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_async_log_configuration"
description: |-
  Manages a FunctionGraph async log configuration resource under a specified region within HuaweiCloud.
---

# huaweicloud_fgs_async_log_configuration

Manages a FunctionGraph async log configuration resource under a specified region within HuaweiCloud.

-> After the resource is deployed, a log group named `function-async-log-group-{project_id}` will be automatically
   created in the corresponding region, and a log stream named `function-async-log-stream-{project_id}` will be
   automatically created in it.

## Example Usage

```hcl
resource "huaweicloud_fgs_async_log_configuration" "test" {
  force_delete = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the async log configuration is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `force_delete` - (Optional, Bool) Specifies whether to force delete the LTS resources corresponding to the async log
  configuration.  
  Defaults to **false**.

  -> Forced delete action will clear all data under the log group corresponding to all current asynchronous logs.<br>
     Please back up the data before deletion.

* `max_retries` - (Optional, Int) Specifies the maximum number of retries for the create operation when encountering
  internal service errors.  
  Defaults to **3**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The (random) resource ID.

* `group_id` - The LTS log group ID used to manage async status logs.

* `group_name` - The LTS log group name used to manage async status logs.

* `stream_id` - The LTS log stream ID used to manage async status logs.

* `stream_name` - The LTS log stream name used to manage async status logs.

## Import

Async log configuration can be imported using the `id` (any string ID format is allowed), e.g.

```bash
$ terraform import huaweicloud_fgs_async_log_configuration.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes includes: `force_delete`, `max_retries`.
It is generally recommended running `terraform plan` after importing this resource.
You can then decide if changes should be applied to the resource, or the script definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_fgs_async_log_configuration" "test" {
  ...

  lifecycle {
    ignore_changes = [
      force_delete,
      max_retries,
    ]
  }
}
```
