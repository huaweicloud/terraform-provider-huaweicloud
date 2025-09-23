---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_domain_name_parse_ip_list"
description: |-
  Use this data source to get the DNS resolution result of a domain name.
---

# huaweicloud_cfw_domain_name_parse_ip_list

Use this data source to get the DNS resolution result of a domain name.

## Example Usage

### DNS resolution result of a domain name

```hcl
data "huaweicloud_cfw_domain_name_parse_ip_list" "test" {
  domain_name = "www.baidu.com"
}
```

### DNS resolution result of a domain name in a domain name group

```hcl
variable "fw_instance_id" {}
variable "domain_name_group_id" {}
variable "domain_address_id" {}

data "huaweicloud_cfw_domain_name_parse_ip_list" "test" {
  fw_instance_id    = var.fw_instance_id
  group_id          = var.domain_name_group_id
  domain_address_id = var.domain_address_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_name` - (Optional, String) Specifies the domain name.

* `fw_instance_id` - (Optional, String) Specifies the firewall ID.

* `group_id` - (Optional, String) Specifies the domain name group ID.

* `domain_address_id` - (Optional, String) Specifies the domain name ID in a domain name group.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `address_type` - (Optional, String) Specifies the address type.
  The valid value can be **0** (IPv4) or **1** (IPv6).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The IP address list for domain name resolution.
