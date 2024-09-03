---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_data_classes"
description: |-
  Use this data source to get the list of SecMaster data classes.
---

# huaweicloud_secmaster_data_classes

Use this data source to get the list of SecMaster data classes.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `name` - (Optional, String) Specifies the name of the data class. Fuzzy matching is supported.

* `business_code` - (Optional, String) Specifies the business code of the data class. Fuzzy matching is supported.

* `description` - (Optional, String) Specifies the data class description. Fuzzy matching is supported.

* `is_built_in` - (Optional, String) Specifies whether the data class is built in SecMaster.
  The value can be  **true** or **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_classes` - The data class list.

  The [data_classes](#data_classes_struct) structure is documented below.

<a name="data_classes_struct"></a>
The `data_classes` block supports:

* `id` - The ID of the data class.

* `name` - The name of the data class.

* `description` - The description of the data class.

* `business_code` - The business code of the data class.

* `is_built_in` - Whether the data class is built in SecMaster.

* `type_num` - The quantity of sub-type data classes.

* `subscribed_version` - The subscribed version of the data class.

* `parent_id` - The parent data class ID.

* `workspace_id` - The workspace ID.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `creator_id` - The creator ID.

* `modifier_id` - The modifier ID.

* `creator` - The creator.

* `modifier` - The modifier.

* `domain_id` - The domain ID.

* `region_id` - The region ID.

* `project_id` - The project ID.
