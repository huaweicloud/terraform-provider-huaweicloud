---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_quotas"
description: |-
  Use this data source to get the list of DNS resource quotas within HuaweiCloud.
---

# huaweicloud_dns_quotas

Use this data source to get the list of DNS resource quotas within HuaweiCloud.

## Example Usage

```hcl
variable domain_id {}

data "huaweicloud_dns_quotas" "test" {
  domain_id = var.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, String) Specifies the account ID of IAM user.

* `type` - (Optional, String) Specifies the resource type.  
  The valid values are as follows:
  + **zone**
  + **private_zone**
  + **record_set**
  + **ptr_record**
  + **custom_line**
  + **line_group**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - All quotas that match the filter parameters.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `type` - The resource type corresponding to quota.

* `max` - The maximum quota of resource.

* `used` - The used quota of resource.

* `unit` - The unit of the quota.
