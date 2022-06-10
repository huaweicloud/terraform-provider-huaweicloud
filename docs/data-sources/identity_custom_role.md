---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud_identity_custom_role

Use this data source to get the ID of an IAM **custom policy**.

~> You *must* have IAM read privileges to use this data source.

## Example Usage

```hcl
data "huaweicloud_identity_custom_role" "role" {
  name = "custom_role"
}
```

## Argument Reference

* `name` - (Optional, String) Name of the custom policy.

* `id` - (Optional, String) ID of the custom policy.

* `domain_id` - (Optional, String) The domain the policy belongs to.

* `description` - (Optional, String) Description of the custom policy.

* `type` - (Optional, String) Display mode. Valid options are *AX*: Account level and *XA*: Project level.

* `references` - (Optional, Int) The number of citations for the custom policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy` - Document of the custom policy.

* `catalog` - The catalog of the custom policy.
