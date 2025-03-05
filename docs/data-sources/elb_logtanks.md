---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_logtanks"
description: |-
  Use this data source to get the list of ELB logtanks.
---

# huaweicloud_elb_logtanks

Use this data source to get the list of ELB logtanks.

## Example Usage

```hcl
variable "logtank_id" {}

data "huaweicloud_elb_logtanks" "test" {
  logtank_id = var.logtank_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source. If omitted, the provider-level
  region will be used.

* `logtank_id` - (Optional, String) Specifies the ID of the log tank.

* `loadbalancer_id` - (Optional, String) Specifies the ID of a load balancer

* `log_group_id` - (Optional, String) Specifies the log group ID.

* `log_topic_id` - (Optional, String) Specifies the log topic ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logtanks` - Lists the logtanks.
  The [logtanks](#Elb_logtanks) structure is documented below.

<a name="Elb_logtanks"></a>
The `logtanks` block supports:

* `id` - The log ID.

* `loadbalancer_id` - The ID of a load balancer.

* `log_group_id` - The log group ID.

* `log_topic_id` - The log topic ID.
