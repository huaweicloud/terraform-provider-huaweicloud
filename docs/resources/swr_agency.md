---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_agency"
description: |-
  Manages a SWR agency resource within HuaweiCloud.
---

# huaweicloud_swr_agency

Manages a SWR agency resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_swr_agency" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
