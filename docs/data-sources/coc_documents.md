---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_documents"
description: |-
  Use this data source to get the list of COC documents.
---

# huaweicloud_coc_documents

Use this data source to get the list of COC documents.

## Example Usage

```hcl
data "huaweicloud_coc_documents" "test" {}
```

## Argument Reference

The following arguments are supported:

* `name_like` - (Optional, String) Specifies the document name and support fuzzy query.

* `creator` - (Optional, String) Specifies the creator.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `document_type` - (Optional, String) Specifies the type of document being executed.
  The value can be **PUBLIC** or **NORMAL**. The default value is **NORMAL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the documents list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `document_id` - Indicates the document ID.

* `name` - Indicates the document name.

* `create_time` - Indicates the creation time.

* `update_time` - Indicates the update time.

* `version` - Indicates the document version, such as **v1**.

* `creator` - Indicates the creator.

* `modifier` - Indicates the modifier.

* `enterprise_project_id` - Indicates the enterprise project ID.
