---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_config_map"
description: |-
  Manages a CCI v2 ConfigMap resource within HuaweiCloud.
---

# huaweicloud_cciv2_config_map

Manages a CCI v2 ConfigMap resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "name" {}
variable "data" {}

data "huaweicloud_cciv2_config_map" "test" {
  namespace = var.namespace
  name      = var.name
  data      = var.data
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the CCI ConfigMap.

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace.

* `binary_data` - (Optional, Map) Specifies the binary data of the CCI ConfigMap.

* `data` - (Optional, Map) Specifies the data of the CCI ConfigMap.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `annotations` - The annotations of the CCI ConfigMap.

* `api_version` - The API version of the CCI ConfigMap.

* `creation_timestamp` - The creation timestamp of the CCI ConfigMap.

* `immutable` - The immutable of the CCI ConfigMap.

* `kind` - The kind of the CCI ConfigMap.

* `labels` - The labels of the CCI ConfigMap.

* `resource_version` - The resource version of the CCI ConfigMap.

* `uid` - The uid of the CCI ConfigMap.

## Import

The CCI v2 ConfigMap can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cciv2_config_map.test <namespace>/<name>
```
