---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_checkitem"
description: |-
  Manages a customized checkitem resource within HuaweiCloud.
---

# huaweicloud_secmaster_checkitem

Manages a customized checkitem resource within HuaweiCloud.

## Example Usage

### Create a manual check item

```hcl
variable "workspace_id" {}

resource "huaweicloud_secmaster_checkitem" "test" {
  workspace_id = var.workspace_id
  name         = "test_checkitem"
  description  = "test checkitem description"
  level        = "medium"
  cloud_server = "IAM"
  method       = 0
  source       = 3
}
```

### Create an automatic-playbook check item

```hcl
variable "workspace_id" {}
variable "workflow_id" {}

resource "huaweicloud_secmaster_checkitem" "test" {
  workspace_id = var.workspace_id
  name         = "test_checkitem"
  description  = "test checkitem description"
  level        = "medium"
  cloud_server = "IAM"
  method       = 0
  source       = 2
  workflow_id  = var.workflow_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the checkitem belongs.

* `name` - (Required, String) Specifies the name of the checkitem.

* `description` - (Required, String) Specifies the description of the checkitem.

* `level` - (Required, String) Specifies the severity level of the checkitem.
  Valid values: **informational**, **low**, **medium**, **high**, **fatal**.

* `cloud_server` - (Required, String) Specifies the cloud service to which the checkitem belongs.

* `method` - (Required, Int) Specifies the check method of the checkitem.
  Valid values: **0** (manual), **1** (automatic), **3** (automatic-playbook),
  **4** (automatic-HSS), **5** (automatic-CSS).

* `audit_procedure` - (Optional, String) Specifies the audit procedure of the checkitem.

* `aggregation_handle_status` - (Optional, String) Specifies the aggregation handle status of the checkitem.

* `impact` - (Optional, String) Specifies the impact of the checkitem.

* `source` - (Optional, Int) Specifies the source of the checkitem.
  Valid values: **0** (automatic), **2** (automatic-playbook), **3** (manual),
  **4** (automatic-HSS), **5** (automatic-CSS).

* `workflow_id` - (Optional, String) Specifies the workflow ID of the checkitem.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the checkitem name.

* `uuid` - The UUID of the checkitem.

## Import

The checkitem can be imported using the `workspace_id` and the `name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_checkitem.test <workspace_id>/<checkitem_name>
```
