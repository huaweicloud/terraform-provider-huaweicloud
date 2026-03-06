---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_policies"
description: |-
  Use this data source to get the list of Organizations policies within HuaweiCloud.
---

# huaweicloud_organizations_policies

Use this data source to get the list of Organizations policies within HuaweiCloud.

## Example Usage

### Query all policies

```hcl
data "huaweicloud_organizations_policies" "test" {}
```

### Query policies by build type

```hcl
data "huaweicloud_organizations_policies" "test" {
  build_type = "system"
}
```

## Argument Reference

The following arguments are supported:

* `build_type` - (Optional, String) Specifies the build type of the policy.  
  The valid values are as follows:
  + **system**: System policy.
  + **custom**: Custom policy.

* `name` - (Optional, String) Specifies the name of the policy.

* `type` - (Optional, String) Specifies the type of the policy.  
  The valid values are as follows:
  + **service_control_policy**: Service control policy.
  + **tag_policy**: Tag policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - The list of policies that match the filter parameters.  
  The [policies](#Policies_Policy) structure is documented below.

<a name="Policies_Policy"></a>
The `policies` block supports:

* `id` - Indicates the unique ID of the policy.

* `name` - Indicates the name of the policy.

* `type` - Indicates the type of the policy.

* `urn` - Indicates the uniform resource name of the policy.

* `description` - Indicates the description of the policy.

* `build_type` - Indicates the build type of the policy.
