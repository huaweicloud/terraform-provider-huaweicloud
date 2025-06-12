---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_overviews_statistics"
description: |-
  Use this data source to query statistics of requests and attacks.
---

# huaweicloud_waf_overviews_statistics

Use this data source to query statistics of requests and attacks.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_waf_overviews_statistics" "test" {
  from = var.start_time
  to   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `from` - (Required, Int) Specifies the query start time.
  The format is 13-digit timestamp in millisecond.

* `to` - (Required, Int) Specifies the query end time.
  The format is 13-digit timestamp in millisecond.

-> The parameters `from` and `to` must be used together.

* `hosts` - (Optional, String) Specifies the ID of the domain.

* `instances` - (Optional, String) Specifies the ID of the dedicated WAF instances.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `statistics` - The statistics about requests and attacks.

  The [statistics](#statistics_struct) structure is documented below.

<a name="statistics_struct"></a>
The `domain` block supports:

* `key` - The type of requests or attacks.
  The options are **ACCESS** for total requests, **CRAWLER** for bot mitigation, **ATTACK** for total attacks,
  **WEB_ATTACK** for basic web protection, **PRECISE** for precise protection, and **CC** for CC attack protection.

* `num` - The number of times requests or attacks.
