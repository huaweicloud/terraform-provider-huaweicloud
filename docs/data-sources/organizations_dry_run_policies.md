---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_dry_run_policies"
description: |-
  Use this data source to get the list of Organizations dry run policies within HuaweiCloud.
---

# huaweicloud_organizations_dry_run_policies

Use this data source to get the list of Organizations dry run policies within HuaweiCloud.

## Example Usage

### Query all dry run policies

```hcl
data "huaweicloud_organizations_dry_run_policies" "test" {}
```

### Query dry run policies by attached entity ID

```hcl
variable "attached_entity_id" {}

data "huaweicloud_organizations_dry_run_policies" "test" {
  attached_entity_id = var.attached_entity_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_type` - (Optional, String) Specifies the type of the dry run policies to be queried.  
  The valid values are as follows:
  + **service_control_policy**: Service control policy.

* `attached_entity_id` - (Optional, String) Specifies the ID of the entity associated with the dry run policy.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - The list of dry run policies that matched the filter parameters.  
  The [policies](#dry_run_policies) structure is documented below.

<a name="dry_run_policies"></a>
The `policies` block supports:

* `id` - The unique ID of the dry run policy.

* `name` - The name of the dry run policy.

* `type` - The type of the dry run policy.

* `urn` - The uniform resource name of the dry run policy.

* `description` - The description of the dry run policy.

* `is_builtin` - Whether the dry run policy is a built-in policy.
