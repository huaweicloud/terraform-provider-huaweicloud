---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_protect_data_types"
description: |-
  Use this data source to get the list of data protection lifecycle types within HuaweiCloud.
---

# huaweicloud_dsc_protect_data_types

Use this data source to get the list of data protection lifecycle types within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_protect_data_types" "test" {
  life_cycle = "TRANSMISSION"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `life_cycle` - (Required, String) Specifies the data lifecycle stage.
  The valid values are as follows:
  + `COLLECTION`: Collection.
  + `TRANSMISSION`: Transmission.
  + `STORAGE`: Storage.
  + `USAGE`: Usage.
  + `SHARING`: Sharing.
  + `DESTRUCTION`: Destruction.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of data protection type details.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `category` - The category of the data protection measure.

* `create_time` - The creation timestamp.

* `data_type` - The data protection type.

* `data_type_id` - The classification ID when the protection type is LEVEL.

* `protect_id` - The unique ID of the data protection measure.

* `is_deleted` - Whether the record is logically deleted.

* `life_cycle` - The lifecycle stage to which it belongs.

* `update_time` - The update timestamp.
