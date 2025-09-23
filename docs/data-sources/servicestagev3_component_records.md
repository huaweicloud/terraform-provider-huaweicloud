---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_component_records"
description: |-
  Use this data source to query the list of component execution records within HuaweiCloud.
---

# huaweicloud_servicestagev3_component_records

Use this data source to query the list of component execution records within HuaweiCloud.

## Example Usage

```hcl
variable "application_id" {}
variable "component_id" {}

data "huaweicloud_servicestagev3_component_records" "test" {
  application_id = var.application_id
  component_id   = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the application and component are located.

* `application_id` - (Required, String) Specifies the ID of the application to which the component belongs.

* `component_id` - (Required, String) Specifies the ID of the component to which the records belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `records` - The list of component execution record.  
  The [records](#servicestage_v3_component_records) structure is documented below.

<a name="servicestage_v3_component_records"></a>
The `applications` block supports:

* `begin_time` - The begin time of the component execution record.

* `end_time` - The end time of the component execution record.

* `description` - The description of the component execution record.

* `instance_id` - The instance ID of the component execution record.

* `version` - The version number of the component execution record.

* `current_used` - Whether version is current used.

* `status` - The status of the component execution record.

* `deploy_type` - The deploy type of the component execution record.

* `jobs` - The list of component jobs.  
  The [jobs](#servicestage_v3_component_record_jobs) structure is documented below.

<a name="servicestage_v3_component_record_jobs"></a>
The `jobs` block supports:

* `sequence` - The sequence of the job execution.

* `job_id` - The job ID.

* `job_info` - The job detail.
  The [job_info](#servicestage_v3_component_record_job_info) structure is documented below.

<a name="servicestage_v3_component_record_job_info"></a>
The `job_info` block supports:

* `source_url` - The source URL of the component.

* `first_batch_weight` - The weight of the first batch execution.`

* `first_batch_replica` - The replica of the first batch execution.

* `replica` - The total replica number.

* `remaining_batch` - The remaining batch number.
