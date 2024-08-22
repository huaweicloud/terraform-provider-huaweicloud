---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_logtank"
description: ""
---

# huaweicloud_elb_logtank

Manage an ELB logtank resource within HuaweiCloud.

## Example Usage

```hcl
variable "loadbalancer_id" {}
variable "group_id" {}
variable "topic_id" {}

resource "huaweicloud_elb_logtank" "test" {
  loadbalancer_id = var.loadbalancer_id
  log_group_id    = var.group_id
  log_topic_id    = var.topic_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the logtank resource.
  If omitted, the provider-level region will be used. Changing this creates a new logtank.

* `loadbalancer_id` - (Required, String, ForceNew) Specifies the ID of a loadbalancer. Changing this
  creates a new logtank

* `log_group_id` - (Required, String) Specifies the ID of a log group. It is provided by other service.

* `log_topic_id` - (Required, String) Specifies the ID of the subscribe topic.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The logtank ID.

## Import

ELB logtank can be imported using the logtank ID, e.g.

```bash
$ terraform import huaweicloud_elb_logtank.test 2f148a75-acd3-4ce7-8f63-d5c9fadab3a0
```
