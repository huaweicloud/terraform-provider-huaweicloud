---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_dedicated_domains"
description: |-
  Use this data source to get a list of WAF dedicated domains.
---

# huaweicloud_waf_dedicated_domains

Use this data source to get a list of WAF dedicated domains.

## Example Usage

```hcl
variable "domain" {}
variable "enterprise_project_id" {}

data "huaweicloud_waf_dedicated_domains" "test" {
  domain                = var.domain
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the WAF dedicated domains.
  If omitted, the provider-level region will be used.

* `domain` - (Optional, String) Specifies the protected domain name or IP address (port allowed).

* `policy_name` - (Optional, String) Specifies the policy name associated with the domain.

* `protect_status` - (Optional, String) Specifies the protection status of domain. Valid values are:
  + `0`: The WAF protection is suspended. WAF only forwards requests for the domain name but does not detect attacks.
  + `1`: The WAF protection is enabled. WAF detects attacks based on the policy you configure.

  If omitted, all domains in different protection status will be queried.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domains` - A list of WAF dedicated domains.

The `domains` block supports:

* `id` - The ID of WAF dedicated domain.

* `domain` - The protected domain name or IP address (port allowed).

* `pci_3ds` - The status of the PCI 3DS compliance certification check.

* `pci_dds` - The status of the PCI DSS compliance certification check.

* `is_dual_az` - The status of the WAF support dual AZ mode.

* `ipv6_enable` - Whether to support IPv6.

* `description` - The description of WAF dedicated domain.

* `policy_id` - The policy ID associated with the domain.

* `protect_status` - The protection status of domain, `0`: suspended, `1`: enabled.

* `access_status` - Whether a domain name is connected to WAF. Valid values are:
  + `0` - The domain name is not connected to WAF,
  + `1` - The domain name is connected to WAF.

* `website_name` - The website name.
