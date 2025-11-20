---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_protection_types"
description: |-
  Use this data source to get the list of protection types.
---

# huaweicloud_waf_rules_protection_types

Use this data source to get the list of protection types.

## Example Usage

```hcl
data "huaweicloud_waf_rules_protection_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The protection types list.
