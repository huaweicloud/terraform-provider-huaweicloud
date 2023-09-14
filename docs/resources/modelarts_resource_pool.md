---
subcategory: "AI Development Platform (ModelArts)"
---

# huaweicloud_modelarts_resource_pool

Manages a ModelArts dedicated resource pool resource within HuaweiCloud.  

## Example Usage

```hcl
variable "modelarts_network_id" {}

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "demo"
  description = "This is a demo"
  scope       = ["Train", "Infer", "Notebook"]
  network_id  = var.modelarts_network_id

  resources {
    flavor_id = "modelarts.vm.cpu.16u64g.d"
    count     = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The name of the dedicated resource pool.  
    The name can contain `4` to `32` characters, only lowercase letters, digits and hyphens (-) are allowed.
    The name must start with a lowercase letter and end with a lowercase letter or digit.

  Changing this parameter will create a new resource.

* `scope` - (Required, List) List of job types supported by the resource pool.  
  The options are as follows:
    + **Train**: training job.
    + **Infer**: inference job.
    + **Notebook**: Notebook job.

* `resources` - (Required, List) List of resource specifications in the resource pool.  
  Including resource flavors and the number of resources of the corresponding flavors.
  The [resources](#ModelartsResourcePool_ResourceFlavor) structure is documented below.

* `network_id` - (Required, String, ForceNew) The ModelArt network ID of the resource pool.

  Changing this parameter will create a new resource.

* `workspace_id` - (Optional, String, ForceNew) Workspace ID, which defaults to 0.  

  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description of the dedicated resource pool.  

<a name="ModelartsResourcePool_ResourceFlavor"></a>
The `resources` block supports:

* `flavor_id` - (Required, String) The resource flavor ID.  

* `count` - (Required, Int) Number of resources of the corresponding flavors.

* `azs` - (Optional, List) AZs for resource pool nodes.
  The [azs](#ModelartsResourcePool_ResourceFlavor_azs) structure is documented below.

<a name="ModelartsResourcePool_ResourceFlavor_azs"></a>
The `azs` block supports:

* `az` - (Optional, String) The AZ name.

* `count` - (Optional, Int) Number of nodes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The ModelArts resource pool can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_resource_pool.test 0ce123456a00f2591fabc00385ff1234
```
