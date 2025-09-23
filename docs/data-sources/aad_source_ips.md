---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_source_ips"
description: |-
  Use this data source to get the list of Advanced Anti-DDos source IP addresses within HuaweiCloud.
---

# huaweicloud_aad_source_ips

Use this data source to get the list of Advanced Anti-DDos source IP addresses within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_source_ips" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ips` - The list of back-to-origin IP address.
