---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_resource_group"
description: ""
---

# huaweicloud_ces_resource_group

Manages a CES resource group resource within HuaweiCloud.

## Example Usage

### Add resources manually

```hcl
variable "image_id" {}
variable "flavor_id" {}
variable "security_group_id" {}
variable "availability_zone" {}
variable "subnet_id" {}

resource "huaweicloud_compute_instance" "vm_1" {
  name               = "ecs-test"
  image_id           = var.image_id
  flavor_id          = var.flavor_id
  security_group_ids = [ var.security_group_id ]
  availability_zone  = var.availability_zone

  network {
    uuid = var.subnet_id
  }
}

resource "huaweicloud_ces_resource_group" "test" {
  name = "test"

  resources {
    namespace = "SYS.ECS"
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.vm_1.id
    }
  }

  resources {
    namespace = "SYS.EVS"
    dimensions {
      name  = "disk_name"
      value = "${huaweicloud_compute_instance.vm_1.id}-sda"
    }
  }
}
```

### Add resources from enterprise projects

```hcl
variable "eps_id" {}

resource "huaweicloud_ces_resource_group" "test" {
  name               = "test"
  type               = "EPS"
  associated_eps_ids = [ var.eps_id ]
}
```

### Add resources by tags

```hcl
resource "huaweicloud_ces_resource_group" "test" {
  name = "test"
  type = "TAG"
  tags = {
    key = "value"
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the resource group name.
  This parameter can contain a maximum of 128 characters, which may consist of letters,
  digits, hyphens (-), underscore (_) and Chinese characters.

* `type` - (Optional, String, ForceNew) Specifies the resource group type.
  The value can be **EPS** and **TAG**. If not specified, that means add resources manually.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the resource group.

* `tags` - (Optional, Map) Specifies the key/value to match resources.
  It's required if the value of type is **TAG**.

* `associated_eps_ids` - (Optional, List, ForceNew) Specifies the enterprise project IDs where the resources from.
  It's required if the value of type is **EPS**.

  Changing this parameter will create a new resource.

* `resources` - (Optional, List) Specifies the list of resources to add into the group.
  The [ResourcesOpts](#ResourceGroup_ResourcesOpts) structure is documented below.

<a name="ResourceGroup_ResourcesOpts"></a>
The `resources` block supports:

* `namespace` - (Required, String) Specifies the namespace in **service.item** format.
  **service** and **item** each must be a string that starts with a letter and contains only letters, digits, and
  underscores (_). For details,
  see [Services Interconnected with Cloud Eye](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html).

* `dimensions` - (Required, List) Specifies the list of dimensions.
  The [DimensionOpts](#ResourceGroup_DimensionOpts) structure is documented below.

<a name="ResourceGroup_DimensionOpts"></a>
The `dimensions` block supports:

* `name` - (Required, String) Specifies the dimension name.
  The value can be a string of 1 to 32 characters that must start with a letter
  and contain only letters, digits, underscores (_), and hyphens (-).

* `value` - (Required, String) Specifies the dimension value.
  The value can be a string of 1 to 64 characters that must start with a letter or a number
  and contain only letters, digits, underscores (_), and hyphens (-).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time.

## Import

The resource group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ces_resource_group.test 0ce123456a00f2591fabc00385ff1234
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `resources`.
It is generally recommended running `terraform plan` after importing a resource group.
You can then decide if changes should be applied to the resource group, or the resource definition should be updated to
align with the resource group. Also you can ignore changes as below.

```hcl
resource "huaweicloud_ces_resource_group" "test" {
    ...

  lifecycle {
    ignore_changes = [
      resources,
    ]
  }
}
```
