---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_lts_config"
description: |-
  Manages a resource to config LTS to DRS within HuaweiCloud.
---

# huaweicloud_drs_lts_config

Manages a resource to config LTS to DRS within HuaweiCloud.

-> Deleting this resource will disable the LTS switch.

## Example Usage

```hcl
variable "job_id" {}
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_drs_lts_config" "test" { 
  job_id        = var.job_id
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `log_group_id` - (Optional, String) Specifies the log group ID.

* `log_stream_id` - (Optional, String) Specifies the log stream ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `job_id`).

## Import

The DRS LTS config can be imported by `id`. e.g.

```bash
$ terraform import huaweicloud_drs_lts_config.test <id>
```
