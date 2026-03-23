---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_agency"
description: |-
  Manages an agency resource within HuaweiCloud.
---

# huaweicloud_css_agency

Manages an agency resource within HuaweiCloud.

-> If the built-in CSS agency does not exist, use this resource, the system automatically creates an agency and
  grants the permissions required by CSS to it.

## Example Usage

```hcl
variable "domain_id" {}
variable "domain_name" {}
variable "type" {}

resource "huaweicloud_css_agency" "test" {
  domain_id   = var.domain_id
  domain_name = var.domain_name
  type        = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the cluster resource. If omitted, the
  provider-level region will be used. Changing this creates a new cluster resource.

* `domain_id` - (Required, String, NonUpdatable) Specifies the account ID.

* `domain_name` - (Required, String, NonUpdatable) Specifies the account name.

* `type` - (Required, String, NonUpdatable) Specifies the type of the agency.
  The valid values are as follows:
  + **obs**: Indicates agency permissions required for creating snapshots and log backups.
  + **vpc**: Indicates agency permissions required for version upgrade, AZ change, scale-in, and node replacement.
  + **elb**: Indicates agency permissions required for using the alerting plug-in.
  + **smn**: Indicates agency permissions required for using load balancing.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The agency can be imported by `id`, e.g.

```bash
$ terraform import huaweicloud_css_agency.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `domain_name`, `type`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_css_agency" "test" {
  ...

  lifecycle {
    ignore_changes = [
      domain_name, type,
    ]
  }
}
```
