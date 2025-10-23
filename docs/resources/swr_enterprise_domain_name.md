---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_domain_name"
description: |-
  Manages a SWR enterprise instance domain name resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_domain_name

Manages a SWR enterprise instance domain name resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "domain_name" {}
variable "certificate_id" {}

resource "huaweicloud_swr_enterprise_domain_name" "test" {
  instance_id    = var.instance_id
  domain_name    = var.domain_name
  certificate_id = var.certificate_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `domain_name` - (Required, String, NonUpdatable) Specifies the domain name.

* `certificate_id` - (Required, String) Specifies the SCM certificate ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `domain_name_id` - Indicates the domain name ID.

* `type` - Indicates the domain name type.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

## Import

The domain name can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_domain_name.test <id>
```
