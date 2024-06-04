---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_features"
description: |-
  Use this data source to get the list of the features under the APIG instance within HuaweiCloud.
---

# huaweicloud_apig_instance_features

Use this data source to get the list of the features under the APIG instance within HuaweiCloud.

## Example Usage

```hcl
variable instance_id {}

data "huaweicloud_apig_instance_features" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specified the ID of the dedicated instance to which the features belong.

* `name` - (Optional, String) Specified the name of the feature.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `features` - All instance features that match the filter parameters.
  The [features](#instance_features) structure is documented below.

<a name="instance_features"></a>
The `features` block supports:

* `id` - The ID of the feature.

* `name` - The name of the feature.

* `enabled` - Whether the feature is enabled.

* `config` - The detailed configuration of the instance feature.

* `updated_at` - The latest update time of the feature, in RFC3339 format.
