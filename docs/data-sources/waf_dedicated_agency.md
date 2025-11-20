---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_dedicated_agency"
description: |-
  Use this data source to get the dedicated WAF agency within HuaweiCloud.
---

# huaweicloud_waf_dedicated_agency

Use this data source to get the dedicated WAF agency within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_waf_dedicated_agency" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `agency_id` - The agency ID.

* `name` - The agency name.

* `duration` - The agent existence time period.

* `domain_id` - The domain ID.

* `is_valid` - Whether the agency is legal。

* `version` - The version.
