---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_reference_tables"
description: |-
  Use this data source to get a list of WAF reference tables.
---

# huaweicloud_waf_reference_tables

Use this data source to get a list of WAF reference tables.

## Example Usage

```hcl
variable "enterprise_project_id" {}

data "huaweicloud_waf_reference_tables" "test" {
  name                  = "reference_table_name"
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the WAF reference table resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the reference table. The value is case-sensitive and matches exactly.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of WAF reference tables.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tables` - A list of WAF reference tables.

The `tables` block supports:

* `id` - The ID of the reference table.

* `name` - The name of the reference table. The maximum length is `64` characters.

* `type` - The type of the reference table, The options are: `url`, `user-agent`, `ip`, `params`, `cookie`, `referer`
  and `header`.

* `conditions` - The conditions of the reference table.

* `description` - The description of the reference table.

* `creation_time` - The server time when reference table was created.
