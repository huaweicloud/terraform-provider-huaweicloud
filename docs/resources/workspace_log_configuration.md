---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_log_configuration"
description: |-
  Use this resource to manage Workspace log configuration for user events within HuaweiCloud.
---

# huaweicloud_workspace_log_configuration

Use this resource to manage Workspace log configuration for user events within HuaweiCloud.

## Example Usage

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_workspace_log_configuration" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the log configuration is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `log_group_id` - (Required, String) Specifies the ID of the log group.  
  The length ranges from `1` to `64` characters.

* `log_stream_id` - (Required, String) Specifies the ID of the log stream.  
  The length ranges from `1` to `64` characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the project ID of the region.

## Import

The log configuration can be imported using the `id` (it can be any UUID), e.g.

```bash
$ terraform import huaweicloud_workspace_log_configuration.test <id>
```
