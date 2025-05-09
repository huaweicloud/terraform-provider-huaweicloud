---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_secrets"
description: |-
  Use this data source to get the list of CCI secrets within HuaweiCloud.
---

# huaweicloud_cciv2_secrets

Use this data source to get the list of CCI secrets within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_secrets" "test" {
  name = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `secrets` - The secrets list.
  The [secrets](#secrets) structure is documented below.

<a name="secrets"></a>
The `secrets` block supports:

* `annotations` - The annotations of the CCI secret.

* `creation_timestamp` - The creation timestamp of the CCI secret.

* `data` - The data of the CCI secret.

* `immutable` - The immutable of the CCI secret.

* `labels` - The labels of the CCI secret.

* `name` - The name of the CCI secret.

* `namespace` - The namespace of the CCI secret.

* `resource_version` - The resource version of the CCI secret.

* `string_data` - The string data of the CCI secret.

* `type` - The type of the CCI secret.

* `uid` - The uid of the CCI secret.
