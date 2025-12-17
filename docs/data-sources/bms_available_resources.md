---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_available_resources"
description: |-
  Use this data source to get the BMS available resources.
---

# huaweicloud_bms_available_resources

Use this data source to get the BMS available resources.

## Example Usage

```hcl
data "huaweicloud_bms_available_resources" "demo" {
  availability_zone = ["cn-north-4a"]
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Required, List) Specifies the list of availability zone.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `available_resource` - Indicates the available resource.
  The [available_resource](#available_resource_struct) structure is documented below.

<a name="available_resource_struct"></a>
The `available_resource` block supports:

* `availability_zone` - Indicates the availability zone.

* `flavors` - Indicates the flavors.
  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `flavor_id` - Indicates the ID of the flavor.

* `status` - Indicates the status of the flavor.
