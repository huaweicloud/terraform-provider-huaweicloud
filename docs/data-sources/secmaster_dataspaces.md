---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_dataspaces"
description: |-
  Use this data source to get the list of dataspaces.
---

# huaweicloud_secmaster_dataspaces

Use this data source to get the list of dataspaces.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_dataspaces" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `dataspace_id` - (Optional, String) Specifies the dataspace ID.

* `dataspace_name` - (Optional, String) Specifies the dataspace name.

* `sort_key` - (Optional, String) Specifies sorting field.

* `sort_dir` - (Optional, String) Specifies sorting order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of dataspaces.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `data` block supports:

* `dataspace_id` - The dataspace ID.

* `dataspace_name` - The dataspace name.

* `dataspace_type` - The dataspace type.

* `description` - The dataspace description.

* `domain_id` - The account ID.

* `project_id` - The project ID.

* `region_id` - The region ID.

* `create_by` - The dataspace creator.

* `update_by` - The dataspace updater.

* `create_time` - The creation time.

* `update_time` - The update time.
