---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_domains"
description: |-
  Use this data source to get the list of LIVE domain names within HuaweiCloud.
---

# huaweicloud_live_domains

Use this data source to get the list of LIVE domain names within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "enterprise_project_id" {}

data "huaweicloud_live_domains" "test" {
  name                  = var.domain_name
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the domain name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If omitted, all domain names will be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domains` - The domain name list.

  The [domains](#domains_struct) structure is documented below.

<a name="domains_struct"></a>
The `domains` block supports:

* `vendor` - The CDN vendor.

* `region` - The region to which the live broadcast source station belongs.

* `ingest_domain_name` - The ingest domain name associated with the streaming domain name.

* `created_at` - The time when the domain name was created.

* `status_describe` - The status description.

* `service_area` - The domain name acceleration region. Valid values are:
  + **mainland_china**: Chinese mainland.
  + **outside_mainland_china**: Outside the Chinese mainland.
  + **global**: Global acceleration.

* `name` - The domain name.

* `type` - The domain name type.

* `is_ipv6` - Whether the IPv6 function is enabled.
  + **true**: Indicates that IPv6 is enabled.
  + **false**: Indicates that IPv6 is disabled

* `enterprise_project_id` - The enterprise project ID.

* `cname` - The CNAME of the domain name.

* `status` - The status of the domain name.
