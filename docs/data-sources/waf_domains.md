---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_domains"
description: |-
  Use this data source to get a list of WAF domains.
---

# huaweicloud_waf_domains

Use this data source to get a list of WAF domains.

## Example Usage

```hcl
variable "domain" {}
variable "enterprise_project_id" {}

data "huaweicloud_waf_domains" "test" {
  domain                = var.domain
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the WAF domains.
  If omitted, the provider-level region will be used.

* `domain` - (Optional, String) Specifies the protected domain name or IP address (port allowed).

* `policy_name` - (Optional, String) Specifies the policy name associated with the domain.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domains` - A list of WAF domains.

The `domains` block supports:

* `id` - The ID of WAF domain.

* `description` - The description of WAF domain.

* `proxy` - Whether the protected domain name uses a proxy.
  Valid values are **true** and **false**.

* `domain` - The protected domain name or IP address (port allowed).

* `policy_id` - The policy ID associated with the domain.

* `pci_3ds` - The status of the PCI 3DS compliance certification check.

* `pci_dss` - The status of the PCI DSS compliance certification check.

* `ipv6_enable` - Whether to support IPv6.

* `protect_status` - The protection status of domain, `0`: suspended, `1`: enabled.

* `access_status` - Whether a domain name is connected to WAF. Valid values are:
  + `0` - The domain name is not connected to WAF.
  + `1` - The domain name is connected to WAF.

* `access_code` - The CNAME prefix. The CNAME suffix is `.vip1.huaweicloudwaf.com`.

* `charging_mode` - The charging mode of the domain.
  Valid values are **prePaid** and **postPaid**.

* `website_name` - The website name.

* `proxy_layer` - Type of front-end proxy. Valid values are:
  + `0`: No proxy. No proxy products are deployed in front of WAF.
  + `4`: Layer-4 proxy. Web proxy products, such as layer-4 anti-DDoS,
  that will not change the source or destination IP addresses are deployed in front of WAF.
  + `7`: Layer-7 proxy. Web proxy products, such as layer-7 anti-DDoS, CDN,
  and other cloud acceleration services, that will change the source and
  destination IP addresses are deployed in front of WAF.
  If a layer-7 proxy is configured, WAF reads the real client IP address
  from the related fields in the header.

* `created_at` - The creation time of domain.
