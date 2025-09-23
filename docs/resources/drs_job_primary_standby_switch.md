---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_primary_standby_switch"
description: |-
  Manages a DRS disaster recovery job primary standby switch resource within HuaweiCloud.
---

# huaweicloud_drs_job_primary_standby_switch

Manages a DRS disaster recovery job primary standby switch resource within HuaweiCloud.

-> Only applies to disaster recovery tasks in progress or disaster recovery failed status. After the primary standby
switchover, DRS job direction will change. Please check.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_job_primary_standby_switch" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `job_id` - (Required, String, ForceNew) Specifies the DRS disaster recovery job ID.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
