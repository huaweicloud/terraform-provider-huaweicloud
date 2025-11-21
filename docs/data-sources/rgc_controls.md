---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_controls"
description: |-
  Use this data source to list controls in Resource Governance Center.
---

# huaweicloud_rgc_controls

Use this data source to list controls in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_controls" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `controls` - Information about the controls list.

The [controls](#controls) structure is documented below.

<a name="controls"></a>
The `controls` block supports:

* `identifier` - The identifier of the control.

* `name` - The name of the control.

* `description` - The description of the control.

* `guidance` - The guidance of the control.

* `resource` - The resources associated with the control.

* `framework` - The frameworks associated with the control.

* `service` - The service associated with the control.

* `implementation` - The implementation of the control.

* `behavior` - The behavior of the control.

* `owner` - The owner of the control.

* `severity` - The severity of the control.

* `control_objective` - The objective of the control.

* `version` - The version of the control.

* `release_date` - The release date of the control.
