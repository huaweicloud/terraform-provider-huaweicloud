---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_sessions_close"
description: |-
  Manages a resource to close all sessions of all nodes on a GeminiDB Redis instance within HuaweiCloud.
---

# huaweicloud_geminidb_sessions_close

Manages a resource to close all sessions of all nodes on a GeminiDB Redis instance within HuaweiCloud.

-> 1. This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.
  <br/>2. This resource only supports GeminiDB Redis instance.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_geminidb_sessions_close" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to close sessions. If omitted, the provider-level region
  will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the GeminiDB Redis instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to the instance ID.
