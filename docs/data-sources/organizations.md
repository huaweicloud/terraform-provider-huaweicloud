---
subcategory: "Organizations"
---

# huaweicloud_organizations

Use this data source to get the Organization info and the root info.

## Example Usage

```hcl
data "huaweicloud_organizations" "test"{
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the uniform resource name of the organization.

* `account_id` - Indicates the unique ID of the organization's management account.

* `account_name` - Indicates the name of the organization's management account.

* `created_at` - Indicates the time when the organization was created.

* `root_id` - Indicates the ID of the root.

* `root_name` - Indicates the name of the root.

* `root_urn` - Indicates the urn of the root.

* `root_tags` - Indicates the key/value to attach to the root.
