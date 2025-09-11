---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_diagnosis_task_node_detail"
description: |-
  Use this data source to get the diagnosis task node detail.
---

# huaweicloud_coc_diagnosis_task_node_detail

Use this data source to get the diagnosis task node detail.

## Example Usage

```hcl
variable "task_id" {}
variable "code" {}
variable "instance_id" {}

data "huaweicloud_coc_diagnosis_task_node_detail" "test" {
  task_id     = var.task_id
  code        = var.code
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Specifies the diagnostic work order ID.

* `code` - (Required, String) Specifies the diagnostic step code.
  Values can be as follows:
  + **holmesInstall**: Installs the Holmes diagnostic plugin.
  + **dataCollection**: Collects data.
  + **diagnosisFault**: Performs fault diagnosis.
  + **holmesUnInstall**: Uninstalls the Holmes diagnostic plugin.
  + **rdsDiagnosis**: Performs diagnostics for the RDS database service.
  + **dcsDiagnosis**: Performs diagnostics for the Distributed Cache Service (DCS).
  + **dmsDiagnosis**: Performs diagnostics for the Distributed Messaging Service (DMS).
  + **elbDiagnosis**: Performs diagnostics for the Elastic Load Balancer (ELB).

* `instance_id` - (Required, String) Specifies the ID of the instance being diagnosed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - Indicates the diagnosis step name.

* `name_zh` - Indicates the Chinese name of the diagnostic step.

* `status` - Indicates the execution status of the diagnostic task.

* `diagnosis_record_id` - Indicates the primary key ID of the diagnosis step.

* `start_time` - Indicates the start time of the diagnostic step.

* `end_time` - Indicates the end time of the diagnostic step.

* `message` - Indicates the diagnostic step execution logging is performed.

* `diagnostic_task_node_id` - Indicates the diagnostic task node ID.
