---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_migrate"
description: |-
  Using this resource to migrate CBR resources within HuaweiCloud.
---

# huaweicloud_cbr_migrate

Using this resource to migrate CBR resources within HuaweiCloud.

-> This resource is a one-time action resource to migrate CBR resources. Deleting this resource will
not affect the migration result, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_cbr_migrate" "test" {
  all_regions = false
  reservation = 0.5
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `all_regions` - (Required, Bool, NonUpdatable) Specifies whether to trigger migration in other regions.

* `reservation` - (Required, Float, NonUpdatable) Specifies the default expansion ratio of the vault.
  The value must be a float between `0` and `1`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
