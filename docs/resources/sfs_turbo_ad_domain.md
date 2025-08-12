---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_ad_domain"
description: |-
  Use this resource to manage the AD domain of the SFS turbo within HuaweiCloud.
---

# huaweicloud_sfs_turbo_ad_domain

Use this resource to manage the AD domain of the SFS turbo within HuaweiCloud.

## Example Usage

```hcl
variable "share_id" {}
variable "service_account" {}
variable "password" {}
variable "domain_name" {}
variable "system_name" {}
variable "dns_server" {
  type = list(string)
}
variable "organization_unit" {}
variable "vpc_id" {}

resource "huaweicloud_sfs_turbo_ad_domain" "test" {
  share_id               = var.share_id
  service_account        = var.service_account
  password               = var.password
  domain_name            = var.domain_name
  system_name            = var.system_name
  dns_server             = var.dns_server
  overwrite_same_account = false
  organization_unit      = var.organization_unit
  vpc_id                 = var.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `share_id` - (Required, String, NonUpdatable) Specifies the ID of the SFS Turbo.

* `service_account` - (Required, String) Specifies the service account, which is specified when the domain
  server is created, **administrator** is used normally.

* `password` - (Required, String) Specifies the password of the service account.

* `domain_name` - (Required, String) Specifies the domain name of the domain controller. It is specified when the
  domain server is created.

* `system_name` - (Required, String) Specifies the name of the file storage system in the AD domain.

* `dns_server` - (Required, List) Specifies the IP address of the DNS server. It is used to resolve the AD domain
  name.

* `overwrite_same_account` - (Optional, Bool) Whether overwrite the existing information in the domain controller.
  If the option is enabled and the domain controller already has the file system name you specified, the information
  you specified will be overwrited.

* `organization_unit` - (Optional, String) Specifies the  group of domain objects, such as users, computers,
  and printers. If you add the file system to an organizational unit (OU), it will become a member of that OU.
  If this parameter is left blank, the file system will be added to the computers OU.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The current status of the AD domain. Possible values are: **JOINING**, **AVAILABLE**, **EXITING**
  and **FAILED**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The SFS Turbo AD domain can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_sfs_turbo_ad_domain.test <id>
```

```
Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `share_id`, `service_account`, `password`,
`overwrite_same_account`. It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the snapshot group. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_sfs_turbo_ad_domain" "test" {
    ...
  lifecycle {
    ignore_changes = [
      share_id,
      service_account,
      password,
      overwrite_same_account,
    ]
  }
}
```
