---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_reference_table"
description: |-
  Manages a WAF reference table resource within HuaweiCloud.
---

# huaweicloud_waf_reference_table

Manages a WAF reference table resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The reference table resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable "name" {}
variable "enterprise_project_id" {}

resource "huaweicloud_waf_reference_table" "test" {
  name                  = var.name
  type                  = "url"
  enterprise_project_id = var.enterprise_project_id

  conditions = [
    "/admin",
    "/manage"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the WAF reference table resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the reference table. Only letters, digits, hyphens (-),
  underscores(_) and dots(.) are allowed. The maximum length is `64` characters.

* `type` - (Required, String, ForceNew) Specifies the type of the reference table.
  The valid values are **url**, **user-agent**, **ip**, **params**, **cookie**, **referer** and **header**.
  Changing this parameter will create a new resource.

* `conditions` - (Required, List) Specifies the conditions of the reference table.

* `description` - (Optional, String) Specifies the description of the reference table.
  The maximum length is `128` characters.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the reference
  table belongs. For enterprise users, if omitted, default enterprise project will be used.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `creation_time` - The creation time of the reference table, in UTC format.

## Import

There are two ways to import WAF reference table state.

* Using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_reference_table.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_reference_table.test <id>/<enterprise_project_id>
```
