---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_thesaurus"
description: ""
---

# huaweicloud_css_thesaurus

Manages CSS thesaurus resource within HuaweiCloud

-> Only one thesaurus resource can be created for the specified cluster.

## Example Usage

### Create a thesaurus

```hcl
variable "cluster_id" {}
variable "bucket_name" {}
variable "bucket_obj_key" {}

resource "huaweicloud_css_thesaurus" "test" {
  cluster_id  = var.cluster_id
  bucket_name = var.bucket_name
  main_object = var.bucket_obj_key
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the thesaurus resource. If omitted, the
  provider-level region will be used. Changing this creates a new thesaurus resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the CSS cluster ID for configuring the thesaurus.
  Changing this parameter will create a new resource.

* `bucket_name` - (Required, String, ForceNew) Specifies the OBS bucket where the thesaurus files are stored
 (the bucket type must be standard storage or low-frequency storage, and archive storage is not supported).

* `main_object` - (Optional, String) Specifies the path of the main thesaurus file object.

* `stop_object` - (Optional, String) Specifies the path of the stop word library file object.

* `synonym_object` - (Optional, String) Specifies the path of the synonyms thesaurus file object.

-> Specifies at least one of `main_object`, `stop_object`, `synonym_object`.
**nil** or **Default** indicates no change, **""** or **Unused** indicates that this value is cleared.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a resource ID in UUID format.

* `status` - Indicates the status of the thesaurus loading

* `update_time` - Specifies the time (UTC) when the thesaurus was modified. The format is ISO8601:YYYY-MM-DDThh:mm:ssZ

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

CSS thesaurus can be imported by `id`, e.g.

```bash
terraform import huaweicloud_css_thesaurus.test <id>
```
