---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_applications"
description: |-
  Use this data source to query the applications within HuaweiCloud.
---

# huaweicloud_apig_applications

Use this data source to query the applications within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "application_name" {}

data "huaweicloud_apig_applications" "test" {
  instance_id = var.instance_id
  name        = var.application_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the applications belong.

* `application_id` - (Optional, String) Specifies the ID of the application to be queried.

* `name` - (Optional, String) Specifies the name of the application to be queried.

* `app_key` - (Optional, String) Specifies the key of the application to be queried.

* `created_by` - (Optional, String) Specifies the creator of the application to be queried.  
  The valid values are as follows:
  + **USER**: The user created.
  + **MARKET**: The cloud store allocation.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - All applications that match the filter parameters.
  The [applications](#attrblock_applications) structure is documented below.

<a name="attrblock_applications"></a>
The `applications` block supports:

* `id` - The ID of the application.

* `name` - The name of the application.

* `status` - The status of the application.

* `description` - The description of the application.

* `app_key` - The key of the application.

* `app_secret` - The secret of the application.

* `app_type` - The type of the application.

* `bind_num` - The number of bound APIs.

* `created_by` - The creator of the application.

* `created_at` - The creation time of the application, in RFC3339 format.

* `updated_at` - The latest update time of the application, in RFC3339 format.
