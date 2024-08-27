---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_table"
description: ""
---

# huaweicloud_dli_table

Manages DLI Table resource within HuaweiCloud

## Example Usage

### Create a Table

```hcl
variable "database_name" {}

resource "huaweicloud_dli_database" "test" {
  name = var.database_name
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "table_1"
  data_location = "DLI"
  description   = "SQL table_1 description"

  columns {
    name        = "column_1"
    type        = "string"
    description = "the first column"
  }

  columns {
    name        = "column_2"
    type        = "string"
    description = "the second column"
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the dli table resource. If omitted,
  the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the table name. The name can contain only digits, letters,
 and underscores, but cannot contain only digits or start with an underscore. Length range: 1 to 128 characters.
 Changing this parameter will create a new resource.

* `database_name` - (Required, String, ForceNew) Specifies the database name which the table belongs to.
 Changing this parameter will create a new resource.

* `data_location` - (Required, String, ForceNew) Specifies data storage location. Changing this parameter will create
  a newresource. The options are as follows:
  + **DLI**: Data stored in DLI tables is applicable to delay-sensitive services, such as interactive queries.
  + **OBS**: Data stored in OBS tables is applicable to delay-insensitive services, such as historical data statistics
   and analysis.

* `description` - (Optional, String, ForceNew) Specifies description of the table.
  Changing this parameter will create a new resource.

* `columns` - (Optional, List, ForceNew) Specifies Columns of the new table. Structure is documented below.
  Changing this parameter will create a new resource.

* `data_format` - (Optional, String, ForceNew) Specifies type of the data to be added to the OBS table.
 The options: parquet, orc, csv, json, carbon, and avro. Changing this parameter will create a new resource.

* `bucket_location` - (Optional, String, ForceNew) Specifies storage path of data which will be import to the OBS table.
 Changing this parameter will create a new resource.
 -> If you need to import data stored in OBS to the OBS table, set this parameter to the path of a folder. If the table
  creation path is a file, data fails to be imported. which must be a path on OBS and must begin with obs.

* `with_column_header` - (Optional, Bool, ForceNew) Specifies whether the table header is included in the data file.
  Only data in CSV files has this attribute. Changing this parameter will create a new resource.

* `delimiter` - (Optional, String, ForceNew) Specifies data delimiter. Only data in CSV files has this
  attribute. Changing this parameter will create a new resource.

* `quote_char` - (Optional, String, ForceNew) Specifies reference character. Double quotation marks (`\`)
 are used by default. Only data in CSV files has this attribute. Changing this parameter will create a new resource.

* `escape_char` - (Optional, String, ForceNew) Specifies escape character. Backslashes (`\\`) are used by
 default. Only data in CSV files has this attribute. Changing this parameter will create a new resource.

* `date_format` - (Optional, String, ForceNew) Specifies date type. `yyyy-MM-dd` is used by default. Only
 data in CSV and JSON files has this attribute. Changing this parameter will create a new resource.

* `timestamp_format` - (Optional, String, ForceNew) Specifies timestamp type. `yyyy-MM-dd HH:mm:ss` is used by default.
 Only data in CSV and JSON files has this attribute. Changing this parameter will create a new resource.

The `column` block supports:

  * `name` - (Required, String, ForceNew) Specifies the name of column. Changing this parameter will create a new
   resource.
  * `type` - (Required, String, ForceNew) Specifies data type of column. Changing this parameter will create a new
   resource.
  * `description` - (Required, String, ForceNew) Specifies the description of column. Changing this parameter will
   create a new resource.
  * `is_partition` - (Optional, Bool, ForceNew) Specifies whether the column is a partition column. The value
    `true` indicates a partition column, and the value false indicates a non-partition column. The default value
     is false. Changing this parameter will create a new resource.
  
  -> When creating a partition table, ensure that at least one column in the table is a non-partition column.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in format of **database_name/table_name**. It is composed of the name of database which table
 belongs and the name of table, separated by a slash.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

DLI table can be imported by `id`. It is composed of the name of database which table belongs and the name of table,
 separated by a slash. For example,

```bash
terraform import huaweicloud_dli_table.example <database_name>/<table_name>
```
