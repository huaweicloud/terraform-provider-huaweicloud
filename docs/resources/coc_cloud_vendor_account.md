---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_cloud_vendor_account"
description: |-
  Manages a COC cloud vendor (such as Alibaba Cloud, AWS, Azure, and Huawei Cloud Stack) account resource within HuaweiCloud.
---

# huaweicloud_coc_cloud_vendor_account

Manages a COC cloud vendor (such as Alibaba Cloud, AWS, Azure, and Huawei Cloud Stack) account resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_coc_cloud_vendor_account" "test" {
  vendor       = "HCS"
  account_id   = "test_account_id"
  account_name = "test_account_name"
  ak           = "test_ak"
  sk           = "test_sk"
}
```

## Argument Reference

The following arguments are supported:

* `vendor` - (Required, String, NonUpdatable) Specifies the cloud vendor information.
  Value options::
  + **RMS**: Huawei Cloud.
  + **AZURE**: Microsoft Azure
  + **ALI**: Alibaba Cloud
  + **VMWARE**: VMware
  + **OPENSTACK**: OpenStack cloud platform
  + **HCS**: Huawei hybrid cloud solution Huawei Cloud Stack
  + **OTHER**: other cloud vendors

* `account_id` - (Required, String, NonUpdatable) Specifies the account ID of a supplier. The value is a string. It
  contains 0 to 64 characters.

* `account_name` - (Required, String) Specifies the account name. The value is a string. It contains 0 to 64 characters.

* `ak` - (Required, String) Specifies the AK. The value is a string. It contains 0 to 64 characters.

* `sk` - (Required, String) Specifies the SK. The value is a string. It contains 0 to 64 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `domain_id` - Indicates the tenant ID.

* `sync_status` - Indicates the task status.

* `failure_msg` - Indicates the error message.

* `sync_date` - Indicates the synchronization time.

* `create_time` - Indicates the creation time.

* `update_time` - Indicates the update time.

## Import

The COC cloud vendor account can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_cloud_vendor_account.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `sk`. It is generally recommended running
`terraform plan` after importing a cloud vendor account. You can then decide if changes should be applied to the cloud
vendor account, or the resource definition should be updated to align with the cloud vendor account. Also you can ignore
changes as below.

```hcl
resource "huaweicloud_coc_cloud_vendor_account" "test" {
    ...

  lifecycle {
    ignore_changes = [
      sk
    ]
  }
}
```
