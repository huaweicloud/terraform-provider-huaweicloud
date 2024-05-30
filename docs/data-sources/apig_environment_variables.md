---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_environment_variables"
description: |-
  Use this data source to get the list of environment variables under specified group of the APIG instance within HuaweiCloud.
---

# huaweicloud_apig_environment_variables

Use this data source to get the list of environment variables under specified group of the APIG instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_id" {}

resource "huaweicloud_apig_environment_variables" "test" {
  instance_id = var.instance_id
  group_id    = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the environment variables belong.

* `group_id` - (Required, String) Specifies the ID of the group to which the environment variables belong.

* `env_id` - (Optional, String) Specifies the ID of the environment to which the environment variables belong.

* `name` - (Optional, String) Specifies the name of the environment variable.
  Fuzzy search is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `variables` - All environment variables that match the filter parameters.

  The [variables](#variables_struct) structure is documented below.

<a name="variables_struct"></a>
The `variables` block supports:

* `id` - The ID of the environment variable.

* `group_id` - The group ID corresponding to the environment variable.

* `env_id` - The environment ID corresponding to the environment variable.

* `name` - The name of the environment variable.

* `value` - The value of the environment variable.
