---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_phone_flavors"
description: |-
  Use this data source to get available flavors of CPH phone.
---

# huaweicloud_cph_phone_flavors

Use this data source to get available flavors of CPH phone.

## Example Usage

```hcl
data "huaweicloud_cph_phone_flavors" "test" {
  type = "1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) The flavor status. Defaults to **1**.  
  The options are as follows:
  + **0**: offline.
  + **1**: normal.

* `type` - (Optional, String) The cloud phone type.  
  The options are as follows:
  + **0**: Cloud phone.
  + **1**: Cloud mobile gaming.

* `vcpus` - (Optional, Int) The vcpus of the CPH phone.

* `memory` - (Optional, Int) The ram of the CPH phone in MB.

* `server_flavor_id` - (Optional, String) The CPH server flavor.

* `image_label` - (Optional, String) The label of image.
  The valid values are **cloud_phone**, **cloud_game**, **qemu_phone**, **cloud_phone_1620**, and **cloud_game_1620**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of flavor detail.
  The [Flavors](#phoneFlavors_Flavors) structure is documented below.

<a name="phoneFlavors_Flavors"></a>
The `Flavors` block supports:

* `flavor_id` - The name of the flavor.

* `server_flavor_id` - The name of the CPH server flavor.

* `vcpus` - The vcpus of the CPH phone.

* `memory` - The ram of the CPH phone in MB.

* `disk` - The storage size in GB.

* `resolution` - The resolution of the CPH phone.

* `phone_capacity` - The number of cloud phones of the current flavor.

* `status` - The flavor status.  
  The options are as follows:
  + **0**: offline.
  + **1**: normal.

* `type` - The cloud phone type.  
  The options are as follows:
  + **0**: Cloud phone.
  + **1**: Cloud mobile gaming.

* `image_label` - (Optional, String) The label of image.
  The valid values are **cloud_phone**, **cloud_game**, **qemu_phone**, **cloud_phone_1620**, and **cloud_game_1620**.

* `extend_spec` - The extended description, which is a string in JSON format and can contain a maximum of 512 bytes.
