---
subcategory: "Enterprise Router (ER)"
---

# huaweicloud_er_availability_zones

Use this data source to query the list of Availability Zones in which the Enterprise Router instance can be created.

## Example Usage

```HCL
resource "huaweicloud_er_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `names` - The name list of the availability zones.
