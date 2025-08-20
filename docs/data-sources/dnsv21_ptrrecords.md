---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dnsv21_ptrrecords"
description: |-
  Use this data source to get the list of DNS PTR records.
---

# huaweicloud_dnsv21_ptrrecords

Use this data source to get the list of DNS PTR records.

## Example Usage

```hcl
data "huaweicloud_dnsv21_ptrrecords" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) The enterprise project ID corresponding to the PTR record.

* `resource_type` - (Optional, String) The resource type.

* `status` - (Optional, String) The status of the PTR record.

* `tags` - (Optional, Map) The key/value pairs to associate with the PTR record.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ptrrecords` - Indicates the PTR records list.
  The [ptrrecords](#attrblock--ptrrecords) structure is documented below.

<a name="attrblock--ptrrecords"></a>
The `ptrrecords` block supports:

* `address` - The address of the EIP.

* `description` - The description of the PTR record.

* `enterprise_project_id` - The enterprise project ID of the PTR record.

* `id` - The ID of the PTR record.

* `names` - The domain names of the PTR record.

* `publicip_id` - The ID of the EIP.

* `status` - The status of the PTR record.

* `tags` - The key/value pairs to associate with the PTR record.

* `ttl` - The time to live (TTL) of the record set (in seconds).
