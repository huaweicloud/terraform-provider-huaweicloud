---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alert_convert_incident"
description: |-
  Manages a resource to convert alert into incident within HuaweiCloud.
---

# huaweicloud_secmaster_alert_convert_incident

Manages a resource to convert alert into incident within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

```hcl
variable "workspace_id" {}
variable "alert_ids" {
  type = list(string)
}
variable "incident_title" {}

resource "huaweicloud_secmaster_alert_convert_incident" "test" {
  workspace_id = var.workspace_id
  ids          = var.alert_ids
  title        = var.incident_title

  incident_type {
    category      = "DDoS"
    incident_type = "ACK Flood"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace.

* `ids` - (Required, List, NonUpdatable) Specifies the IDs of the alerts to be converted into incidents.

* `incident_type` - (Required, List, NonUpdatable) Specifies the incident type details.
  The [incident_type](#convert_incident_type) structure is documented below.

* `title` - (Optional, String, NonUpdatable) Specifies the converted incident name.

<a name="convert_incident_type"></a>
The `incident_type` block supports:

* `id` - (Optional, String, NonUpdatable) Specifies the incident type ID.

* `category` - (Optional, String, NonUpdatable) Specifies the parent incident type.

* `incident_type` - (Optional, String, NonUpdatable) Specifies the child incident type.

-> Exactly one of `id`, `category` or  `incident_type` must be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
