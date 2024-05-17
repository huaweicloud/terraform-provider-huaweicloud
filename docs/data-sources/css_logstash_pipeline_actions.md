---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_pipeline_actions"
description: |-
  Use this data source to get the list of CSS logstash pipeline actions.
---

# huaweicloud_css_logstash_pipeline_actions

Use this data source to get the list of CSS logstash pipeline actions.

## Example Usage

```hcl
variable "cluster_id" {}
variable "action_id" {}

data "huaweicloud_css_logstash_pipeline_actions" "test" {
  cluster_id = var.cluster_id
  action_id  = var.action_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the CSS logstash cluster.

* `action_id` - (Optional, String) Specifies the ID of the action.

* `type` - (Optional, String) Specifies the type of the action.
  The values can be **start**, **hotStart**, **hotStop** and **stop**.

* `status` - (Optional, String) Specifies the status of the action.
  The values can be **running**, **success** and **failed**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `actions` - The list of the actions.

  The [actions](#actions_struct) structure is documented below.

<a name="actions_struct"></a>
The `actions` block supports:

* `id` - The ID of the action.

* `type` - The type of the action.

* `status` - The status of the action.

* `conf_content` - The configuration file content.

* `updated_at` - The update time.

* `error_msg` - The error message of the action.

* `message` - The message of the action.
