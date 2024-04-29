---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_triggers"
description: ""
---

# huaweicloud_swr_image_triggers

Use this data source to get the list of SWR image triggers.

## Example Usage

```hcl
variable "organization" {}
variable "repository" {}

data "huaweicloud_swr_image_triggers" "test" {
  organization = var.organization
  repository   = var.repository
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `organization` - (Required, String) Specifies the name of the organization.

* `repository` - (Required, String) Specifies the name of the repository.

* `name` - (Optional, String) Specifies the name of the image trigger.

* `enabled` - (Optional, String) Specifies whether to enable the trigger.
  The valid values are **true** and **false**.
  
* `condition_type` - (Optional, String) Specifies the trigger condition type.
  The valid values are **all**, **tag**, **regular**.

* `cluster_name` - (Optional, String) Specifies the name of the triggered cluster in CCE.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `triggers` - All triggers that match the filter parameters.
  The [triggers](#attrblock--triggers) structure is documented below.

<a name="attrblock--triggers"></a>
The `triggers` block supports:

* `action` - The trigger action.

* `workload_name` - The name of the application.

* `workload_type` - The type of the application.

* `cluster_id` - The ID of the cluster in CCE.

* `cluster_name` - The name of the cluster in CCE.

* `namespace` - The namespace where the application is located.

* `name` - The trigger name.

* `type` - The trigger type.

* `condition_type` - The trigger condition type.

* `condition_value` - The trigger condition value.

* `container` - The name of the container to be updated.

* `enabled` - Whether to enable the trigger.

* `created_at` - The creation time of the trigger.

* `created_by` - The creator name of the trigger.

* `histories` - All histories of the trigger.
  The [histories](#attrblock--triggers--histories) structure is documented below.

<a name="attrblock--triggers--histories"></a>
The `histories` block supports:

* `result` - The triggered result.

* `tag` - The triggered image tag.

* `detail` - The triggered detail.
