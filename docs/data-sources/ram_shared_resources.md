---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_shared_resources"
description: |
  Use this data source to get the list of RAM shared resources.
---

# huaweicloud_ram_shared_resources

Use this data source to get the list of RAM shared resources.

## Example Usage

```hcl
data "huaweicloud_ram_shared_resources" "test" {
  resource_owner = "self"
}
```

## Argument Reference

The following arguments are supported:

* `resource_owner` - (Required, String) Specifies the owner associated with the RAM share.
  Value options are as follows:
    + **self**: Shared to other users by myself.
    + **other-accounts**: Shared to me by other users.

* `principal` - (Optional, String) Specifies the principal associated with the RAM share.
  The principal could be account ID or organization ID.
  + If set to account ID, please make sure the account ID is not your owner account ID.
  + If set to organization ID, you first need to use the RAM console to enable sharing with Organization. Please refer
  to the [document](https://support.huaweicloud.com/intl/en-us/qs-ram/ram_02_0004.html).

* `resource_ids` - (Optional, List) Specifies the list of resource IDs associated with the RAM share.

* `resource_region` - (Optional, String) Specifies the resource region associated with the RAM share.

* `resource_share_ids` - (Optional, List) Specifies the list of resource share IDs.

* `resource_type` - (Optional, String) Specifies the resource type associated with the RAM share.

* `resource_urns` - (Optional, List) Specifies one or more resources urns associated with the
  RAM share. The format of URN is: `<service-name>:<region>:<account-id>:<type-name>:<resource-path>`.
  Sharable cloud services and resource types refer to
  [document](https://support.huaweicloud.com/intl/en-us/productdesc-ram/ram_01_0007.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `shared_resources` - The list of shared resources.
  The [shared_resources](#attrblock--shared_resources) structure is documented below.

<a name="attrblock--shared_resources"></a>
The `shared_resources` block supports:

* `resource_share_id` - The resource share ID.

* `resource_type` - The resource type associated with the RAM share.

* `resource_urn` - The resource urn associated with the RAM share.

* `status` - The status of the RAM share.

* `created_at` - The creation time of the RAM share.

* `updated_at` - The latest update time of the RAM share.
