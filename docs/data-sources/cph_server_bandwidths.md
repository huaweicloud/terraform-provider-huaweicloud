---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_server_bandwidths"
description: |-
  Use this data source to get bandwidth list of CPH server.
---

# huaweicloud_cph_server_bandwidths

Use this data source to get bandwidth list of CPH server.

## Example Usage

```hcl
data "huaweicloud_cph_server_bandwidths" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidths` - The bandwidth list of CPH server.

  The [bandwidths](#bandwidths_struct) structure is documented below.

<a name="bandwidths_struct"></a>
The `bandwidths` block supports:

* `bandwidth_type` - The bandwidth type.

* `bandwidth_name` - The bandwidth name.

* `bandwidth_id` - The bandwidth ID.

* `bandwidth_size` - The bandwidth size.

* `bandwidth_charge_mode` - The bandwidth charge mode.
  + **0**: bandwidth
  + **1**: traffic

* `created_at` - The bandwidth creation time.

* `updated_at` - The bandwidth update time.
