---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_stored_query"
description: |-
  Manages a stored ResourceQL query within HuaweiCloud Resource Center.
---

# huaweicloud_rc_stored_query

Manages a stored ResourceQL query within HuaweiCloud Resource Center.

## Example Usage

```hcl
variable "query_name" {}
variable "query_expression" {}
variable "query_description" {}

resource "huaweicloud_rc_stored_query" "test" {
  name        = var.query_name
  expression  = var.query_expression
  description = var.query_description
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the ResourceQL name.

* `expression` - (Required, String) Specifies the ResourceQL expression.

* `description` - (Optional, String) Specifies the ResourceQL description.

## Attribute Reference

* `id` - The unique ID of the operation.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `deleted_at` - The deletion time.

* `expired_at` - The expiration time.

* `created` - The ResourceQL created time.

* `updated` - The ResourceQL updated time.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rc_stored_query.test <id>
```
