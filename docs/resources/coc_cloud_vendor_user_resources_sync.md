---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_change_update"
description: |-
  Manages a COC cloud vendor user resources sync resource within HuaweiCloud.
---

# huaweicloud_coc_change_update

Manages a COC cloud vendor user resources sync resource within HuaweiCloud.

~> Deleting cloud vendor user resources sync resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
resource "huaweicloud_coc_cloud_vendor_user_resources_sync" "test" {
  vendor = "HCS"
}
```

## Argument Reference

The following arguments are supported:

* `vendor` - (Required, String, NonUpdatable) Specifies the cloud vendor information.
  Values options:
  + **RMS**: Huawei Cloud
  + **AZURE**: Microsoft Azure
  + **ALI**: Alibaba Cloud
  + **VMWARE**: VMware
  + **OPENSTACK**: OpenStack cloud platform
  + **HCS**: Huawei hybrid cloud solution Huawei Cloud Stack
  + **OTHER**: other cloud vendors

* `account_id` - (Optional, String, NonUpdatable) Specifies the account ID of a supplier.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
