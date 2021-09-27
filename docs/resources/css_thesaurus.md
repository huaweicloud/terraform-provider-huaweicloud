---
subcategory: "Cloud Search Service (CSS)"
---

# huaweicloud_css_thesaurus

Manages CSS thesaurus resource within HuaweiCloud

-> Only one thesaurus resource can be created for the specified cluster

## Example Usage

### Create a thesaurus

```hcl
resource "huaweicloud_css_thesaurus" "test" {
  cluster_id  = {{ css_cluster_id }}
  bucket_name = {{ bucket_name }}
  main_object = {{ bucket_obj_key }}
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

-> Specifies at least one of `main_object`,`stop_object`,`synonym_object`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a resource ID in UUID format.

* `status` - Indicates the status of the thesaurus loading

* `update_time` - Specifies the time (UTC) when the thesaurus was modified. The format is ISO8601:YYYY-MM-DDThh:mm:ssZ

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

CSS thesaurus can be imported by `id`. For example,

```
terraform import huaweicloud_css_thesaurus.example e9ee3f48-f097-406a-aa74-cfece0af3e31
```
