---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_associated_acl_policies"
description: |-
  Use this data source to query the ACL policies associated with the specified API within HuaweiCloud.
---

# huaweicloud_apig_api_associated_acl_policies

Use this data source to query the ACL policies associated with the specified API within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "associated_api_id" {}

data "huaweicloud_apig_api_associated_acl_policies" "test" {
  instance_id = var.instance_id
  api_id      = var.associated_api_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the associated ACL policies.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the ACL policies belong.

* `api_id` - (Required, String) Specifies the ID of the API bound to the ACL policy.

* `policy_id` - (Optional, String) Specifies the ID of the ACL policy.

* `name` - (Optional, String) Specifies the name of the ACL policy.

* `type` - (Optional, String) Specifies the type of the ACL policy.

* `env_id` - (Optional, String) Specifies the ID of the environment where the API is published.

* `env_name` - (Optional, String) Specifies the name of the environment where the API is published.

* `entity_type` - (Optional, String) Specifies the entity type of the ACL policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - All ACL policies that match the filter parameters.
  The [policies](#api_associated_acl_policies) structure is documented below.

<a name="api_associated_acl_policies"></a>
The `policies` block supports:

* `id` - The ID of the ACL policy.

* `name` - The name of the ACL policy.

* `type` - The type of the ACL policy.

* `value` - One or more objects from which the access will be controlled.

* `env_id` - The ID of the environment where the API is published.

* `env_name` - The name of the environment where the API is published.

* `entity_type` - The entity type of the ACL policy.

* `bind_id` - The bind ID.

* `bind_time` - The time that the ACL policy is bound to the API.
