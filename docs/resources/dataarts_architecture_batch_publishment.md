---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_batch_publishment"
description: |-
  Manages a DataArts Architecture batch publish resource within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_batch_publishment

Manages a DataArts Architecture batch publish resource within HuaweiCloud.

-> Destroying this resource will not change the status of the published resource.

## Example Usage

```hcl
variable "workspace_id" {}
variable "approver_user_id" {}
variable "approver_user_name" {}
variable "biz_id" {}

resource "huaweicloud_dataarts_architecture_batch_publishment" "test" {
  workspace_id       = var.workspace_id
  approver_user_id   = var.approver_user_id
  approver_user_name = var.approver_user_name
  fast_approval      = true

  biz_infos {
    biz_id   = var.biz_id
    biz_type = "SUBJECT"
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

* `biz_infos` - (Required, List, ForceNew) Specifies the list of the business information.
  Changing this creates a new resource.
  The [biz_infos](#publishment_biz_infos) structure is documented below.

  -> If the parameter contains published objects, the resource creation will fail, but the remaining objects
     will continue to be published.

* `approver_user_id` - (Required, String, ForceNew) Specifies the user ID of the architecture reviewer.
  Changing this creates a new resource.

* `approver_user_name` - (Required, String, ForceNew) Specifies the user name of the architecture reviewer.
  Changing this creates a new resource.
  
* `fast_approval` - (Optional, Bool, ForceNew) Whether to automatically review.
  Changing this creates a new resource.

  -> This parameter is available only when the current user has approval authority.

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
