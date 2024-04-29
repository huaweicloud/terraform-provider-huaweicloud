---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_policies"
description: ""
---

# huaweicloud_organizations_policies

Use this data source to get the list of policies in an organization.

## Example Usage

```hcl
data "huaweicloud_organizations_policies" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `build_type` - (Optional, String) Specifies the build type of the policy.
  + **system**: system policy.
  + **custom**: custom policy.

* `name` - (Optional, String) Specifies the name of the policy.

* `type` - (Optional, String) Specifies the type of the policy. Value options:
  + **service_control_policy**: service control policy.
  + **tag_policy**: tag policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - List of policies in an organization.
  The [policies](#Policies_Policy) structure is documented below.

<a name="Policies_Policy"></a>
The `policies` block supports:

* `id` - Indicates the unique ID of the policy.

* `name` - Indicates the name of the policy.

* `type` - Indicates the type of the policy.

* `urn` - Indicates the uniform resource name of the policy.

* `description` - Specifies the description of the policy.

* `build_type` - Indicates the build type of the policy.
