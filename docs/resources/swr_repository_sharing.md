---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_repository_sharing"
description: ""
---

# huaweicloud_swr_repository_sharing

Manages a SWR repository sharing resource within HuaweiCloud.

## Example Usage

```hcl
variable "organization_name" {} 
variable "repository_name" {}
variable "sharing_account" {}

resource "huaweicloud_swr_repository_sharing" "test" {
  organization    = var.organization_name
  repository      = var.repository_name
  sharing_account = var.sharing_account
  permission      = "pull"
  deadline        = "forever"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `organization` - (Required, String, ForceNew) Specifies the name of the organization (namespace) the repository belongs.
  Changing this creates a new resource.

* `repository` - (Required, String, ForceNew) Specifies the name of the repository to be shared.
  Changing this creates a new resource.

* `sharing_account` - (Required, String, ForceNew) Specifies the name of the account for repository sharing.
  Changing this creates a new resource
  -> **NOTE:** `sharing_account` should be an existing HuaweiCloud account.

* `deadline` - (Required, String) Specifies the end date of image sharing (UTC time in YYYY-MM-DD format,
  for example `2021-10-01`). When the value is set to forever, the image will be permanently available for the domain.
  The validity period is calculated by day. The shared images expire at 00:00:00 on the day after the end date.

* `permission` - (Optional, String) Specifies the permission to be granted. Currently, only the **pull** permission is supported.
  Default value is **pull**.

* `description` - (Optional, String) Specifies the description of the repository sharing.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the repository sharing. The value is the value of `sharing_account`.

* `status` - Indicates the repository sharing is valid (true) or expired (false).

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 2 minutes.
* `update` - Default is 2 minutes.
* `delete` - Default is 2 minutes.

## Import

Repository sharing can be imported using the organization name, repository name and sharing account
separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_swr_repository_sharing.test org-name/repo-name/sharing-account
```
