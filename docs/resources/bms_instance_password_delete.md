---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_instance_password_delete"
description: |-
  Manages a BMS instance password delete resource within HuaweiCloud.
---

# huaweicloud_bms_instance_password_delete

Manages a BMS instance password delete resource within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}

resource "huaweicloud_bms_instance_password_delete" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the provider-level
  region will be used. Changing this creates a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the BMS ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the BMS instance ID.
