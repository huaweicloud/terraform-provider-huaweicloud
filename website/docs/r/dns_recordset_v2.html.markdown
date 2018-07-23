---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_recordset_v2"
sidebar_current: "docs-huaweicloud-resource-dns-recordset-v2"
description: |-
  Manages a DNS record set in the HuaweiCloud DNS Service
---

# huaweicloud\_dns\_recordset_v2

Manages a DNS record set in the HuaweiCloud DNS Service.

## Example Usage

### Automatically detect the correct network

```hcl
resource "huaweicloud_dns_zone_v2" "example_zone" {
  name = "example.com."
  email = "email2@example.com"
  description = "a zone"
  ttl = 6000
  zone_type = "public"
}

resource "huaweicloud_dns_recordset_v2" "rs_example_com" {
  zone_id = "${huaweicloud_dns_zone_v2.example_zone.id}"
  name = "rs.example.com."
  description = "An example record set"
  ttl = 3000
  type = "A"
  records = ["10.0.0.1"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 DNS client.
    If omitted, the `region` argument of the provider is used.
    Changing this creates a new DNS record set.

* `zone_id` - (Required) The ID of the zone in which to create the record set.
  Changing this creates a new DNS  record set.

* `name` - (Required) The name of the record set. Note the `.` at the end of the name.
  Changing this creates a new DNS record set.

* `type` - (Optional) The type of record set. The options include `A`, `AAAA`, `MX`,
  `CNAME`, `TXT`, `NS`, `SRV`, and `PTR`. Changing this creates a new DNS record set.

* `ttl` - (Optional) The time to live (TTL) of the record set (in seconds). The value
  range is 300–2147483647. The default value is 300.

* `description` - (Optional) A description of the record set.

* `records` - (Required) An array of DNS records.

* `value_specs` - (Optional) Map of additional options. Changing this creates a
  new record set.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `type` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `description` - See Argument Reference above.
* `records` - See Argument Reference above.
* `zone_id` - See Argument Reference above.
* `value_specs` - See Argument Reference above.

## Import

This resource can be imported by specifying the zone ID and recordset ID,
separated by a forward slash.

```
$ terraform import huaweicloud_dns_recordset_v2.recordset_1 <zone_id>/<recordset_id>
```
