---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_dataset_versions"
description: ""
---

# huaweicloud_modelarts_dataset_versions

Use this data source to get a list of ModelArts dataset versions.

## Example Usage

```hcl
variable "dataset_id" {}
data "huaweicloud_modelarts_dataset_versions" "test" {
  dataset_id = var.dataset_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available dataset versions in the current region.
All dataset versions that meet the filter criteria will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain dataset versions. If omitted, the
provider-level region will be used.

* `dataset_id` - (Required, String) Specifies the ID of dataset.

* `split_ratio` - (Optional, String) Specifies the range of splitting ratio which randomly divides a labeled sample
into a training set and a validation set. Separate the minimum and maximum split ratios with commas,
for example: "0.0,1.0".

* `name` - (Optional, String) Specifies the name of the dataset version.

## Attribute Reference

The following attributes are exported:

* `id` - Indicates a data source ID.

* `versions` - Indicates a list of all dataset versions found. Structure is documented below.

The `versions` block contains:

* `id` - The ID of the dataset version.

* `name` - The name of the dataset version.

* `description` - The description of the dataset version.

* `split_ratio` - The ratio of splitting which randomly divides a labeled sample into a training set and
a validation set.

* `status` - Dataset version status. Valid values are as follows:
  + **0**: Creating.
  + **1**: Normal.
  + **2**: Deleting.
  + **3**: Deleted.
  + **4**: Exception.
  
* `files` - The total number of samples.

* `storage_path` - The path to save the manifest file of the version.

* `is_current` - Whether this version is current version.

* `created_at` - The creation time, in UTC format.

* `updated_at` - The last update time, in UTC format.
