---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_dns_domains"
description: |-
  Use this data source to get the list of DNS domains.
---

# huaweicloud_waf_dns_domains

Use this data source to get the list of DNS domains.

## Example Usage

```hcl
variable "enterprise_project_id" {}

data "huaweicloud_waf_dns_domains" "test" {
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the DNS domain belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The list of DNS domains.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID of the DNS domain.

* `domain` - The domain name of the DNS domain.

* `servers` - The list of servers.

  The [servers](#servers_struct) structure is documented below.

* `protect_port` - The protected port of the DNS domain.

<a name="servers_struct"></a>
The `servers` block supports:

* `type` - The type of the server.

* `address` - The address of the server.
