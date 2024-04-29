---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_environments"
description: ""
---

# huaweicloud_cae_environments

Use this data source to get the list of CAE environments.

## Example Usage

```hcl
variable "environment_name" {}

data "huaweicloud_cae_environments" "test" {
  name = var.environment_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `environment_id` - (Optional, String) Specifies the ID of the environment to be queried.

* `name` - (Optional, String) Specifies the name of the environment to be queried.

* `status` - (Optional, String) Specifies the status of the environment to be queried.
  The valid values are **finish**, **freeze** and **police_freeze**.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the environments belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `environments` - All environments that match the filter parameters.
  The [environments](#CAE_environments) structure is documented below.

<a name="CAE_environments"></a>
The `environments` block supports:

* `id` - The ID of the environment.

* `name` - The name of the environment.

* `status` - The status of the environment.

* `annotations` - The additional attributes of the environment.

* `created_at` - The creation time of the environment.

* `updated_at` - The latest update time of the environment.
