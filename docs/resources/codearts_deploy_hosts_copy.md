---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_hosts_copy"
description: |-
  Manages a CodeArts deploy hosts copy resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_hosts_copy

Manages a CodeArts deploy hosts copy resource within HuaweiCloud.

## Example Usage

```hcl
variable "source_group_id" {}
variable "host_uuids" {}
variable "target_group_id" {}

resource "huaweicloud_codearts_deploy_hosts_copy" "test" {
  source_group_id = var.source_group_id
  host_uuids      = var.host_uuids
  target_group_id = var.target_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `source_group_id` - (Required, String, ForceNew) Specifies the source group ID.
  Changing this creates a new resource.

* `host_uuids` - (Required, List, ForceNew) Specifies the host IDs list.
  Changing this creates a new resource.

* `target_group_id` - (Required, String, ForceNew) Specifies the target group ID.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
