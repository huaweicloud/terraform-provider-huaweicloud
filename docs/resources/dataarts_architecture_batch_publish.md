---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_batch_publish"
description: |-
  Manages a DataArts Architecture batch publish resource within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_batch_publish

Manages a DataArts Architecture batch publish resource within HuaweiCloud.

-> 1. Repeated publishing is not supported when the objects has been published or pending approval status.
   <br>2. This resource is only a one-time action resource for publshing objects. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "approver_user_id" {}
variable "approver_user_name" {}
variable "batch_publish_objects" {
  type = list(object({
    object_id   = string
    object_type = string
  }))
}

resource "huaweicloud_dataarts_architecture_batch_publish" "test" {
  workspace_id       = var.workspace_id
  approver_user_id   = var.approver_user_id
  approver_user_name = var.approver_user_name
  fast_approval      = true

  dynamic "biz_infos" {
    for_each = var.batch_publish_objects
    content {
      biz_id   = biz_infos.value["object_id"]
      biz_type = biz_infos.value["object_type"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, ForceNew) The ID of DataArts Studio workspace.
  Changing this creates a new resource.

* `biz_infos` - (Required, List, ForceNew) Specifies the list of objects to be published.
  Changing this creates a new resource.
  The [biz_infos](#publishment_biz_infos) structure is documented below.

  -> If the parameter contains published objects, the resource creation will fail, but the remaining objects
     will continue to be published.

* `approver_user_id` - (Required, String, ForceNew) Specifies the user ID of the architecture reviewer.
  Changing this creates a new resource.

* `approver_user_name` - (Required, String, ForceNew) Specifies the user name of the architecture reviewer.
  Changing this creates a new resource.
  
* `fast_approval` - (Optional, Bool, ForceNew) Specifies whether to automatically review.
  Changing this creates a new resource. The default value is `false`.

  -> This parameter is available only when the current user has approval permission.

* `schedule_time` - (Optional, String, ForceNew) Specifies scheduling time of the DataArts quality job.
  The format is `mm_hh`. e.g. `30_18`, it means `18:30`.
  Changing this creates a new resource.

<a name="publishment_biz_infos"></a>
The `biz_infos` block supports:

* `biz_id` - (Required, String, ForceNew) Specifies the ID of the object to be published.
  Changing this creates a new resource.

* `biz_type` - (Required, String, ForceNew) Specifies the type of the object to be published.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **AGGREGATION_LOGIC_TABLE**
  + **ATOMIC_INDEX**
  + **ATOMIC_METRIC**
  + **BIZ_CATALOG**
  + **BIZ_METRIC**
  + **CODE_TABLE**
  + **COMMON_CONDITION**
  + **COMPOUND_METRIC**
  + **CONDITION_GROUP**
  + **DEGENERATE_DIMENSION**
  + **DERIVATIVE_INDEX**
  + **DERIVED_METRIC**
  + **DIMENSION**
  + **DIMENSION_ATTRIBUTE**
  + **DIMENSION_HIERARCHIES**
  + **DIMENSION_LOGIC_TABLE**
  + **DIMENSION_TABLE_ATTRIBUTE**
  + **DIRECTORY**
  + **FACT_ATTRIBUTE**
  + **FACT_DIMENSION**
  + **FACT_LOGIC_TABLE**
  + **FACT_MEASURE**
  + **FUNCTION**
  + **INFO_ARCH**
  + **MODEL**
  + **QUALITY_RULE**
  + **SECRECY_LEVEL**
  + **STANDARD_ELEMENT**
  + **STANDARD_ELEMENT_TEMPLATE**
  + **SUBJECT**
  + **SUMMARY_DIMENSION_ATTRIBUTE**
  + **SUMMARY_INDEX**
  + **SUMMARY_TIME**
  + **TABLE_MODEL**
  + **TABLE_MODEL_ATTRIBUTE**
  + **TABLE_MODEL_LOGIC**
  + **TABLE_TYPE**
  + **TAG**
  + **TIME_CONDITION**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
