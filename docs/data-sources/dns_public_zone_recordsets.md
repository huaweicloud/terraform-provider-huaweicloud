---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_public_zone_recordsets"
description: |-
  Use this data source to get the list of email or website recordsets of the public zone within HuaweiCloud.
---

# huaweicloud_dns_public_zone_recordsets

Use this data source to get the list of email or website recordsets of the public zone within HuaweiCloud.

## Example Usage

### Query email recordsets

```hcl
variable "zone_id" {}

data "huaweicloud_dns_public_zone_recordsets" "test" {
  zone_id = var.zone_id
  type    = "email"
}
```

### Query website recordsets

```hcl
data "huaweicloud_dns_public_zone_recordsets" "test" {
  zone_id = var.zone_id
  type    = "website"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Specifies the ID of the public zone to be queried.

* `type` - (Required, String) Specifies the type of the domain name to be queried.  
  The valid values are as follows:
  + **email**
  + **website**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `recordsets` - The list of recordsets.  
  The [recordsets](#public_zone_recordsets) structure is documented below.

<a name="public_zone_recordsets"></a>
The `recordsets` block supports:

* `id` - The ID of the recordset.

* `name` - The name of the recordset.

* `zone_id` - The ID of the zone to which the recordset belongs.

* `type` - The type of the recordset.

* `default` - Whether the recordset is default.

* `created_at` - The creation time of the recordset, in RFC3339 format.

* `updated_at` - The update time of the recordset, in RFC3339 format.
