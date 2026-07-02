---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_instance_parameter"
description: |-
  Manages a GeminiDB instance parameter resource within HuaweiCloud.
---

# huaweicloud_geminidb_instance_parameter

Manages a GeminiDB instance parameter resource within HuaweiCloud.
This resource allows you to modify the value of a single parameter on a GeminiDB instance.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_geminidb_instance_parameter" "test" {
  instance_id = var.instance_id
  name        = "concurrent_reads"
  value       = "32"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, ForceNew) Specifies the GeminiDB instance ID.

* `name` - (Required, String, ForceNew) Specifies the parameter name.
  Only system-defined, modifiable parameters are valid. Other parameters will be ignored.

* `value` - (Required, String) Specifies the parameter value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in the format of `<instance_id>/<name>`.

* `restart_required` - Whether the instance needs to be restarted for the parameter change to take effect.

* `readonly` - Whether the parameter is read-only.

* `value_range` - The valid value range of the parameter.

* `type` - The parameter type. Valid values: `string`, `integer`, `boolean`, `list`, `float`.

* `description` - The parameter description.

## Timeouts

This resource does not have any timeouts.

## Import

The GeminiDB instance parameter can be imported using the `instance_id` and `name` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_geminidb_instance_parameter.test <instance_id>/<name>
```
