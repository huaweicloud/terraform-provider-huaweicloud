---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_namespace"
description: |-
  Manages a CCI v2 namespace resource within HuaweiCloud.
---

# huaweicloud_cciv2_namespace

Manages a CCI v2 namespace resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_cciv2_namespace" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the unique name of the namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is also the name of the namespace.

* `api_version` - The API version of the namespace.

* `kind` - The kind of the namespace.

* `annotations` - The annotations of the namespace.

* `labels` - The labels of the namespace.

* `creation_timestamp` - The creation timestamp of the namespace.

* `resource_version` - The resource version of the namespace.

* `uid` - The uid of the namespace.

* `finalizers` - The finalizers of the namespace.

* `status` - The status of the namespace.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 3 minutes.

## Import

The CCI v2 namespace can be imported using `name`, e.g.

```bash
$ terraform import huaweicloud_cciv2_namespace.test <name>
```
