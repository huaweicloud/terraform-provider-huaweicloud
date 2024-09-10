---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_dataset_version"
description: ""
---

# huaweicloud_modelarts_dataset_version

Manages ModelArts dataset version resource within HuaweiCloud.

## Example Usage

```hcl
variable "dataset_id" {}

resource "huaweicloud_modelarts_dataset_version" "v001" {
  name        = "v001"
  dataset_id  = var.dataset_id
  description = "Created by demo"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the dataset version. The name consists of `1` to `32`
  characters. Only letters, Chinese characters, digits underscores (_) and hyphens (-) are allowed.
  Changing this parameter will create a new resource.

* `dataset_id` - (Required, String, ForceNew) Specifies the ID of dataset.
  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of dataset version. It contains a maximum of
  `256` characters and cannot contain special characters `!<>=&"'`. Changing this parameter will create a new resource.

* `split_ratio` - (Optional, String, ForceNew) Specifies the ratio of splitting which randomly divides a labeled sample
  into a training set and a validation set. Changing this parameter will create a new resource.

-> Before you enable splitting, ensure each label has at least five labeled samples. Ensure there are at least two
  multi-label samples, if any.

* `hard_example` - (Optional, Bool, ForceNew) Specifies whether to enable ModelArts to write the hard example
  attributes (difficult, hard-coefficient, and hard-reasons) into the XML and manifest labeling files. ModelArts will
  use these attributes to optimize hard example filtering. Default value is `false`.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of **dataset_id/version_id**. It is composed of dataset ID and version ID,
  separated by a slash.

* `version_id` - The version ID.

* `status` - The status of the dataset version. Valid values are as follows:
  + **0**: Creating.
  + **1**: Normal.
  + **2**: Deleting.
  + **3**: Deleted.
  + **4**: Exception.

* `verification` - Whether the data has been verified by the verification algorithm before publishing.

* `labeling_type` - The label type of the dataset version. Valid values are as follows:
  + **multi**: Indicates that there are multi-label samples.
  + **single**: Indicates that all samples are single-label.
  + **unlabeled**: Indicates that all samples are unlabeled.

* `files` - The total number of samples.

* `storage_path` - The path to save the manifest file of the version.

* `is_current` - Whether this version is current version.

* `created_at` - The creation time, in UTC format.

* `updated_at` - The last update time, in UTC format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.

## Import

The dataset versions can be imported by dataset ID and version ID, separated by a slash, e.g.

```bash
terraform import huaweicloud_modelarts_dataset_version.test yiROKoTTjtwjvP71yLG/wieeeoTrtrtjvn67yLm
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `hard_example`. It is generally
recommended running `terraform plan` after importing a dataset. You can then decide if changes should be applied to the
dataset, or the resource definition should be updated to align with the dataset. Also you can ignore changes as below.

```hcl
resource "huaweicloud_modelarts_dataset_version" "test" {
    ...

  lifecycle {
    ignore_changes = [
      hard_example,
    ]
  }
}
```
