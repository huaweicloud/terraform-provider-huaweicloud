---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_domain_status"
description: |-
  Use this data source to query the domain running status information.
---

# huaweicloud_waf_domain_status

Use this data source to query the domain running status information.

## Example Usage

```hcl
variable "domain_id" {}

data "huaweicloud_waf_domain_status" "test" {
  host_id = var.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Required, String) Specifies the domain ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - The domain name.

* `status` - The protection status of the domain.
  The value can be **enabled**, **disabled**, or **bypassed**.
  + **enabled**: The WAF protection is enabled. WAF detects attacks based on the policy you configure.
  + **disabled**: The WAF protection is suspended. WAF only forwards requests destined for the domain and does not
  detect attacks.
  + **bypassed**: The WAF protection is bypassed. Requests of the domain are directly sent to the backend server and
  do not pass through WAF.

* `waf_instance_id` - The domain ID.
