---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_policy_versions"
description: |-
  Use this data source to get a list of IAM V5 policy versions.
---

# huaweicloud_identityv5_policy_versions

Use this data source to get a list of IAM V5 policy versions.

## Example Usage

```hcl
variable "policy_id" {}

data "huaweicloud_identityv5_policy_versions" "all" {}

data "huaweicloud_identityv5_policy_versions" "test" {
  policy_id  = var.policy_id
  version_id = "v1"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String) Specifies the ID of the policy.

* `version_id` - (Optional, String) Specifies the version ID of the policy. If this parameter is omitted,
  all versions of the policy will be retrieved.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `versions` - The list of IAM policy versions.
  The [versions](#Identityv5_policy_versions) structure is documented below.

<a name="Identityv5_policy_versions"></a>
The `versions` block supports:

* `version_id` - Indicates the version ID of the policy.

* `is_default` - Indicates whether the version is the default version.

* `created_at` - Indicates the time when the policy version was created.

* `document` - Indicates the document of the policy version.
