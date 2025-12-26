---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_policies"
description: |-
  Use this data source to get a list of IAM V5 policies.
---


# huaweicloud_identityv5_policies

Use this data source to get a list of IAM V5 policies.

## Example Usage

```hcl
data "huaweicloud_identityv5_policies" "all" {}
```

## Argument Reference

The following arguments are supported:

* `policy_type` - (Optional, String) The type of identity policy, can be "custom" or "system".

* `path_prefix` - (Optional, String) The resource path prefix, composed of segments of strings,
  each segment contains one or more letters, digits, ".", ",", "+", "@", "=", "_", or "-",
  ending with "/", for example "foo/bar/".

* `only_attached` - (Optional, Bool) Whether to list only policies that are attached to entities.

* `language` - (Optional, String) The language.

* `policy_id` - (Optional, String) The policy ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policies` - The list of IAM policies.
  The [policies](#IdentityPolicies_List) structure is documented below.

<a name="IdentityPolicies_List"></a>
The `policies` block supports:

* `attachment_count` - Indicates the number of entities attached to the policy.

* `default_version_id` - Indicates the ID of the default version of the policy.

* `path` - Indicates the path of the policy.

* `policy_id` - Indicates the ID of the policy.

* `policy_name` - Indicates the name of the policy.

* `policy_type` - Indicates the type of the policy.

* `updated_at` - Indicates the time when the default version of the policy was last updated.

* `urn` - Indicates the URN of the policy.

* `created_at` - Indicates the time when the policy was created.

* `description` - Indicates the description of the policy.
