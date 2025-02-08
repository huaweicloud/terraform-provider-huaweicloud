---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_instance_features"
description: |-
  Use this data source to get the advanced features of the GaussDB OpenGauss instance.
---

# huaweicloud_gaussdb_opengauss_instance_features

Use this data source to get the advanced features of the GaussDB OpenGauss instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_instance_features" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `features` - Indicates the advanced features.

  The [features](#features_struct) structure is documented below.

<a name="features_struct"></a>
The `features` block supports:

* `name` - Indicates the feature name.

* `status` - Indicates whether the feature is enabled.

* `description` - Indicates the feature description.

* `type` - Indicates the feature value type.
  The value can be: **integer**, **string**, **boolean**.

* `value` - Indicates the feature value.

* `range` - Indicates the feature value range.

* `range_description` - Indicates the feature scope description.
