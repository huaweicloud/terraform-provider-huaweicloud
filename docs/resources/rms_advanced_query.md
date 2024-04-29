---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_advanced_query"
description: ""
---

# huaweicloud_rms_advanced_query

Manages a RMS advanced query resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_rms_advanced_query" "test" {
  name       = "advanced_query_name"
  expression = "select * from table_test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the advanced query name. It contains 1 to 64 characters.

  Changing this parameter will create a new resource.

* `expression` - (Required, String) Specifies the advanced query expression. It contains 1 to 4096 characters.

* `description` - (Optional, String) Specifies the advanced query description. It contains 1 to 512 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The RMS advanced query can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rms_advanced_query.test <id>
```
