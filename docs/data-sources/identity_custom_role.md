---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud\_identity\_custom\_role

Use this data source to get the ID of an HuaweiCloud custom role.

The Role in Terraform is the same as Policy on console. however,
The policy name is the display name of Role, the Role name cannot
be found on Console. 

```hcl
data "huaweicloud_identity_custom_role" "role" {
  name = "custom_role"
}
```

## Argument Reference

* `name` - (Optional, String) Name of the custom policy. 

* `id` - (Optional, String) ID of the custom policy.

* `domain_id` - (Optional, String) The domain the policy belongs to.

* `references` - (Optional, Int) The number of citations for the custom policy.

* `description` - (Optional, String) Description of the custom policy.

* `type` - (Optional, String) Display mode. Valid options are AX: Account level and XA: Project level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy` - Document of the custom policy.

* `catalog` - The catalog of the custom policy.
