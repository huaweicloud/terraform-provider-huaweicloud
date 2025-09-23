---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_custom_authorizer"
description: ""
---

# huaweicloud_apig_custom_authorizer

Manages an APIG custom authorizer resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "authorizer_name" {}
variable "function_urn" {}

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id      = var.instance_id
  name             = var.authorizer_name
  function_urn     = var.function_urn
  function_version = "latest"
  type             = "FRONTEND"
  cache_age        = 60

  identity {
    name     = "user_name"
    location = "QUERY"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the custom authorizer resource.
  If omitted, the provider-level region will be used.
  Changing this will create a new custom authorizer resource.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the
  custom authorizer belongs to.
  Changing this will create a new custom authorizer resource.

* `name` - (Required, String) Specifies the name of the custom authorizer.
  The custom authorizer name consists of `3` to `64` characters, starting with a letter.
  Only letters, digits and underscores (_) are allowed.

* `function_urn` - (Required, String) Specifies the uniform function URN of the function graph resource.

* `function_version` - (Required, String) Specifies the version of the FGS function.

* `network_type` - (Optional, String) Specifies the framework type of the function.
  + **V1**: Non-VPC network architecture.
  + **V2**: VPC network architecture.

  Defaults to **V1**.

* `function_alias_uri` - (Optional, String) Specifies the version alias URI of the FGS function.

* `type` - (Optional, String, ForceNew) Specifies the custom authorize type.
  The valid values are **FRONTEND** and **BACKEND**. Defaults to **FRONTEND**.
  Changing this will create a new custom authorizer resource.

* `is_body_send` - (Optional, Bool) Specifies whether to send the body.

* `cache_age` - (Optional, Int) Specifies the maximum cache age.  
  The valid value is range from `1` to `3,600`.

* `user_data` - (Optional, String) Specifies the user data, which can contain a maximum of `2,048` characters.
  The user data is used by APIG to invoke the specified authentication function when accessing the backend service.

  -> **NOTE:** The user data will be displayed in plain text on the console.

* `identity` - (Optional, List) Specifies an array of one or more parameter identities of the custom authorizer.
  The [object](#authorizer_identity) structure is documented below.

<a name="authorizer_identity"></a>
The `identity` block supports:

* `name` - (Required, String) Specifies the name of the parameter to be verified.
  The parameter includes front-end and back-end parameters.

* `location` - (Required, String) Specifies the parameter location, which support **HEADER** and **QUERY**.

* `validation` - (Optional, String) Specifies the parameter verification expression.
  If omitted, the custom authorizer will not perform verification.
  The valid value is range form `1` to `2,048`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the custom authorizer.

* `created_at` - The creation time of the custom authorizer.

## Import

Custom Authorizers of the APIG can be imported using their `name` and related dedicated instance IDs, separated by a
slash, e.g.

```shell
$ terraform import huaweicloud_apig_custom_authorizer.test <instance_id>/<name>
```
