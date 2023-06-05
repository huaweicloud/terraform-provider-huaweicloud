---
subcategory: "Organizations"
---

# huaweicloud_organizations

Manages an Organizations resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_organizations" "test"{
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `root_tags` - (Optional, Map) Specifies the key/value pairs to associate with the root.

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

## Import

The organizations can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations.test <id>
```
