---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_recordsets"
description: ""
---

# huaweicloud_dns_recordsets

Use this data source to get the list of DNS recordsets.

## Example Usage

```hcl
variable "zone_id" {}

data "huaweicloud_dns_recordsets" "test" {
  zone_id = var.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `zone_id` - (Required, String) Specifies the zone ID.

* `line_id` - (Optional, String) Specifies the resolution line ID. This parameter is only valid when `zone_id` is a
  public zone ID.

  -> You can use custom line or get more information about default resolution lines
  from [Resolution Lines](https://support.huaweicloud.com/intl/en-us/api-dns/en-us_topic_0085546214.html).

* `tags` - (Optional, String) Specifies the resource tag. The format is as follows: key1,value1|key2,value2.
  Multiple tags are separated by vertical bar (|). The key and value of each tag are separated by comma (,).

* `status` - (Optional, String) Specifies the status of the recordset to be queried. Valid values are as follows:
  + **ACTIVE**: Normal.
  + **ERROR**: Failed.
  + **FREEZE**: Frozen.
  + **DISABLE**: Disabled.
  + **POLICE**: Frozen due to security reasons.
  + **ILLEGAL**: Frozen due to abuse.

* `type` - (Optional, String) Specifies the recordset type.
  + If the `zone_id` is a public zone ID, valid values are **A**, **AAAA**, **MX**, **CNAME**, **TXT**, **NS**, **SRV**
  and **CAA**.
  + If the `zone_id` is a private zone ID, valid values are **A**, **AAAA**, **MX**, **CNAME**, **TXT** and **SRV**.

* `name` - (Optional, String) Specifies the name of the recordset to be queried. Fuzzy matching will work.

* `recordset_id` - (Optional, String) Specifies the ID of the recordset to be queried. Fuzzy matching will work.

* `search_mode` - (Optional, String) Specifies the search mode for `name` and `recordset_id`. Valid values are as follows:
  + **like**: Fuzzy matching.
  + **equal**: Accurate matching.

  If not specified, fuzzy matching will be used.

* `sort_key` - (Optional, String) Specifies the sorting field for the list of the recordsets to be queried.  
  The parameter is left blank by default, indicating that the query results are not sorted.  
  The valid values are as follows:
  + **name**: The name of the recordset.
  + **type**: The type of the recordset.

* `sort_dir` - (Optional, String) Specifies the sorting mode for the list of the recordsets to be queried.  
  The parameter is left blank by default, indicating that the query results are not sorted.  
  The valid values are as follows:
  + **asc**: Ascending order.
  + **desc**: Descending order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `recordsets` - The list of recordsets.
  The [recordsets](#DNSRecordsets_Recordsets) structure is documented below.

<a name="DNSRecordsets_Recordsets"></a>
The `recordsets` block supports:

* `id` - The recordset ID.

* `name` - The recordset name.

* `description` - The recordset description.

* `zone_id` - The zone ID of the recordset.

* `zone_name` - The zone name of the recordset.

* `type` - The recordset type. The value can be **A**, **AAAA**, **MX**, **CNAME**, **TXT**, **NS**, **SRV**, or **CAA**.

* `ttl` - The recordset caching duration (in seconds) on a local DNS server. The longer the duration is, the slower the
  update takes effect.

* `records` - The values of domain name resolution.

* `status` - The recordset status.

* `default` - Whether the record set is created by default. A default record set cannot be deleted.

* `line_id` - The resolution line ID. This attribute is only valid when `zone_id` is a public zone ID.

* `weight` - The weight of the recordset. This attribute is only valid when `zone_id` is a public zone ID.

* `created_at` - The creation time of the recordset, in RFC3339 format.

* `updated_at` - The latest update time of the recordset, in RFC3339 format.
