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
* `region` - (Optional) The region in which to obtain the cce addon resource. If omitted, the provider-level region will work as default. Changing this creates a new cce addon resource.
* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.
* `template_name` - (Required) Name of the addon template. Changing this parameter will create a new resource.
* `version` - (Required) Version of the addon. Changing this parameter will create a new resource.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

 * `id` -  ID of the addon instance.
 * `status` - Addon status information.
 * `description` - Description of addon instance.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 3 minute.

