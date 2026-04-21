---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_data_recognition_rule_group"
description: |-
  Use this resource to manage a data recognition rule group for DataArts Studio Security within HuaweiCloud.
---

# huaweicloud_dataarts_security_data_recognition_rule_group

Use this resource to manage a data recognition rule group for DataArts Studio Security within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "group_name" {}
variable "data_recognition_rule_ids" {
  type = list(string)
}

resource "huaweicloud_dataarts_security_data_recognition_rule_group" "test" {
  workspace_id = var.workspace_id
  name         = var.group_name
  rule_ids     = var.data_recognition_rule_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the data recognition rule group is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the rule group belongs.

* `name` - (Required, String) Specifies the name of the data recognition rule group.

* `rule_ids` - (Required, List) Specifies the list of data recognition rule IDs that the rule group contains.

* `description` - (Optional, String) Specifies the description of the data recognition rule group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_by` - The creator of the data recognition rule group.

* `created_at` - The creation time of the data recognition rule group, in RFC3339 format.

* `updated_by` - The updater of the data recognition rule group.

* `updated_at` - The update time of the data recognition rule group, in RFC3339 format.

## Import

The resource can be imported using `workspace_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_security_data_recognition_rule_group.test <workspace_id>/<id>
```
