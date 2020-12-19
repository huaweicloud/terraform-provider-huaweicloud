---
subcategory: "Cloud Container Engine (CCE)"
---

## Example Usage
```hcl
variable "cluster_id" { }

resource "huaweicloud_cce_addon" "addon_test" {
    cluster_id    = var.cluster_id
    template_name = "metrics-server"
    version       = "1.0.0"
}
``` 

## Argument Reference
The following arguments are supported:
* `region` - (Optional, String, ForceNew) The region in which to create the cce addon resource. If omitted, the provider-level region will be used. Changing this creates a new cce addon resource.
* `cluster_id` - (Required, String, ForceNew) ID of the cluster. Changing this parameter will create a new resource.
* `template_name` - (Required, String, ForceNew) Name of the addon template. Changing this parameter will create a new resource.
* `version` - (Required, String, ForceNew) Version of the addon. Changing this parameter will create a new resource.
* `values` - (Optional, List, ForceNew) Add-on template installation parameters. These parameters vary depending on the add-on.

The `values` block supports:
* `basic` - (Required, Map) Key/Value pairs vary depending on the add-on.
* `custom` - (Optional, Map) Key/Value pairs vary depending on the add-on.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

 * `id` -  ID of the addon instance.
 * `status` - Addon status information.
 * `description` - Description of addon instance.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 3 minute.

