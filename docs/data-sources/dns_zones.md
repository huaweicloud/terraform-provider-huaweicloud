---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_zones"
description: ""
---

# huaweicloud_dns_zones

Use this data source to get the list of DNS zones.

## Example Usage

```hcl
variable "zone_type" {}
variable "enterprise_project_id" {}

data "huaweicloud_dns_zones" "test" {
  zone_type             = var.zone_type
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `zone_type` - (Required, String) Specifies the zone type. The value can be **public** or **private**.

* `zone_id` - (Optional, String) Specifies the ID of the zone.

* `tags` - (Optional, String) Specifies the resource tag. The format is as follows: key1,value1|key2,value2.
  Multiple tags are separated by vertical bar (|). The key and value of each tag are separated by comma (,). The tags
  are in `AND` relationship. Exact matching will work. If the value starts with an asterisk (*), fuzzy matching will
  work for the string following the asterisk.

* `status` - (Optional, String) Specifies the zone status. Valid values are as follows:
  + **ACTIVE**: Normal.
  + **ERROR**: Failed.
  + **FREEZE**: Frozen.
  + **DISABLE**: Disabled.
  + **POLICE**: Frozen due to security reasons.
  + **ILLEGAL**: Frozen due to abuse.

* `name` - (Optional, String) Specifies the name of the zone to be queried. Fuzzy matching will work.

* `search_mode` - (Optional, String) Specifies the search mode for `name`. Valid values are as follows:
  + **like**: Fuzzy matching.
  + **equal**: Accurate matching.
  
  If not specified, fuzzy matching will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID which the zone associated.

* `router_id` - (Optional, String) Specifies the ID of the VPC associated with the private zone.  
  This parameter is available only when the `zone_type` parameter is set to **private**.

* `sort_key` - (Optional, String) Specifies the sorting filed for the list of the zones.  
  The valid values are as follows:
  + **name**: The zone name.
  + **created_at**: The creation time of the zone.
  + **updated_at** The update time of the zone.

* `sort_dir` - (Optional, String) Specifies the sorting mode for the list of the zones.  
  The valid values are as follows:
  + **DESC**: Descending order.
  + **ASC**: Ascending order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `zones` - The list of zones.
  The [zones](#DNSZones_Zones) structure is documented below.

<a name="DNSZones_Zones"></a>
The `zones` block supports:

* `id` - The zone ID.

* `name` - The zone name.

* `description` - The zone description.

* `email` - The email address of the administrator managing the zone.

* `zone_type` - The zone type. Valid values are **public** and **private**.

* `ttl` - The time to live (TTL) of the zone. The unit is seconds.

* `enterprise_project_id` - The enterprise project ID.

* `status` - The zone status.

* `record_num` - The number of record sets in the zone.

* `masters` - The master DNS servers, from which the slave servers get DNS information.

* `tags` - The key/value pairs to associate with the zone.

* `routers` - The list of VPCs associated with the zone. This attribute is only valid when `zone_type` is **private**.
  The [routers](#Zones_routers) structure is documented below.

* `proxy_pattern` - The recursive resolution proxy mode for subdomains of the private zone.
  + **AUTHORITY**: The recursive resolution proxy is disabled for the private zone.
  + **RECURSIVE**: The recursive resolution proxy is enabled for the private zone.

* `pool_id` - The ID of the pool to which the zone belongs, assigned by the system.

* `created_at` - The creation time of the zone, in RFC3339 format.

* `updated_at` - The latest update time of the zone, in RFC3339 format.

<a name="Zones_routers"></a>
The `routers` block supports:

* `router_id` - The ID of the VPC associated with the zone.

* `router_region` - The region of the VPC.
