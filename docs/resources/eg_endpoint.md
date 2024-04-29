---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_endpoint"
description: ""
---

# huaweicloud_eg_endpoint

Manages an EG endpoint resource within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_eg_endpoint" "test" {
  name        = "test"
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  description = "created by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the endpoint.
  The value can contain no more than 128 characters, including letters, digits, underscores (_), hyphens (-),
  and periods (.), and must start with a character or letter.

  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC to which the endpoint belongs.

  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the subnet to which the endpoint belongs.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the endpoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `domain` - The domain of the endpoint.

* `status` - The status of the endpoint.

* `created_at` - The creation time of the endpoint.

* `updated_at` - The last update time of the endpoint.

## Import

The endpoint can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_eg_endpoint.test 32a6a33f-ac15-4548-a328-8dc91ed22c3c
```
