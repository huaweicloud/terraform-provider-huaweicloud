---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_overviews_abnormal"
description: |-
  Use this data source to query the top service exceptions.
---

# huaweicloud_waf_overviews_abnormal

Use this data source to query the top service exceptions.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_waf_overviews_abnormal" "test" {
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

* `top` - (Optional, Int) The first several results to query.
  The valid value ranges from `1` to `10`. Defaults to `5`.

* `code` - (Optional, Int) Specifies the abnormal status code.
  The value can be **404**, **500** or **502**. Defaults to **404**.

* `hosts` - (Optional, String) Specifies the ID of the domain.

* `instances` - (Optional, String) Specifies the ID of the dedicated WAF instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  The default value is **0**.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The abnormal request information.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `key` - The attack type.

* `num` - The attack count.

* `host` - The protected domain.
