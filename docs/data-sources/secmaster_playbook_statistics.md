---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_statistics"
description: |-
  Use this data source to get the list of SecMaster playbook statistics.
---

# huaweicloud_secmaster_playbook_statistics

Use this data source to get the list of SecMaster playbook statistics.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_playbook_statistics" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The playbook statistics.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `enabled_num` - The number of enabled playbooks.

* `unapproved_num` - The number of unapproved playbooks.

* `disabled_num` - The number of playbooks that are not enabled.
