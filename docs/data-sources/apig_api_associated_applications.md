---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_associated_applications"
description: |-
  Use this data source to query the applications associated with the specified API within HuaweiCloud.
---

# huaweicloud_apig_api_associated_applications

Use this data source to query the applications associated with the specified API within HuaweiCloud.

## Example Usage

### Query the contents of all applications bound to the current API

```hcl
variable "instance_id" {}
variable "associated_api_id" {}

data "huaweicloud_apig_api_associated_applications" "test" {
  instance_id = var.instance_id
  api_id      = var.associated_api_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the associated applications.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the applications belong.

* `api_id` - (Required, String) Specifies the ID of the API bound to the application.

* `application_id` - (Optional, String) Specifies the ID of the application.

* `name` - (Optional, String) Specifies the name of the application.

* `env_id` - (Optional, String) Specifies the ID of the environment where the API is published.

* `env_name` - (Optional, String) Specifies the name of the environment where the API is published.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - All applications that match the filter parameters.
  The [applications](#api_associated_applications) structure is documented below.

<a name="api_associated_applications"></a>
The `applications` block supports:

* `id` - The ID of the application.

* `name` - The name of the application.

* `description` - The description of the application.

* `env_id` - The ID of the environment where the API is published.

* `env_name` - The name of the environment where the API is published.

* `bind_id` - The bind ID.

* `bind_time` - The time that the application is bound to the API.
