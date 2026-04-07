---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_firewall_config_quota"
description: |-
  Use this data source to get the CFW firewall configuration quota within HuaweiCloud.
---

# huaweicloud_cfw_firewall_config_quota

Use this data source to get the CFW firewall configuration quota within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "config_type" {}

data "huaweicloud_cfw_firewall_config_quota" "test" {
  fw_instance_id = var.fw_instance_id
  config_type    = var.config_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `config_type` - (Required, String) Specifies the firewall quota type.  
  The valid values are as follows:
  + **ACL**: ACL rule quota.
  + **DNS_DOMAIN_SET**: Network domain set quota.
  + **DOMAIN**: Domain name member quota in a domain set.
  + **DOMAIN_DEVICE**: Domain device quota.
  + **DNS_DOMAIN_SET_PARSE_IP**: Parse IP quota for a network domain set.
  + **APPLICATION_DOMAIN_SET**: Application domain set quota.
  + **APPLICATION_DOMAIN_SET_ITEM**: Domain name member quota in an application domain set.
  + **APPLICATION_DOMAIN_SET_ITEM_DEVICE**: Application domain set device quota.
  + **ADDR_SET**: Address group quota.
  + **ADDR_SET_ITEM**: IP address member quota in an address group.
  + **ADDR_SET_ITEM_DEVICE**: Address group IP address device quota.
  + **SERV_SET**: Service group quota.
  + **SERV_SET_ITEM**: Service member quota in a service group.
  + **SERV_SET_ITEM_DEVICE**: Service group service device quota.
  + **BLACKLIST**: Blacklist quota.
  + **WHITELIST**: Whitelist quota.
  + **PRIVATE_NETWORK_SEGMENT**: Private network segment quota.

* `set_id` - (Optional, String) Specifies the group ID. This parameter is required when querying quotas for IP address
  group members, domain group members, or service group members.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The firewall configuration quota.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `item_info` - The firewall quota information.

  The [item_info](#item_info_struct) structure is documented below.

* `max_quota` - The maximum member quota of the firewall.

* `quota_type` - The firewall member quota type.

* `used_quota` - The used quota.

<a name="item_info_struct"></a>
The `item_info` block supports:

* `max_quota` - The maximum quota.

* `used_quota` - The used quota.

* `extras_info` - The additional parameters, ACL and network domain usage.
