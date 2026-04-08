---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_domain_resolve_ip_list"
description: |-
  Use this data source to get the domain resolve IP list within HuaweiCloud.
---

# huaweicloud_cfw_domain_resolve_ip_list

Use this data source to get the domain resolve IP list within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {} 
variable "domain_address_id" {}

data "huaweicloud_cfw_domain_resolve_ip_list" "test" { 
  fw_instance_id    = var.fw_instance_id 
  domain_address_id = var.domain_address_id 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `domain_address_id` - (Required, String) Specifies the domain address ID.

* `address_type` - (Optional, Int) Specifies the address type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The data of domain resolve IP list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `excess_ip` - The list of excess IPs.

* `parsed_success_ip` - The list of parsed success IPs.
