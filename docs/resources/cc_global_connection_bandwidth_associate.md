---
subcategory: Cloud Connect (CC)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_global_connection_bandwidth_associate"
description: ""
---

# huaweicloud_cc_global_connection_bandwidth_associate

Associate a global connection bandwidth to GEIP or some other resource instances within HuaweiCloud.

## Example Usage

```hcl
variable "project_id" {}
variable "gcb_id" {}
variable "resource_id" {}

resource "huaweicloud_cc_global_connection_bandwidth_associate" "test" {
  gcb_id = var.gcb_id

  gcb_binding_resources {
    resource_id   = var.resource_id
    resource_type = "GEIP"
    region_id     = "global"
    project_id    = var.project_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `gcb_binding_resources` - (Required, List) The resources to associate with the global connection bandwidth.
  The [gcb_binding_resources](#GCB_Binding_Resources) structure is documented below.

* `gcb_id` - (Required, String, ForceNew) The global connection bandwidth ID.
  Changing this creates a new resource.

<a name="GCB_Binding_Resources"></a>
The `gcb_binding_resources` block supports:

* `resource_id` - (Required, String) The ID of the resource to associate with the global connection bandwidth.

* `resource_type` - (Required, String) The type of the resource to associate with the global connection bandwidth.
  Currently, only **GEIP** is supported.

* `project_id` - (Optional, String) The project ID of the resource to associate with the global connection bandwidth.

* `region_id` - (Optional, String) The region ID of the resource to associate with the global connection bandwidth.
  If the value of `resource_type` is **GEIP**, the valid value is **global**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the global connection bandwidth ID.

## Import

The global connection bandwidth associate resource can be imported using the global connection bandwidth ID, e.g.

```bash
$ terraform import huaweicloud_cc_global_connection_bandwidth_associate.test <gcb_id>
```
