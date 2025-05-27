---
subcategory: "EVS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volume_metadata"
description: |-
  Use this resource to manage EVS volume metadata within HuaweiCloud.
---

# huaweicloud_evs_volume_metadata

Use this resource to manage EVS volume metadata within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_evs_volume_metadata" "test" {
  volume_id = volume_id
  
  metadata {
    key1 = "value1"
    key2 = "value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `volume_id` - (Required, String) Specifies the disk ID.

* `metadata` - (Required, Map) Specifies the snapshot metadata, which is made up of key-value pairs.
  The [metadata](#metadata_struct) structure is documented below.

<a name="metadata_struct"></a>
The `metadata` block supports:

* `key` - (Optional, String) Specifies the key of the metadata. Possible values are:

* `value` - (Optional, String) Specifies the value of the metadata.

  -> The field has the following restriction:
  <br/>`key` or `value` under metadata can contain no more than `85` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.
