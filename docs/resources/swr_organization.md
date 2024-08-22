---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_organization"
description: ""
---

# huaweicloud_swr_organization

Manages a SWR organization resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_swr_organization" "test" {
  name = "terraform-test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the organization. The organization name must be globally
  unique.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the organization.

* `creator` - The creator user name of the organization.

* `permission` - The permission of the organization, the value can be Manage, Write, and Read.

* `login_server` - The URL that can be used to log into the container registry.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

Organizations can be imported using the `name`, e.g.

```bash
$ terraform import huaweicloud_swr_organization.test org-name
```
