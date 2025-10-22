---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_domain_names"
description: |-
  Use this data source to get the list of SWR enterprise instance domain names.
---

# huaweicloud_swr_enterprise_domain_names

Use this data source to get the list of SWR enterprise instance domain names.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_domain_names" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `uid` - (Optional, String) Specifies the domain name ID.

* `domain_name` - (Optional, String) Specifies the domain name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domain_name_infos` - Indicates the domain name infos.

  The [domain_name_infos](#domain_name_infos_struct) structure is documented below.

<a name="domain_name_infos_struct"></a>
The `domain_name_infos` block supports:

* `uid` - Indicates the domain name ID.

* `domain_name` - Indicates the domain name.

* `type` - Indicates the domain name type.

* `certificate_id` - Indicates the SCM certificate ID.

* `created_at` - Indicates the create time of the domain name.

* `updated_at` - Indicates the update time of the domain name.
