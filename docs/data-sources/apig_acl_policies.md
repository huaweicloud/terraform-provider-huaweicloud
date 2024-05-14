---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_acl_policies"
description: |-
  Use this data source to query the ACL policies within HuaweiCloud.
---

# huaweicloud_apig_acl_policies

Use this data source to query the ACL policies within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "policy_name" {}

data "huaweicloud_apig_acl_policies" "test" {
  instance_id = var.instance_id
  name        = var.policy_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the ACL policies belong.

* `policy_id` - (Optional, String) Specifies the ID of the ACL policy to be queried.

* `name` - (Optional, String) Specifies the name of the ACL policy to be queried.

* `type` - (Optional, String) Specifies the type of the ACL policy to be queried.  
  The valid values are as follows:
  + **PERMIT**: The whitelist type Strategies.
  + **DENY**: The blacklist type Strategies.

* `entity_type` - (Optional, String) Specifies the entity type of the ACL policy to be queried.  
  The valid values are as follows:
  + **IP**
  + **DOMAIN**
  + **DOMAIN_ID**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - All ACL policies that match the filter parameters.
  The [policies](#acl_policies) structure is documented below.

<a name="acl_policies"></a>
The `policies` block supports:

* `id` - The ID of the ACL policy.

* `name` - The name of the ACL policy.

* `type` - The type of the ACL policy.

* `value` - The value of the ACL policy.

* `bind_num` - The number of bound APIs.

* `entity_type` - The entity type of the ACL policy.

* `updated_at` - The latest update time of the policy.
