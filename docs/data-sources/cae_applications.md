---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_applications"
description: |-
  Use this data source to get the list of CAE applications within HuaweiCloud.
---

# huaweicloud_cae_applications

Use this data source to get the list of CAE applications within HuaweiCloud.

## Example Usage

### Query the all applications under the default enterprise project

```hcl
variable "environment_id" {}

data "huaweicloud_cae_applications" "test" {
  environment_id = var.environment_id
}
```

### Query the application under the specified enterprise project by application ID

```hcl
variable "environment_id" {}
variable "application_id" {}
variable "enterprise_project_id" {}

data "huaweicloud_cae_applications" "filter_by_application_id" {
  environment_id        = var.environment_id
  application_id        = var.application_id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the CAE application list.
  If omitted, the provider-level region will be used.

* `environment_id` - (Required, String) Specifies the ID of the environment to which the applications belong.

* `application_id` - (Optional, String) Specifies the ID of the application to be queried.

* `name` - (Optional, String) Specifies the name of the application to be queried.
  The name can contain `2` to `64` characters, only lowercase letters, digits, and hyphens (-) allowed.
  The name must start with a lowercase letter and end with lowercase letters and digits.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the applications
  belong.  
  If the `environment_id` belongs to the non-default enterprise project, this parameter is required and is only valid for
  enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - All applications that match the filter parameters.
  The [applications](#CAE_applications) structure is documented below.

<a name="CAE_applications"></a>
The `applications` block supports:

* `id` - The ID of the application.

* `name` - The name of the application.

* `created_at` - The creation time of the application.

* `updated_at` - The latest update time of the application.
