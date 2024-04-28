---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_environments"
description: ""
---

# huaweicloud_apig_environments

Use this data source to query the environment list under the APIG instance within Huaweicloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "environment_name" {}

data "huaweicloud_apig_environments" "test" {
  instance_id = var.instance_id
  name        = var.environment_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the APIG environment list.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies an ID of the APIG dedicated instance to which the API
  environment belongs.

* `name` - (Optional, String) Specifies the name of the API environment. The API environment name consists of 3 to 64
  characters, starting with a letter. Only letters, digits and underscores (_) are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID.

* `environments` - List of APIG environment details. The structure is documented below.

The `environments` block supports:

* `id` - ID of the APIG environment.

* `name` - The environment name.

* `description` - The description about the API environment.

* `create_time` - Time when the APIG environment was created, in RFC-3339 format.
