---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_queueing_records"
description: |-
  Use this data source to get a list of CodeArts pipeline queueing records.
---

# huaweicloud_codearts_pipeline_queueing_records

Use this data source to get a list of CodeArts pipeline queueing records.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}

data "huaweicloud_codearts_pipeline_queueing_records" "test" {
  project_id  = var.codearts_project_id
  pipeline_id = var.pipeline_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `pipeline_id` - (Required, String) Specifies the pipeline ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - Indicates the record list.
  The [records](#attrblock--records) structure is documented below.

<a name="attrblock--records"></a>
The `records` block supports:

* `id` - Indicates the record ID.

* `pipeline_run_id` - Indicates the pipeline run ID.

* `status` - Indicates the status.

* `trigger_type` - Indicates the trigger type.

* `enqueue_time` - Indicates the enqueuing time.

* `creator_id` - Indicates the creator ID.

* `creator_name` - Indicates the creator name.
