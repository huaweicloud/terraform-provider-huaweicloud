---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_custom_role"
description: ""
---

# huaweicloud_identity_custom_role

Use this data source to get details of the specified IAM **custom policy**.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

```hcl
data "huaweicloud_identity_custom_role" "policy" {
  name = "custom_policy"
}
```

## Argument Reference

* `name` - (Optional, String) Specifies the name of the custom policy. It's required if `id` is not specified.

* `id` - (Optional, String) Specifies the ID of the custom policy. It's required if `name` is not specified.

* `description` - (Optional, String) Specifies the description of the custom policy.

* `type` - (Optional, String) Specifies the display mode of the custom policy. Valid options are as follows:
  + **AX**: the global service project.
  + **XA**: region-specific projects.

* `domain_id` - (Optional, String) Specifies the domain the policy belongs to.

* `references` - (Optional, Int) Specifies the number of citations for the custom policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `policy` - The content of the custom policy in JSON format.

* `catalog` - The catalog of the custom policy. The value is **CUSTOMED**.
