---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_database_watermark_embed_task"
description: |-
  Manages a DSC database watermark embed task resource within HuaweiCloud.
---

# huaweicloud_dsc_database_watermark_embed_task

Manages a DSC database watermark embed task resource within HuaweiCloud.

## Example Usage

### Create a database fake-column watermark embed task and immediately execute it

```hcl
variable "task_name" {}
variable "water_mark" {}
variable "database_id" {}
variable "database_name" {}
variable "instance_id" {}
variable "instance_name" {}
variable "source_db_table_name" {}
variable "target_db_table_name" {}
variable "watermark_key" {}

resource "huaweicloud_dsc_database_watermark_embed_task" "test" {
  task_name         = var.task_name
  water_mark        = var.water_mark
  watermark_version = "V2"

  db_water_param {
    embed_mode    = "EMBED_FAKE_COLUMN"
    watermark_key = var.watermark_key

    params {
      new_column_name = "field1"
      new_column_type = "date"
      fake_strategy   = "date"

      fake_param {
        date_begin = "2026-08-03T14:20:55+08:00"
        date_end   = "2026-08-04T14:20:55+08:00"
      }
    }
    params {
      new_column_name = "field2"
      new_column_type = "number"
      fake_strategy   = "number_random"

      fake_param {
        random_begin    = "2"
        random_end      = "3"
        random_accuracy = 1
      }
    }
  }

  source_db_info {
    db_id      = var.database_id
    db_name    = var.database_name
    db_type    = "MySQL"
    ins_id     = var.instance_id
    ins_name   = var.instance_name
    table_name = var.source_db_table_name
  }

  target_db_info {
    db_id      = var.database_id
    db_name    = var.database_name
    db_type    = "MySQL"
    ins_id     = var.instance_id
    ins_name   = var.instance_name
    table_name = var.target_db_table_name
  }

  error_code    = 1
  start_now     = true
  schedule_type = "ONCE"
}
```

### Create a database fake-row watermark embed task and execute it at a specific time

```hcl
variable "task_name" {}
variable "water_mark" {}
variable "row_spacing" {}
variable "watermark_key" {}
variable "instance_id" {}
variable "instance_name" {}
variable "database_id" {}
variable "database_name" {}
variable "source_db_table_name" {}
variable "target_db_table_name" {}
variable "start_time" {}

resource "huaweicloud_dsc_database_watermark_embed_task" "test" {
  task_name         = var.task_name
  water_mark        = var.water_mark
  watermark_version = "V2"

  db_water_param {
    embed_mode    = "EMBED_FAKE_ROW"
    row_spacing   = var.row_spacing
    watermark_key = var.watermark_key
  }

  source_db_info {
    db_id      = var.database_id
    db_name    = var.database_name
    db_type    = "MySQL"
    ins_id     = var.instance_id
    ins_name   = var.instance_name
    table_name = var.source_db_table_name
  }

  target_db_info {
    db_id      = var.database_id
    db_name    = var.database_name
    db_type    = "MySQL"
    ins_id     = var.instance_id
    ins_name   = var.instance_name
    table_name = var.target_db_table_name
  }

  error_code    = 1
  schedule_type = "DAY"
  start_time    = var.start_time
}
```

### Create a database lossy column watermark embed task

