---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_business_partners"
description: |-
  Use this data source to get the list of RDS business partners.
---

# huaweicloud_rds_business_partners

Use this data source to get the list of RDS business partners.

## Example Usage

```hcl
data "huaweicloud_rds_business_partners" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `business_partners` - Indicates the business partner list.

  The [business_partners](#business_partners_struct) structure is documented below.

<a name="business_partners_struct"></a>
The `business_partners` block supports:

* `order` - Indicates the priority, integer value range **1-100**, the smaller the value, the higher the priority.

* `international` - Indicates whether it is an international site service provider.

* `bp_domain_id` - Indicates the service provider ID.

* `bp_name` - Indicates the service provider name.
