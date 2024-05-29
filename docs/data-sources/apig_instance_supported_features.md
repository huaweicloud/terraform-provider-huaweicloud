---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_supported_features"
description: |-
  Use this data source to get the list of supported features under the APIG instance within HuaweiCloud.
---

# huaweicloud_apig_instance_supported_features

Use this data source to get the list of supported features under the APIG instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_apig_instance_supported_features" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) The ID of the dedicated instance to which the features belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `features` - All names of the supported features that match the filter parameters.
