---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_model"
description: ""
---

# huaweicloud_dataarts_architecture_model

Manages DataArts Architecture ER Model resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "name" {}

resource "huaweicloud_dataarts_architecture_model" "physical_model"{
  workspace_id = var.workspace_id
  name         = var.name
  type         = "THIRD_NF"
  physical     = true
  dw_type      = "DWS"
}

resource "huaweicloud_dataarts_architecture_model" "logic_model"{
  workspace_id = var.workspace_id
  name         = var.name
  type         = "THIRD_NF"
  physical     = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the model.
  If omitted, the provider-level region will be used. Changing this parameter will create a new model.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID which the model in.
  Changing this parameter will create a new model.

* `physical` - (Required, Bool, ForceNew) Specifies the model is physical or logical.
  When the value is **true**, it means physical model.
  Changing this parameter will create a new model.

* `name` - (Required, String) Specifies the model name.

* `type` - (Required, String) Specifies the model type. The valid values are **THIRD_NF** and **DIMENSION**.

* `description` - (Optional, String) Specifies the description of model.

* `dw_type` - (Optional, String) Specifies the data connection type. This parameter is mandatory when
  `physical` is **true**. The valid values are:
  + **DWS**
  + **DLI**
  + **MRS_HIVE**
  + **POSTGRESQL**
  + **MRS_SPARK**
  + **CLICKHOUSE**
  + **MYSQL**
  + **ORACLE**

* `level` - (Optional, String) Specifies the data warehouse layer. Valid values are **SDI** and **DWI**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `dw_type` - The data connection type.

* `created_at` - The create time of the model.

* `updated_at` - The update time of the model.

* `created_by` - The person creating the model.

* `updated_by` - The person updating the model.

## Import

DataArts Studio architecture model can be imported using `<workspace_id>/<name>`, e.g.

```sh
terraform import huaweicloud_dataarts_architecture_model.test <workspace_id>/<name>
```
