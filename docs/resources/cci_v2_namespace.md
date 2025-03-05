---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cci_v2_namespace"
description: ""
---

# huaweicloud_cci_v2_namespace

Manages a CCI namespace resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace_name" {}

resource "huaweicloud_cci_namespace" "test" {
  name         = var.namespace_name
  flavor       = "gpu-accelerated"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the unique name of the CCI namespace.

* `flavor` - (Required, String, NonUpdatable) Specifies the CCI namespace flavor.
  The valid values are **general-computing** and **gpu-accelerated**.

* `warm_pool_size` - (Optional, String, NonUpdatable) Specifies the size of IP pool to warm-up.

* `container_network_enabled` - (Optional, Bool, NonUpdatable) Specifies whether container network is enabled.
  Enable this option if you want CCI to start the container network in advance so that containers can connect to the
  network as soon as they are started. Default to **false**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies a unique ID in UUID format of enterprise project.

* `rbac_enabled` - (Optional, Bool, NonUpdatable) Specifies whether Role-based access control is enabled. After the RBAC
  permission is enabled, the user's use of resources under the namespace will be controlled by the RBAC permission.

* `warm_pool_recycle_interval` - (Optional, String, NonUpdatable) Specifies the IP address recycling interval, in hour.
  The idle IP resources from the elastic expansion of the IP resource pool can be recycled within this time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is also the name of the namespace.

* `created_at` - The time when the namespace was created, in UTC format, e.g., **2021-09-27T01:30:39Z**.

* `status` - The namespace status.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 3 minutes.

## Import

The CCI Namespaces can be imported using `name`, e.g.

```bash
$ terraform import huaweicloud_cci_v2_namespace.test <name>
```
