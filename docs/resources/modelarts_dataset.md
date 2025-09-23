---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_dataset"
description: ""
---

# huaweicloud_modelarts_dataset

Manages ModelArts dataset resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "output_obs_path" {}
variable "input_obs_path" {}

resource "huaweicloud_modelarts_dataset" "test" {
  name        = var.name
  type        = 1
  output_path = var.output_obs_path
  description = "Terraform Demo"

  data_source {
    path = var.input_obs_path
  }

  labels {
    name = "foo"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
 provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the dataset. The name consists of `1` to `100` characters,
 starting with a letter. Only letters, chinese characters, digits underscores (_) and hyphens (-) are allowed.

* `type` - (Required, Int, ForceNew) Specifies the type of dataset. The options are as follows:
  + **0**: Image classification, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **1**: Object detection, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **3**: Image segmentation, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **100**: Text classification, supported formats: `.txt`, `.csv`.
  + **200**: Sound classification, Supported formats: `.wav`.
  + **400**: Table type, supported formats: Carbon type.
  + **600**: Video, supported formats: `.mp4`.
  + **900**: Free format.

 Changing this parameter will create a new resource.

* `output_path` - (Required, String, ForceNew) Specifies the OBS path for storing output files such as labeled files.
 The path cannot be the same as the import path or subdirectory of the import path.
 Changing this parameter will create a new resource.

* `data_source` - (Required, List, ForceNew)Specifies the data sources which be used to imported the source data (such
 as pictures/files/audio, etc.) in this directory and subdirectories to the dataset. Structure is documented below.
 Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of dataset. It contains a maximum of `256` characters and
 cannot contain special characters `!<>=&"'`.

* `import_labeled_enabled` - (Optional, Bool, ForceNew) Specifies whether to import labeled files.
 Default value is `true`. Changing this parameter will create a new resource.

* `schemas` - (Optional, List, ForceNew) Specifies the schema information of source data when `type` is `400`.
 Structure is documented below. Changing this parameter will create a new resource.

* `label_format` - (Optional, List, ForceNew) Specifies the custom format information of labeled files when import
 labeled files for Text classification. Structure is documented below.
 Changing this parameter will create a new resource.

* `labels` - (Optional, List) Specifies labels information. Structure is documented below.

The `data_source` block supports:

* `data_type` - (Optional, Int, ForceNew) Specifies the type of data source. The options are as follows:
  + **0**: OBS.
  + **1**: GaussDB(DWS).
  + **2**: DLI.
  + **4**: MRS.
  
 Default value is 0. Changing this parameter will create a new resource.

* `path` - (Optional, String, ForceNew) Specifies the OBS path when `data_type` is `0`
 or the hdsf path when `data_type` is `4`. All the file in this directory and subdirectories will be which be imported
 to the dataset. Changing this parameter will create a new resource.

* `with_column_header` - (Optional, Bool, ForceNew) Specifies whether the data contains table header when the type
 of dataset is `400`(Table type). Default value is `true`. Changing this parameter will create a new resource.

* `queue_name` - (Optional, String, ForceNew) Specifies the queue name of DLI when `data_type` is `2`.
 Changing this parameter will create a new resource.

* `database_name` - (Optional, String, ForceNew) Specifies the database name of DWS/DLI when `data_type` is `1` or `2`.
 Changing this parameter will create a new resource.

* `table_name` - (Optional, String, ForceNew) Specifies the table name of DWS/DLI when `data_type` is `1` or `2`.
 Changing this parameter will create a new resource.

* `cluster_id` - (Optional, String, ForceNew) Specifies the cluster ID of DWS/MRS when `data_type` is `1` or `4`.
 Changing this parameter will create a new resource.

* `user_name` - (Optional, String, ForceNew) Specifies the user name of database when `data_type` is `1`.
 Changing this parameter will create a new resource.

* `password` - (Optional, String, ForceNew) Specifies the password of database when `data_type` is `1`.
 Changing this parameter will create a new resource.

The `schemas` block supports:

* `type` - (Required, String, ForceNew) Specifies the field type. Valid values include: `String`, `Short`, `Int`,
 `Long`, `Double`, `Float`, `Byte`, `Date`, `Timestamp`, `Bool`. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the field name. Changing this parameter will create a new resource.

The `label_format` block supports:

* `type` - (Optional, String, ForceNew) Specifies Label type for text classification.
 The optional values are as follows:

  + **0**: Label and text are separated, distinguished by the suffix `_result`.
   For example: the text file is *abc.txt*, and the label file is *abc_result.txt*.
  + **1**: Default, labels and text are in one file, separated by a delimiter. The separator between text and labels,
   the separator between label and label can be specified by `label_separator` and `text_label_separator`.
  
 Default value is `1`.

* `text_label_separator` - (Optional, String, ForceNew) Specifies the separator between text and label.
 Changing this parameter will create a new resource.

* `label_separator` - (Optional, String, ForceNew) Specifies the separator between label and label.
 Changing this parameter will create a new resource.

The `labels` block supports:

* `name` - (Required, String) Specifies the name of label.

* `property_color` - (Optional, String) Specifies color of label.

* `property_shape` - (Optional, String) Specifies shape of label. Valid values include: `bndbox`, `polygon`,
 `circle`, `line`, `dashed`, `point`, `polyline`.

* `property_shortcut` - (Optional, String) Specifies shortcut of label.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `data_format` - dataset format. Valid values include: `Default`, `CarbonData`: Carbon format(Supported only for
 table type datasets).

* `status` - Dataset status. Valid values are as follows:
  + **0**: Creating.
  + **1**: Completed.
  + **2**: Deleting.
  + **3**: Deleted.
  + **4**: Exception.
  + **5**: Syncing.
  + **6**: Releasing.
  + **7**: Version switching.
  + **8**: Importing.

* `created_at` - The dataset creation time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.
* `delete` - Default is 10 minute.

## Import

The datasets can be imported by `id`.

```bash
terraform import huaweicloud_modelarts_dataset.test yiROKoTTjtwjvP71yLG
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `data_source.0.path`,
`data_source.0.queue_name`, `data_source.0.database_name`, `data_source.0.table_name`, `data_source.0.cluster_id`,
`data_source.0.user_name` and `data_source.0.password`. It is generally recommended running `terraform plan` after
importing a dataset. You can then decide if changes should be applied to the dataset, or the resource definition
should be updated to align with the dataset. Also you can ignore changes as below.

```hcl
resource "huaweicloud_modelarts_dataset" "test" {
    ...

  lifecycle {
    ignore_changes = [
      data_source.0.path, data_source.0.queue_name, data_source.0.database_name, data_source.0.table_name,
      data_source.0.cluster_id, data_source.0.user_name, data_source.0.password,
    ]
  }
}
```
