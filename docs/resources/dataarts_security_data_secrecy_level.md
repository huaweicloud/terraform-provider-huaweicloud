---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_data_secrecy_level"
description: |-
  Manages DataArts Security data secrecy level resource within HuaweiCloud.
---


# huaweicloud_dataarts_security_data_secrecy_level

Manages DataArts Security data secrecy level resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "secrecy_level_name" {}

resource "huaweicloud_dataarts_security_data_secrecy_level" "test" {
  workspace_id = var.workspace_id
  name         = var.secrecy_level_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the data secrecy level belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the data secrecy level.
  The valid length is limited from `1` to `128`, only Chinese and English characters, digits and underscores (_) are
  allowed. Changing this creates a new resource.

* `description` - (Optional, String) Specifies the description of the data secrecy level.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `level_number` - The security level of the data secrecy level. The larger value means the higher security level.

* `created_by` - The creator of the data secrecy level.

* `updated_by` - The user who latest updated the data secrecy level.

* `created_at` - The creation time of the data secrecy level.

* `updated_at` - The latest update time of the data secrecy level.

## Import

The resource can be imported using `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dataarts_security_data_secrecy_level.test <workspace_id>/<id>
```
