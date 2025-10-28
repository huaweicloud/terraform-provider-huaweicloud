---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_cloud_vendor_accounts"
description: |-
  Use this data source to get the list of COC cloud vendor account.
---

# huaweicloud_coc_cloud_vendor_accounts

Use this data source to get the list of COC cloud vendor account.

## Example Usage

```hcl
data "huaweicloud_coc_cloud_vendor_accounts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `vendor` - (Optional, String) Specifies the cloud vendor information.
  Value options::
  + **RMS**: Huawei Cloud.
  + **AZURE**: Microsoft Azure
  + **ALI**: Alibaba Cloud
  + **VMWARE**: VMware
  + **OPENSTACK**: OpenStack cloud platform
  + **HCS**: Huawei hybrid cloud solution Huawei Cloud Stack
  + **OTHER**: other cloud vendors

* `account_id` - (Optional, String) Specifies the account ID of a supplier.

* `account_name` - (Optional, String) Specifies the account name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of cloud vendor account.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the cloud vendor account ID allocated by CloudCMDB

* `vendor` - Indicates the cloud vendor information.

* `account_id` - Indicates the account ID of a supplier.

* `account_name` - Indicates the account name.

* `ak` - Indicates the AK.

* `domain_id` - Indicates the tenant ID.

* `sync_status` - Indicates the task status.

* `failure_msg` - Indicates the error message.

* `sync_date` - Indicates the synchronization time.

* `create_time` - Indicates the creation time.

* `update_time` - Indicates the update time.
