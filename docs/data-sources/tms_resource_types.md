---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_resource_types"
description: ""
---

# huaweicloud_tms_resource_types

Using this data source to query supported resource types information that used to manage resource tags within
HuaweiCloud.

## Example Usage

```hcl
variable "supported_region" {}

data "huaweicloud_tms_resource_types" "test" {
  region = var.supported_region
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region name used to filter resource types information

* `service_name` - (Optional, String) Specifies the service name used to filter resource types information.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `types` - All resource types that match the filter parameters.
  The [types](#tms_resource_types) structure is documented below.

<a name="tms_resource_types"></a>
The `types` block supports:

* `name` - The resource type name.

* `is_global` - Whether the resource corresponding to this type is a global resource.

* `display_name` - The service display name of the resource type.

* `service_name` - The name of the service to which the resource type belong.
