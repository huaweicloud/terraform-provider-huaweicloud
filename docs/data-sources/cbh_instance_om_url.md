---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_instance_om_url"
description: |-
  Use this data source to get the OM URL for a managed host in a CBH instance within HuaweiCloud.
---

# huaweicloud_cbh_instance_om_url

Use this data source to get the OM URL for a managed host in a CBH instance within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}
variable "ip_address" {}
variable "host_account_name" {}

data "huaweicloud_cbh_instance_om_url" "test" {
  server_id         = var.server_id
  ip_address        = var.ip_address
  host_account_name = var.host_account_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the ID of the CBH instance.

* `ip_address` - (Required, String) Specifies the IP address of the managed host.

* `host_account_name` - (Required, String) Specifies the account name of the managed host.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `login_url` - The OM login URL for the managed host.

* `state` - The state of getting the OM URL. Possible values include:
  + **SUCCESS**: Successfully retrieved the OM URL. If the status is not SUCCESS, the task failed. Locate the cause by
    referring to the following error codes and the corresponding descriptions.
  + **IAM_USER_CONFLICT**: IAM user conflict.
  + **HOST_NOT_MANAGE**: The host is not managed.
  + **HOST_ACCOUNT_NOT_EXIST**: The host account does not exist.
  + **IAM_USER_NO_PERMISSION**: IAM user has no permission.
  + **CUR_VERSION_NOT_SUPPORT_OPERATION**: Current version does not support this operation.
  + **INS_WHITE_LIST_INITIALIZING**: Instance whitelist is initializing.
  + **UNKNOWN_ERROR**: Unknown error occurred.

* `description` - The description when failed to get the OM URL.
