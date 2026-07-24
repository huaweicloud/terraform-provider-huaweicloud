---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_measure_info"
description: |-
  Manages a DSC measure info resource within HuaweiCloud.
---

# huaweicloud_dsc_measure_info

Manages a DSC measure info resource within HuaweiCloud.

## Example Usage

```hcl
variable "life_cycle" {}

resource "huaweicloud_dsc_measure_info" "test" {
  name         = "test_measure"
  life_cycle   = var.life_cycle
  description  = "test description"
  measure_type = "CUSTOM"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the measure name.

* `life_cycle` - (Required, String, NonUpdatable) Specifies the data lifecycle stage.
  Valid values are **COLLECTION**, **TRANSMISSION**, **STORAGE**, **USAGE**, **SHARING**, **DESTRUCTION**.

* `description` - (Optional, String) Specifies the measure description.

* `measure_type` - (Optional, String, NonUpdatable) Specifies the measure type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `is_deleted` - Whether the measure is deleted.

* `create_time` - The creation time.

* `update_time` - The update time.

## Import

The measure info can be imported using the `id` and `life_cycle`, separated by a slash, e.g.

```
$ terraform import huaweicloud_dsc_measure_info.test 123/COLLECTION
```
