---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_config_maps"
description: |-
  Use this data source to get the list of CCI config maps within HuaweiCloud.
---

# huaweicloud_cciv2_config_maps

Use this data source to get the list of CCI config maps within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_config_maps" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace of the CCI config maps.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `config_maps` - The config maps.
  The [config_maps](#config_maps) structure is documented below.

<a name="config_maps"></a>
The `config_maps` block supports:

* `annotations` - The annotations.

* `binary_data` - The binary data.

* `creation_timestamp` - The creation time.

* `data` - The data.

* `immutable` - The immutable.

* `labels` - The labels.

* `name` - The name.

* `namespace` - The namespace.

* `resource_version` - The resource version.

* `uid` - The uid.
