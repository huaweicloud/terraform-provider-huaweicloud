---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_source_ips"
description: |-
  Use this data source to get a list of WAF Back-to-Source IP Addresses.
---

# huaweicloud_waf_source_ips

Use this data source to get a list of WAF Back-to-Source IP Addresses.

## Example Usage

```hcl
data "huaweicloud_waf_source_ips" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `source_ips` - Origin server information list.

  The [source_ips](#source_ips_struct) structure is documented below.

<a name="source_ips_struct"></a>
The `source_ips` block supports:

* `ips` - WAF retrieval IP addresses.

* `update_time` - Time the WAF IP addresses are updated, in RFC3339 format.
