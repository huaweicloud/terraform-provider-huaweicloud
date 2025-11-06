---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_observability_configuration"
description: |-
  Manages a CCI v2 observability configuration resource within HuaweiCloud.
---

# huaweicloud_cciv2_observability_configuration

Manages a CCI v2 observability configuration resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_cciv2_observability_configuration" "test" {
  event {
    enable = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `event` - (Required, List) Specifies the observability configuration.
  The [event](#event_struct) structure is documented below.

<a name="event_struct"></a>
The `event` block supports:

* `enable` - (Required, Bool) Specifies whether the observability is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the domain ID.

## Import

The CCI v2 observability configuration can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cciv2_observability_configuration.test <id>
```
