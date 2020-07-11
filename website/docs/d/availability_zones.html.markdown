---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_availability_zones"
sidebar_current: "docs-huaweicloud-datasource-availability-zones"
description: |-
  Get a list of availability zones from HuaweiCloud
---

# huaweicloud\_availability\_zones

Use this data source to get a list of availability zones from HuaweiCloud

## Example Usage

```hcl
data "huaweicloud_availability_zones" "zones" {}
```

## Argument Reference

* `state` - (Optional) The `state` of the availability zones to match, default ("available").


## Attributes Reference

`id` is set to hash of the returned zone list. In addition, the following attributes
are exported:

* `names` - The names of the availability zones, ordered alphanumerically, that match the queried `state`
