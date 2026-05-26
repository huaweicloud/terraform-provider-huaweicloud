---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_data_secrecy_levels"
description: |-
  Use this data source to get the list of DataArts Security data secrecy levels within HuaweiCloud.
---

# huaweicloud_dataarts_security_data_secrecy_levels

Use this data source to get the list of DataArts Security data secrecy levels within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_data_secrecy_levels" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the data secrecy levels are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the data secrecy levels belong.

* `order_by` - (Optional, String) Specifies the field used to sort the data secrecy levels.  
  The valid values are as follows:
  + **name**
  + **description**
  + **createdAt**
  + **updatedAt**
  + **createdBy**
  + **updatedBy**

* `desc` - (Optional, Bool) Specifies whether to sort the data secrecy levels in descending order.
  Default is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `secrecy_levels` - The list of data secrecy levels that matched the filter parameters.  
  The [secrecy_levels](#dataarts_security_data_secrecy_levels_attr) structure is documented below.

<a name="dataarts_security_data_secrecy_levels_attr"></a>
The `secrecy_levels` block supports:

* `id` - The ID of the data secrecy level.

* `name` - The name of the data secrecy level.

* `level_number` - The level number of the data secrecy level.

* `description` - The description of the data secrecy level.

* `instance_id` - The instance ID to which the data secrecy level belongs.

* `created_by` - The creator of the data secrecy level.

* `created_at` - The creation time of the data secrecy level, in RFC3339 format.

* `updated_by` - The updater of the data secrecy level.

* `updated_at` - The latest update time of the data secrecy level, in RFC3339 format.
