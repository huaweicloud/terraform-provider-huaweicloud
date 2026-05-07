---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_workspace_quotas"
description: |-
  Use this data source to query the quotas of a specified ModelArts workspace within Huaweicloud.
---

# huaweicloud_modelarts_workspace_quotas

Use this data source to query the quotas of a specified ModelArts workspace within Huaweicloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_modelarts_workspace_quotas" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the workspace quotas are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of workspace quotas that matched filter parameters.  
  The [quotas](#workspace_quotas) structure is documented below.

<a name="workspace_quotas"></a>
The `quotas` block supports:

* `resource` - The unique identifier of the resource.

* `quota` - The current quota value. If the value is `-1`, it means there is no quota limit.

* `min_quota` - The minimum value allowed for the quota.

* `max_quota` - The maximum value allowed for the quota.

* `used_quota` - The used quota value.

* `name_cn` - The name of the quota in Chinese.

* `name_en` - The name of the quota in English.

* `unit_cn` - The unit of the quota in Chinese.

* `unit_en` - The unit of the quota in English.

* `updated_at` - The last update time of the quota, in RFC3339 format.
