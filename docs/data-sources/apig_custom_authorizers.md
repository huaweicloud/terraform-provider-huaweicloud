---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_custom_authorizers"
description: |-
  Use this data source to query the custom authorizers within HuaweiCloud.
---

# huaweicloud_apig_custom_authorizers

Use this data source to query the custom authorizers within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "authorizer_name" {}

data "huaweicloud_apig_custom_authorizers" "test" {
  instance_id = var.instance_id
  name        = var.authorizer_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the custom authorizers belong.

* `authorizer_id` - (Optional, String) Specifies the ID of the custom authorizer.

* `name` - (Optional, String) Specifies the name of the custom authorizer.  
  The custom authorizer name consists of `3` to `64` characters, starting with a letter.  
  Only letters, digits and underscores (_) are allowed.

* `type` - (Optional, String) Specifies the type of the custom authorizer.  
  The valid values are as follows:
  + **FRONTEND**
  + **BACKEND**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `authorizers` - All custom authorizers that match the filter parameters.
  The [authorizers](#authorizers) structure is documented below.

<a name="authorizers"></a>
The `authorizers` block supports:

* `id` - The ID of the custom authorizer.

* `name` - The name of the custom authorizer.

* `type` - The type of the custom authorizer.

* `function_type` - The type of the FGS function.

* `function_urn` - The URN of the FGS function.

* `network_type` - The network architecture types of function.

* `function_version` - The version of the FGS function.

* `function_alias_uri` - The version alias URI of the FGS function.

* `cache_age` - The maximum cache age of custom authorizer.

* `created_at` - The creation time of custom authorizer.

* `is_body_send` - Whether to send the body of custom authorizer.

* `user_data` - The user data of custom authorizer.

* `identity` - The parameter identities of the custom authorizer.
  The [identity](#authorizers_identity) structure is documented below.

<a name="authorizers_identity"></a>
The `identity` block supports:

* `name` - The name of the parameter to be verified.

* `location` - The parameter location of identity.

* `validation` - The parameter verification expression of identity.
