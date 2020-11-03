---
subcategory: "Elastic Cloud Server (ECS)"
---

# huaweicloud\_compute\_flavors

Use this data source to get the ID of the available Compute Flavors.

## Example Usage

```hcl
data "huaweicloud_compute_flavors" "flavors" {
  availability_zone = "cn-north-4a"
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

# Create ECS instance with the first matched flavor

resource "huaweicloud_compute_instance" "instance" {
  flavor_id = data.huaweicloud_compute_flavors.flavors.ids[0]

  # Other properties...
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the flavors. If omitted, the provider-level region will be used.

* `availability_zone` - (Optional) Specifies the AZ name.

* `performance_type` - (Optional) Specifies the ECS flavor type.

* `generation` - (Optional) Specifies the generation of an ECS type.

* `cpu_core_count` - (Optional) Specifies the number of vCPUs in the ECS flavor.

* `memory_size` - (Optional) Specifies the memory size(GB) in the ECS flavor.


## Attributes Reference

`id` is set to the ID of the found flavors. In addition, the following attributes are exported:

* `ids` - A list of flavor IDs.