```hcl
variable "task_name" {}
variable "water_mark" {}
variable "watermark_key" {}
variable "instance_id" {}
variable "instance_name" {}
variable "database_id" {}
variable "database_name" {}
variable "schema_name" {}
variable "source_db_table_name" {}
variable "target_db_table_name" {}
variable "start_time" {}
variable "selected_fields" {
  type = list(object({
    column_name = string
    column_type = string
  }))
}

resource "huaweicloud_dsc_database_watermark_embed_task" "test" {
  task_name         = var.task_name
  water_mark        = var.water_mark
  watermark_version = "V2"

  db_water_param {
    embed_mode    = "EMBED_COLUMN"
    watermark_key = var.watermark_key
  }

  source_db_info {
    db_id       = var.database_id
    db_name     = var.database_name
    db_type     = "DWS"
    ins_id      = var.instance_id
    ins_name    = var.instance_name
    schema_name = var.schema_name
    table_name  = var.source_db_table_name
  }

  target_db_info {
    db_id       = var.database_id
    db_name     = var.database_name
    db_type     = "DWS"
    ins_id      = var.instance_id
    ins_name    = var.instance_name
    schema_name = var.schema_name
    table_name  = var.target_db_table_name
  }

  dynamic "selected_fields" {
    for_each = var.selected_fields

    content {
      column_name = selected_fields.value.column_name
      column_type = selected_fields.value.column_type
    }
  }

  error_code    = 1
  schedule_type = "DAY"
  start_time    = var.start_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the database watermark embed task is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `task_name` - (Required, String) Specifies the name of the database watermark embed task.  
  The maximum length is `255` characters and must be unique.  
  Only Chinese characters, English letters, digits, underscores (_) and hyphens (-) are allowed.

* `water_mark` - (Required, String) Specifies the watermark content to be embedded into the database.

* `watermark_version` - (Required, String) Specifies the watermark algorithm version.  
  The valid values are as follows:
  + **V1**
  + **V2**

* `db_water_param` - (Required, List) Specifies the database watermark embedding configuration.  
  The [db_water_param](#database_watermark_embed_task_db_water_param) structure is documented below.

* `source_db_info` - (Required, List) Specifies the source database information from which data is read.  
  The [source_db_info](#database_watermark_embed_task_db_info) structure is documented below.

* `target_db_info` - (Required, List) Specifies the target database information to which watermarked data is written.  
  The [target_db_info](#database_watermark_embed_task_db_info) structure is documented below.

* `error_code` - (Required, Int) Specifies the watermark error correction level.  
  The valid values are as follows:
  + **1**
  + **2**
  + **3**
  + **4**

  The smaller the value, the higher the error correction ability.

* `selected_fields` - (Optional, List) Specifies the selected field list used for watermark embedding.  
  The [selected_fields](#database_watermark_embed_task_selected_fields) structure is documented below.  
  This parameter is required when `db_water_param.embed_mode` is set to **EMBED_COLUMN**.

* `watermark_describe` - (Optional, String) Specifies the description of the watermark.

* `schedule_switch` - (Optional, Bool) Specifies whether to enable task scheduling.  
  Defaults to **true**.

* `schedule_type` - (Optional, String) Specifies the schedule type of the task.  
  The valid values are as follows:
  + **ONCE**
  + **DAY**
  + **WEEK**
  + **MONTH**

* `start_now` - (Optional, Bool) Specifies whether to start the watermark embed task immediately after creation.  
  Defaults to **false**.  
  This parameter is valid only when `schedule_type` is set to **ONCE**.

* `start_time` - (Optional, String) Specifies the scheduled start time of the task, in RFC3339 format.  
  `start_now` and `start_time` cannot be specified at the same time.  
  When `start_now` is set to **true**, the value of this parameter is determined by the actual start time of the task.

<a name="database_watermark_embed_task_db_water_param"></a>
The `db_water_param` block supports:

* `embed_mode` - (Required, String) Specifies the watermark embed mode.  
  The valid values are as follows:
  + **EMBED_FAKE_ROW**: Lossless fake-row watermark.
  + **EMBED_FAKE_COLUMN**: Lossless fake-column watermark.
  + **EMBED_COLUMN**: Lossy column watermark.

* `watermark_key` - (Optional, String) Specifies the watermark key used to embed and extract the watermark.

* `params` - (Optional, List) Specifies the fake-column embed parameter list.  
  The [params](#database_watermark_embed_task_params) structure is documented below.  
  This parameter is required when `embed_mode` is set to **EMBED_FAKE_COLUMN**.

* `row_spacing` - (Optional, String) Specifies the row spacing used by fake-row watermark.  
  This parameter is required when `embed_mode` is set to **EMBED_FAKE_ROW**.  
  This value must be an integer greater than `1`.

<a name="database_watermark_embed_task_params"></a>
The `params` block supports:

* `new_column_name` - (Required, String) Specifies the name of the new fake column to be created.

* `new_column_type` - (Required, String) Specifies the data type of the new fake column.  
  The valid values are as follows:
  + **varchar**
  + **date**
  + **number**

* `fake_strategy` - (Optional, String) Specifies the strategy used to generate fake data.

* `fake_param` - (Optional, List) Specifies the configuration of fake data generation.  
  The [fake_param](#database_watermark_embed_task_fake_param) structure is documented below.  
  This parameter is used together with `fake_strategy`.

<a name="database_watermark_embed_task_fake_param"></a>
The `fake_param` block supports:

* `address_accuracy` - (Optional, String) Specifies the accuracy of generated address data.

* `date_begin` - (Optional, String) Specifies the start date of the generated date range, in RFC3339 format.

* `date_end` - (Optional, String) Specifies the end date of the generated date range, in RFC3339 format.

* `random_accuracy` - (Optional, Int) Specifies the precision of generated random numbers.

* `random_begin` - (Optional, String) Specifies the lower bound of the generated random value range.

* `random_distribute` - (Optional, String) Specifies the distribution mode of generated random values.

* `random_end` - (Optional, String) Specifies the upper bound of the generated random value range.

* `string_distribute` - (Optional, String) Specifies the distribution mode of generated string values.

<a name="database_watermark_embed_task_selected_fields"></a>
The `selected_fields` block supports:

* `column_name` - (Optional, String) Specifies the name of the selected column.

* `column_type` - (Optional, String) Specifies the data type of the selected column.

<a name="database_watermark_embed_task_db_info"></a>
The `source_db_info` and `target_db_info` blocks support:

* `db_id` - (Required, String) Specifies the ID of the authorized database.

* `db_name` - (Required, String) Specifies the name of the database.

* `db_type` - (Required, String, NonUpdatable) Specifies the database type.  
  When `db_water_param.embed_mode` is set to **EMBED_COLUMN**, the valid values are as follows:
  + **DWS**
  + **MRS_HIVE**

  When `db_water_param.embed_mode` is set to **EMBED_FAKE_COLUMN** or **EMBED_FAKE_ROW**, the valid values are as follows:
  + **DWS**
  + **PostgreSQL**
  + **MySQL**

* `ins_id` - (Required, String, NonUpdatable) Specifies the ID of the database instance.

* `ins_name` - (Required, String, NonUpdatable) Specifies the name of the database instance.

* `table_name` - (Required, String) Specifies the name of the database table.

* `schema_name` - (Optional, String) Specifies the schema name of the database.  
  This parameter is valid only when `db_type` is **DWS** or **PostgreSQL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `task_state` - The running state of the watermark embed task.  
  + **WAIT**
  + **FINISHED**
  + **CLOSED**
  + **ERROR**

* `task_create_time` - The creation time of the task, in RFC3339 format.

* `task_end_time` - The end time of the task, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dsc_database_watermark_embed_task.test <id>
```
