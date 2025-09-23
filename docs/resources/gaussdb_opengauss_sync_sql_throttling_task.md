---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_sync_sql_throttling_task"
description: |-
  Manages a GaussDB OpenGauss sync SQL throttling task resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_sync_sql_throttling_task

Manages a GaussDB OpenGauss sync SQL throttling task resource within HuaweiCloud.

-> This resource is only a one-time action resource for operating the API.
Deleting this resource will not clear the corresponding request record,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_opengauss_sync_sql_throttling_task" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB OpenGauss instance.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is `instance_id`.
