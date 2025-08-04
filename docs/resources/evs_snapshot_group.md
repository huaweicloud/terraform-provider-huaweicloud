---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshot_group"
description: |-
  Manages an EVS snapshot group resource within HuaweiCloud.
---

# huaweicloud_evs_snapshot_group

Manages an EVS snapshot group resource within HuaweiCloud.

-> Before using this resource, ensure that there is no snapshot being created under the volume.

## Example Usage

```hcl
variable "name" {}
variable "description" {}
variable "volume_ids" {
  type = list(string)
}
variable "enterprise_project_id" {}
variable "server_id" {}
variable "instant_access" {}
variable "incremental" {}

resource "huaweicloud_evs_snapshot_group" "example" {
  name                  = var.name
  description           = var.description
  volume_ids            = var.volume_ids
  enterprise_project_id = var.enterprise_project_id
  server_id             = var.server_id
  instant_access        = var.instant_access
  incremental           = var.incremental

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `server_id` - (Optional, String, NonUpdatable) Specifies the server ID to which the snapshot group are attached.
  Field server_id and field volume_ids cannot be empty at the same time.

* `volume_ids` - (Optional, List, NonUpdatable) Specifies the list of volume IDs to be included in the snapshot group.
  The value is in the **[id1,id2,...,idx]** format. If `server_id` is set, this parameter can only be used to set
  cloud drives mounted within the specified instance. Setting multiple cloud drive IDs across instances is no longer
  supported.

* `instant_access` - (Optional, Bool, NonUpdatable) Specifies whether to enable instant access for the snapshot group.
  Possible values are **true** (enable) and **false** (disable). Default is **false**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID for the snapshot
  group.

* `incremental` - (Optional, Bool, NonUpdatable) Specifies whether to create an incremental snapshot.
  Default is **false**.

* `tags` - (Optional, Map) Specifies the key/value pairs to be associated with the snapshot group.

* `name` - (Optional, String) Specifies the snapshot group name. The maximum length is `255` bytes.

* `description` - (Optional, String) Specifies the snapshot group description. The maximum length is `255` bytes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The time when the snapshot group was created.

* `status` - The snapshot group status.

* `updated_at` - The time when the snapshot group was updated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

EVS snapshot group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_evs_snapshot_group.test <id>
```

```
Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `volume_ids`, `instant_access`,
`incremental`. It is generally recommended running `terraform plan` after importing the resource. You can then decide
if changes should be applied to the resource, or the resource definition should be updated to align with the snapshot
group. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_evs_snapshot_group" "test" {
    ...
  lifecycle {
    ignore_changes = [
      volume_ids,
      instant_access,
      incremental,
    ]
  }
}
```
