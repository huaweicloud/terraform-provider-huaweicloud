---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_recordset"
description: ""
---

# huaweicloud_dns_recordset

Manages a DNS record set resource within HuaweiCloud.

## Example Usage

### Record Set with Multi-line

```hcl
resource "huaweicloud_dns_zone" "example_zone" {
  name        = "example.com."
  email       = "email2@example.com"
  description = "a zone"
  ttl         = 6000
  zone_type   = "public"
}

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.example_zone.id
  name        = "test.example.com."
  type        = "A"
  description = "a recordset description"
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

### Record Set with Public Zone

```hcl
resource "huaweicloud_dns_zone" "example_zone" {
  name        = "example.com."
  email       = "email2@example.com"
  description = "a public zone"
  ttl         = 6000
  zone_type   = "public"
}

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.example_zone.id
  name        = "rs.example.com."
  description = "An example record set"
  ttl         = 3000
  type        = "A"
  records     = ["10.0.0.1"]
}
```

### Record Set with Private Zone

```hcl
resource "huaweicloud_dns_zone" "example_zone" {
  name        = "example.com."
  email       = "email2@example.com"
  description = "a private zone"
  ttl         = 6000
  zone_type   = "private"
}

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.example_zone.id
  name        = "rs.example.com."
  description = "An example record set"
  ttl         = 3000
  type        = "A"
  records     = ["10.0.0.1"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `zone_id` - (Required, String, ForceNew) Specifies the zone ID.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the record set.
  The name suffixed with a zone name, which is a complete host name ended with a dot.

* `type` - (Required, String) Specifies the type of the record set.
  Value options: **A**, **AAAA**, **MX**, **CNAME**, **TXT**, **NS**, **SRV**, **CAA**.

* `records` - (Required, List) Specifies an array of DNS records. The value rules vary depending on the record set type.

* `ttl` - (Optional, Int) Specifies the time to live (TTL) of the record set (in seconds).
  The valid value is range from `1` to `2,147,483,647`. The default value is `300`.

* `status` - (Optional, String) Specifies the status of the record set.
  Value options: **ENABLE**, **DISABLE**. The default value is **ENABLE**.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the DNS recordset.

* `description` - (Optional, String) Specifies the description of the record set.

* `line_id` - (Optional, String, ForceNew) Specifies the resolution line ID.
  Changing this parameter will create a new resource.

-> **NOTE:** Only public zone support. You can use custom line or get more information about default resolution lines
from [Resolution Lines](https://support.huaweicloud.com/intl/en-us/api-dns/en-us_topic_0085546214.html).

* `weight` - (Optional, Int) Specifies the weight of the record set.
  Only public zone support. The valid value is range from `1` to `1,000`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `zone_name` - The zone name of the record set.

* `zone_type` - The type of zone. The value can be **public** or **private**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The DNS recordset can be imported using `zone_id`, `recordset_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_dns_recordset.test <zone_id>/<recordset_id>
```
