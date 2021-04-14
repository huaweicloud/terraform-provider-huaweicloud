---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud\_cce\_addon

To get the specified addon in a cluster.

## Example Usage

```hcl
variable "cluster_id" { }
variable "tempalte_name" { }

data "huaweicloud_cce_addon" "addon" {
  cluster_id = var.cluster_id
  name       = var.tempalte_name
}
```
## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to create the cce addon resource.
 If omitted, the provider-level region will be used.
* `cluster_id` - (Required, String) Specifies the ID of the cluster.
* `addon_id` - (Optional, String) Specifies the ID of the addon.
* `template_name` - (Optional, String) Specifies the name of the addon template.
* `version` - (Optional, String) Specifies the version of the addon.
* `status` - (Optional, String) Specifies the addon status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

 * `id` - ID of the addon instance.
 * `description` - Description of addon instance.
