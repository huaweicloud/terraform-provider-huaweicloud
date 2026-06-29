---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_siem_shipper"
description: |-
  Manages a SIEM shipper resource within HuaweiCloud.
---

# huaweicloud_secmaster_siem_shipper

Manages a SIEM shipper resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "shipper_name" {}

resource "huaweicloud_secmaster_siem_shipper" "test" {
  consumption_type = 0
  dataspace_id     = "xxx"
  dataspace_name   = "xxx"
  domain_id        = "xxx"
  pipe_id          = "xxx"
  pipe_name        = "xxx"
  project_id       = "xxx"
  region           = "cn-north-4"
  shipper_name     = var.shipper_name
  version          = "v1"
  workspace_id     = var.workspace_id
  workspace_name   = "xxx"

  shipper_destination {
    data_param                 = jsonencode({})
    destination_dataspace      = "xxx"
    destination_dataspace_name = "xxx"
    destination_identity_role  = "pipe-strategy"
    destination_pipe           = "xxx"
    destination_pipe_name      = "xxx"
    destination_region         = "cn-north-4"
    destination_shipper_type   = 0
    destination_workspace      = "xxx"
    destination_workspace_name = "xxx"
  }

  shipper_source {
    region                = "cn-north-4"
    source_dataspace      = "xxx"
    source_dataspace_name = "xxx"
    source_identity_role  = "pipe-strategy"
    source_pipe           = "xxx"
    source_pipe_name      = "xxx"
    source_type           = 0
    source_workspace      = "xxx"
    source_workspace_name = "xxx"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the SIEM shipper belongs.

* `shipper_name` - (Required, String, NonUpdatable) Specifies the name of the SIEM shipper.

* `consumption_type` - (Optional, Int, NonUpdatable) Specifies the consumption type.
  **0** means real-time, **1** means scheduled.

* `dataspace_id` - (Optional, String, NonUpdatable) Specifies the data space ID.

* `dataspace_name` - (Optional, String, NonUpdatable) Specifies the data space name.

* `domain_id` - (Optional, String, NonUpdatable) Specifies the domain ID.

* `pipe_id` - (Optional, String, NonUpdatable) Specifies the pipe ID.

* `pipe_name` - (Optional, String, NonUpdatable) Specifies the pipe name.

* `project_id` - (Optional, String, NonUpdatable) Specifies the project ID.

* `shipper_destination` - (Optional, List, NonUpdatable) Specifies the shipper destination configuration.
  The [shipper_destination](#shipper_destination_block) structure is documented below.

* `shipper_source` - (Optional, List, NonUpdatable) Specifies the shipper source configuration.
  The [shipper_source](#shipper_source_block) structure is documented below.

* `version` - (Optional, String, NonUpdatable) Specifies the version information.

* `workspace_name` - (Optional, String, NonUpdatable) Specifies the workspace name.

* `action` - (Optional, String) Specifies the operation to perform on the SIEM shipper. Valid values are:
  + **pause**: Permitted only when the delivery status is "in delivery".
  + **resume**: Permitted only when the delivery status is "Pending".
  + **retry**: Resubmission is only possible when the authorization status is "Cancelled".

<a name="shipper_destination_block"></a>
The `shipper_destination` block supports:

* `data_param` - (Optional, String, NonUpdatable) Specifies the data parameter, usually in JSON format.

* `destination_dataspace` - (Optional, String, NonUpdatable) Specifies the destination data space ID.

* `destination_dataspace_name` - (Optional, String, NonUpdatable) Specifies the destination data space name.

* `destination_identity_role` - (Optional, String, NonUpdatable) Specifies the destination identity role.

* `destination_pipe` - (Optional, String, NonUpdatable) Specifies the destination pipe ID.

* `destination_pipe_name` - (Optional, String, NonUpdatable) Specifies the destination pipe name.

* `destination_region` - (Optional, String, NonUpdatable) Specifies the destination region.

* `destination_shipper_type` - (Optional, Int, NonUpdatable) Specifies the destination shipper type.

* `destination_workspace` - (Optional, String, NonUpdatable) Specifies the destination workspace ID.

* `destination_workspace_name` - (Optional, String, NonUpdatable) Specifies the destination workspace name.

<a name="shipper_source_block"></a>
The `shipper_source` block supports:

* `region` - (Optional, String, NonUpdatable) Specifies the source region.

* `source_dataspace` - (Optional, String, NonUpdatable) Specifies the source data space ID.

* `source_dataspace_name` - (Optional, String, NonUpdatable) Specifies the source data space name.

* `source_identity_role` - (Optional, String, NonUpdatable) Specifies the source identity role.

* `source_pipe` - (Optional, String, NonUpdatable) Specifies the source pipe ID.

* `source_pipe_name` - (Optional, String, NonUpdatable) Specifies the source pipe name.

* `source_type` - (Optional, Int, NonUpdatable) Specifies the source type.

* `source_workspace` - (Optional, String, NonUpdatable) Specifies the source workspace ID.

* `source_workspace_name` - (Optional, String, NonUpdatable) Specifies the source workspace name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the shipper name.

* `shipper_id` - The shipper ID.

* `status` - The status of the SIEM shipper.

* `create_time` - The creation time, in milliseconds.

* `update_time` - The update time, in milliseconds.

## Import

The SIEM shipper can be imported using the `workspace_id` and the `shipper_name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_siem_shipper.test <workspace_id>/<shipper_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `workspace_name`, `action`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource,
or the resource definition should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_siem_shipper" "test" {
  ...

  lifecycle {
    ignore_changes = [
      workspace_name,
      action,
    ]
  }
}
```
