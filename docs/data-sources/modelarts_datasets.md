---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_datasets"
description: ""
---

# huaweicloud_modelarts_datasets

Use this data source to get a list of ModelArts datasets.

## Example Usage

```hcl
data "huaweicloud_modelarts_datasets" "test" {
  name = "your_dataset_name"
  type = 1
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available datasets in the current region.
 All datasets that meet the filter criteria will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain datasets. If omitted, the provider-level region
 will be used.

* `name` - (Optional, String) Specifies the name of datasets.

* `type` - (Optional, Int) Specifies the type of datasets. The options are:
  + **0**: Image classification, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **1**: Object detection, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **3**: Image segmentation, supported formats: `.jpg`, `.png`, `.jpeg`, `.bmp`.
  + **100**: Text classification, supported formats: `.txt`, `.csv`.
  + **200**: Sound classification, Supported formats: `.wav`.
  + **400**: Table type, supported formats: Carbon type.
  + **600**: Video, supported formats: `.mp4`
  + **900**: Free format.

## Attribute Reference

The following attributes are exported:

* `id` - Indicates a data source ID.

* `datasets` - Indicates a list of all datasets found. Structure is documented below.

The `datasets` block contains:

* `id` - The ID of the dataset.

* `name` - The name of the dataset.

* `type` - The type of the dataset.

* `description` - The description of the dataset.

* `output_path` - The OBS path for storing output files such as labeled files.

* `data_source` - The data sources which is used to imported the source data (such as pictures/files/audio, etc.) in
 this directory and subdirectories to the dataset. Structure is documented below.

* `schemas` - The schema information of source data when `type` is `400`(Table Type). Structure is documented below.

* `labels` - The labels information. Structure is documented below.

* `data_format` - The dataset format. Valid values include: `Default`, `CarbonData`: Carbon format(Supported only for
 table type dataset.).

* `created_at` - The dataset creation time.

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

The `data_source` block contains:

* `data_type` - The type of data source. Valid values are as follows:
  + *0*: OBS.
  + *1*: GaussDB(DWS).
  + *2*: DLI.
  + *4*: MRS.
  
* `path` - The OBS path when `data_type` is `0`(OBS) or the HDFS path when `data_type` is `4`(MRS). All the file in this
 directory and subdirectories will be which be imported to the dataset.

* `with_column_header` - Whether the data contains table header when the type of dataset is `400`(Table type).

The `schemas` block contains:

* `type` - The field type. Valid values include: `String`, `Short`, `Int`, `Long`, `Double`, `Float`, `Byte`,
 `Date`, `Timestamp`, `Bool`.

* `name` - The field name.

The `labels` block contains:

* `name` - The name of label.

* `property_color` - The color of label.

* `property_shape` - The shape of label. Valid values include: `bndbox`, `polygon`, `circle`, `line`, `dashed`,
 `point`, `polyline`.

* `property_shortcut` - The shortcut of label.
