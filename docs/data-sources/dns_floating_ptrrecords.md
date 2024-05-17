---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_floating_ptrrecords"
description: |-
  Use this data source to get the list of DNS PTR records.
---

# huaweicloud_dns_floating_ptrrecords

Use this data source to get the list of DNS PTR records.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_dns_floating_ptrrecords" "test"{
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `record_id` - (Optional, String) Specifies the ID of the PTR record.
  The format is `{region}:{floatingip_id}`, `floatingip_id` indicates the EIP ID corresponding to the PTR record.

* `public_ip` - (Optional, String) Specifies the EIP address of the PTR record.

* `domain_name` - (Optional, String) Specifies the domain name of the PTR record.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID corresponding to the PTR record.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the PTR record.

* `status` - (Optional, String) Specifies the status of the PTR record.
  The valid values are **ACTIVE**, **ERROR**, **FREEZE** and **DISABLE**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ptrrecords` - The list of the PTR records.

  The [ptrrecords](#ptrrecords_struct) structure is documented below.

<a name="ptrrecords_struct"></a>
The `ptrrecords` block supports:

* `id` - The ID of the PTR record.

* `public_ip` - The EIP address corresponding to the PTR record.

* `domain_name` - The domain name of the PTR record.

* `ttl` - The valid cache time of the PTR record (in seconds).

* `enterprise_project_id` - The enterprise project ID corresponding to the PTR record.

* `tags` - The key/value pairs to associate with the PTR record.

* `description` - The description of the PTR record.

* `status` - The current status of the PTR record.
