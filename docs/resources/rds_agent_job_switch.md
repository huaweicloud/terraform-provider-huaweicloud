---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_agent_job_switch"
description: |-
  Manages an RDS agent job switch resource within HuaweiCloud.
---

# huaweicloud_rds_agent_job_switch

Manages an RDS agent job switch resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "job_id" {}

resource "huaweicloud_rds_agent_job_switch" "test" {
  instance_id = var.instance_id
  job_id      = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS instance.

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
