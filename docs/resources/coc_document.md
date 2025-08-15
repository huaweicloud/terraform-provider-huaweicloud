---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_document"
description: |-
  Manages a COC document resource within HuaweiCloud.
---

# huaweicloud_coc_document

Manages a COC document resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "content" {}

resource "huaweicloud_coc_document" "test" {
  name                  = var.name
  content               = var.content
  risk_level            = "HIGH"
  enterprise_project_id = "0"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the document name.

* `content` - (Required, String) Specifies the document content, it is DSL statements.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.

* `risk_level` - (Required, String) Specifies the risk level.
  The value can be **LOW**, **MEDIUM** or **HIGH**.

* `description` - (Optional, String) Specifies the document description.

* `tags` - (Optional, Map, NonUpdatable) Specifies the key/value pairs to associate with the document.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - Indicates the creation time.

* `update_time` - Indicates the update time.

* `version` - Indicates the document version, such as **v1**.

* `creator` - Indicates the creator.

* `modifier` - Indicates the modifier.

* `versions` - Indicates the collection of versions.

  The [versions](#versions_struct) structure is documented below.

<a name="versions_struct"></a>
The `versions` block supports:

* `version` - Indicates the version number, such as **v1**.

* `version_uuid` - Indicates the version ID.

* `create_time` - Indicates the version creation time.

## Import

The COC document can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_document.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `tags`.
It is generally recommended running `terraform plan` after importing a document.
You can then decide if changes should be applied to the document, or the resource definition should be updated to
align with the document. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_document" "test" {
    ...

  lifecycle {
    ignore_changes = [
      tags,
    ]
  }
}
```
