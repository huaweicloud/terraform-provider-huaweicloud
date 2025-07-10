---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_domains"
description: |-
  Use this data source to get the list of Advanced Anti-DDos protected domains within HuaweiCloud.
---

# huaweicloud_aad_domains

Use this data source to get the list of Advanced Anti-DDos protected domains within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_domains" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The list of domains.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `domain_name` - The domain name.

* `cname` - The domain cname.

* `protocol` - The domain protocol.

* `real_server_type` - The type of real server.

* `real_servers` - The real servers.

* `waf_status` - The WAF status.

* `enterprise_project_id` - The enterprise project ID.

* `domain_id` - The domain ID.
