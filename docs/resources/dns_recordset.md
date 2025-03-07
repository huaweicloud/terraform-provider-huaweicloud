---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_recordset"
description: |-
  Manages a DNS record set resource within HuaweiCloud.
---

# huaweicloud_dns_recordset

Manages a DNS record set resource within HuaweiCloud.

## Example Usage

### Record Set with Public Zone

```hcl
variable "public_zone_id" {}
variable "public_recordset_name" {}
variable "description" {}

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = var.public_zone_id
  name        = var.public_recordset_name
  type        = "A"
  description = var.description
  status      = "ENABLE"
  ttl         = 300
  records     = ["10.1.0.0"]
  line_id     = "Dianxin_Shanxi"
  weight      = 3

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

### Record Set with Private Zone

```hcl
variable "private_zone_id" {}
variable "private_recordset_name" {}
variable "description" {}

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = var.private_zone_id
  name        = var.private_recordset_name
  description = var.description
  ttl         = 3000
  type        = "A"
  records     = ["10.0.0.1"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `zone_id` - (Required, String, ForceNew) Specifies the ID of the zone to which the record set belongs.  
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the record set.  
  The name suffixed with a zone name, which is a complete host name ended with a dot.

* `type` - (Required, String) Specifies the type of the record set.  
  + For the public record set, the valid values are **A**, **AAAA**, **MX**, **CNAME**, **TXT**, **NS**, **SRV** and **CAA**.
  + For the private record set, the valid values are **A**, **AAAA**, **MX**, **CNAME**, **TXT** and **SRV**.

* `records` - (Required, List) Specifies the list of the records of the record set.  
  The value depends on the `type` parameter, you can refer to this [document](https://support.huaweicloud.com/intl/en-us/usermanual-dns/dns_usermanual_0601.html#dns_usermanual_0601__table936244914119).

* `ttl` - (Optional, Int) Specifies the time to live (TTL) of the record set (in seconds).  
  The valid value is range from `1` to `2,147,483,647`. The default value is `300`.

* `line_id` - (Optional, String, ForceNew) Specifies the resolution line ID.  
  Changing this parameter will create a new resource.

-> Only public zone support. You can use custom line or get more information about default resolution lines
   from [Resolution Lines](https://support.huaweicloud.com/intl/en-us/api-dns/en-us_topic_0085546214.html).

* `status` - (Optional, String) Specifies the status of the record set, defaults to **ENABLE**.  
  The valid values are as follows:
  + **ENABLE**
  + **DISABLE**

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the DNS record set.

* `description` - (Optional, String) Specifies the description of the record set.

* `weight` - (Optional, Int) Specifies the weight of the record set.
  Only public zone support. The valid value is range from `1` to `1,000`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consists of the `zone_id` and the record set ID, separated by a slash.

* `zone_name` - The name of the zone to which the record set belongs.

* `zone_type` - The type of the zone to which the record set belongs.
  + **public**
  + **private**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The DNS recordset can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_dns_recordset.test <id>
```
