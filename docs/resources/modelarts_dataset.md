---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_dataset"
description: |-
  Manages a ModelArts dataset resource within HuaweiCloud.
---

# huaweicloud_modelarts_dataset

Manages a ModelArts dataset resource within HuaweiCloud.

## Example Usage

```hcl
variable "dataset_name" {}
variable "output_obs_path" {}
variable "input_obs_path" {}

resource "huaweicloud_modelarts_dataset" "test" {
  name        = var.dataset_name
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

* `region` - (Optional, String, ForceNew) The region where the dataset is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the dataset.  
  The valid length is limited from `1` to `100`, only letters, chinese characters, digits underscores (_) and
  hyphens (-) are allowed. The name must start with a letter.

* `type` - (Required, Int, NonUpdatable) Specifies the type of dataset.  
  The valid values are as follows:
  + **0**: Image classification, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **1**: Object detection, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **3**: Image segmentation, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **100**: Text classification, supported formats: `.txt`, `.csv`.
  + **200**: Sound classification, Supported formats: `.wav`.
  + **201**: Speech content.
  + **202**: Speech paragraph labeling.
  + **400**: Table type, supported formats: Carbon type.
  + **600**: Video, supported formats: `.mp4`.
  + **900**: Free format.

* `output_path` - (Required, String, NonUpdatable) Specifies the OBS storage path that used to store output files.  
  The path cannot be the same as the import path or subdirectory of the import path.

* `data_source` - (Required, List, NonUpdatable) Specifies the data sources which be used to imported the source data
  (such as pictures/files/audio, etc.) in this directory and subdirectories to the dataset.  
  The [data_source](#modelarts_dataset_datasource) structure is documented below.

* `description` - (Optional, String) Specifies the description of the dataset.  
  The valid length is limited from `0` to `256`, and the description cannot contain special characters `!<>=&"'`.

* `schemas` - (Optional, List, NonUpdatable) Specifies the schema configurations of the dataset.  
  Required if the value of `type` parameter is **400**.  
  The [schemas](#modelarts_dataset_schema) structure is documented below.

* `import_labeled_enabled` - (Optional, Bool, NonUpdatable) Specifies whether to enable the import labeled features.  
  Defaults to **true**.

* `label_format` - (Optional, List, NonUpdatable) Specifies the custom format information of labeled features when
  import labeled files for Text classification.  
  The [label_format](#modelarts_dataset_label_format) structure is documented below.

* `labels` - (Optional, List) Specifies labels of the dataset.  
  The [labels](#modelarts_dataset_labels) structure is documented below.

<a name="modelarts_dataset_datasource"></a>
The `data_source` block supports:

* `data_type` - (Optional, Int, NonUpdatable) Specifies the type of data source.  
  The valid values are as follows:
  + **0**: OBS.
  + **1**: GaussDB(DWS).
  + **2**: DLI.
  + **4**: MRS.

  Defaults to **0**.

* `path` - (Optional, String, NonUpdatable) Specifies the OBS storage path or MRS HDFS path.  
  All files in this directory and its subdirectories will be imported into the dataset.  
  Required if the value of `data_source.data_type` parameter is **0** or **4**.

* `cluster_id` - (Optional, String, NonUpdatable) Specifies the cluster ID of the DWS/MRS cluster.  
  Required if the value of `data_source.data_type` parameter is **1** or **4**.

* `database_name` - (Optional, String, NonUpdatable) Specifies the name of the DWS/DLI database.  
  Required if the value of `data_source.data_type` parameter is **1** or **2**.

* `table_name` - (Optional, String, NonUpdatable) Specifies the name of the DWS/DLI table.  
  Required if the value of `data_source.data_type` parameter is **1** or **2**.

* `user_name` - (Optional, String, NonUpdatable) Specifies the name of the DWS database user.  
  Required if the value of `data_source.data_type` parameter is **1**.

* `password` - (Optional, String, NonUpdatable) Specifies the password of the DWS database user.  
  Required if the value of `data_source.data_type` parameter is **1**.

* `queue_name` - (Optional, String, NonUpdatable) Specifies the name of the DLI queue.  
  Required if the value of `data_source.data_type` parameter is **2**.

* `with_column_header` - (Optional, Bool, NonUpdatable) Specifies whether the data contains table header when the type
 of dataset is **400** (table type).  
  Defaults to **true**.

<a name="modelarts_dataset_schema"></a>
The `schemas` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the field type of the schema.  
  The valid values are as follows:
  + **String**
  + **Short**
  + **Int**
  + **Long**
  + **Double**
  + **Float**
  + **Byte**
  + **Date**
  + **Timestamp**
  + **Bool**

* `name` - (Required, String, NonUpdatable) Specifies the field name of the schema.

<a name="modelarts_dataset_label_format"></a>
The `label_format` block supports:

* `type` - (Optional, String, NonUpdatable) Specifies the type of the label format.  
  The valid values are as follows:
  + **0**: Label and text are separated, distinguished by the suffix `_result`.
    For example: the text file is *abc.txt*, and the label file is *abc_result.txt*.
  + **1**: Default, labels and text are in one file, separated by a delimiter.
    The separator between text and labels, the separator between label and label can be specified by `label_separator`
    and `text_label_separator`.

  Defaults to **1**.

* `text_label_separator` - (Optional, String, NonUpdatable) Specifies the separator between text and label.

* `label_separator` - (Optional, String, NonUpdatable) Specifies the separator between label and label.

<a name="modelarts_dataset_labels"></a>
The `labels` block supports:

* `name` - (Required, String) Specifies the name of the label.

* `property_color` - (Optional, String) Specifies the color of the label.

* `property_shape` - (Optional, String) Specifies the shape of the label.  
  The valid values are as follows:
  + **bndbox**
  + **polygon**
  + **circle**
  + **line**
  + **dashed**
  + **point**
  + **polyline**

* `property_shortcut` - (Optional, String) Specifies shortcut of the label.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `data_format` - The data format of the dataset.
  + **Default**
  + **CarbonData**: Carbon format(Supported only for table type datasets).

* `created_at` - The creation time of the dataset, in RFC3339 format.

* `status` - The status of the dataset.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.
* `delete` - Default is 10 minute.

## Import

The datasets can be imported by `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_dataset.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `data_source.0.password`.
It is generally recommended running `terraform plan` after importing a dataset. You can then decide if changes should be
applied to the dataset, or the resource definition should be updated to align with the dataset.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_modelarts_dataset" "test" {
  ...

  lifecycle {
    ignore_changes = [
      data_source.0.password,
    ]
  }
}
```
