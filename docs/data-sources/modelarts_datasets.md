---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_datasets"
description: ""
---

# huaweicloud_modelarts_datasets

Use this data source to get a list of ModelArts datasets.

## Example Usage

### Query all datasets but each dataset does not contain lables

```hcl
data "huaweicloud_modelarts_datasets" "test" {}
```

### Querying a list of datasets by fuzzy matching of names

```hcl
variable "matched_name_word" {}

data "huaweicloud_modelarts_datasets" "test" {
  name = var.matched_name_word
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the datasets are located.  
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of datasets, in fuzzy match.

  -> When querying a specific dataset by name, labels information will be returned.

* `type` - (Optional, Int) Specifies the type of datasets.  
  The valid values are as follows:
  + **0**: Image classification, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **1**: Object detection, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **3**: Image segmentation, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **100**: Text classification, supported formats: `.txt`, `.csv`.
  + **200**: Sound classification, Supported formats: `.wav`.
  + **201**: Speech content.
  + **202**: Speech paragraph labeling.
  + **400**: Table type, supported formats: Carbon type.
  + **600**: Video, supported formats: `.mp4`
  + **900**: Free format.

## Attribute Reference

The following attributes are exported:

* `id` - Indicates a data source ID.

* `datasets` - The list of datasets that match the filter parameters.  
  The [datasets](#modelarts_datasets_attr) structure is documented below.

<a name="modelarts_datasets_attr"></a>
The `datasets` block contains:

* `id` - The ID of the dataset.

* `name` - The name of the dataset.

* `type` - The type of the dataset.

* `description` - The description of the dataset.

* `output_path` - The OBS storage path that used to store output files.

* `data_source` - The data sources which be used to imported the source data (such as pictures/files/audio, etc.) in
  this directory and subdirectories to the dataset.  
  The [data_source](#modelarts_datasets_datasource_attr) structure is documented below.

* `schemas` - The schema configurations of the dataset.  
  Returns if the value of `type` parameter is `400` (table Type).  
  The [schemas](#modelarts_datasets_schemas_attr) structure is documented below.

* `labels` - The labels of the dataset.  
  The [labels](#modelarts_datasets_labels_attr) structure is documented below.

* `data_format` - The dataset format.
  + **Default**
  + **CarbonData**: Carbon format(Supported only for table type datasets).

* `created_at` - The creation time of the dataset, in RFC3339 format.

* `status` - The status of the dataset.

<a name="modelarts_datasets_datasource_attr"></a>
The `data_source` block contains:

* `data_type` - The type of the data source.
  + **0**: OBS.
  + **1**: GaussDB(DWS).
  + **2**: DLI.
  + **4**: MRS.

* `path` - The OBS storage path or MRS HDFS path.

* `cluster_id` - The ID of the DWS/MRS cluster.

* `database_name` - The name of the DWS/DLI database.

* `table_name` - The name of the DWS/DLI table.

* `user_name` - The name of the DWS database user.

* `queue_name` - The name of the DLI queue.

* `with_column_header` - Whether the data contains table header when the type of dataset is `400` (table type).

<a name="modelarts_datasets_schemas_attr"></a>
The `schemas` block contains:

* `type` - The field type of the schema.
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

* `name` - The field name of the schema.

<a name="modelarts_datasets_labels_attr"></a>
The `labels` block contains:

* `name` - The name of the label.

* `property_color` - The color of the label.

* `property_shape` - The shape of the label.
  + **bndbox**
  + **polygon**
  + **circle**
  + **line**
  + **dashed**
  + **point**
  + **polyline**

* `property_shortcut` - The shortcut of the label.
