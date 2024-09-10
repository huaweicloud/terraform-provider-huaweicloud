---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cci_namespace"
description: ""
---

# huaweicloud_cci_namespace

Manages a CCI namespace resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace_name" {}

resource "huaweicloud_cci_namespace" "test" {
  name         = var.namespace_name
  type         = "gpu-accelerated"
  rbac_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCI namespace resource.
  If omitted, the provider-level region will be used. Changing this will create a new CCI namespace resource.

* `type` - (Required, String, ForceNew) Specifies the CCI namespace type.
  The valid values are **general-computing** and **gpu-accelerated**.
  Changing this will create a new CCI namespace resource.

* `name` - (Required, String, ForceNew) Specifies the unique name of the CCI namespace.  
  This parameter can contain a maximum of `63` characters, which may consist of lowercase letters, digits and
  hyphens (-), and must start and end with lowercase letters and digits.  
  Changing this will create a new CCI namespace resource.

* `auto_expend_enabled` - (Optional, Bool, ForceNew) Specifies whether elastic scheduling is enabled.
  Changing this will create a new CCI namespace resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies a unique ID in UUID format of enterprise project.
  Changing this will create a new CCI namespace resource.

  ->**NOTE:** If the enterprise project selected by namespace is different from the enterprise project owned by the VPC,
  the created namespace may not work normally due to permissions.

* `warmup_pool_size` - (Optional, Int, ForceNew) Specifies the size of IP pool to warm-up.  
  The valid value is range from `1` to `500`.
  Changing this will create a new CCI namespace resource.

* `recycling_interval` - (Optional, Int, ForceNew) Specifies the IP address recycling interval, in hour.
  The idle IP resources from the elastic expansion of the IP resource pool can be recycled within this time.
  Changing this will create a new CCI namespace resource.

* `container_network_enabled` - (Optional, Bool, ForceNew) Specifies whether container network is enabled.
  Enable this option if you want CCI to start the container network in advance so that containers can connect to the
  network as soon as they are started. Default to **false**.
  Changing this will create a new CCI namespace resource.

* `rbac_enabled` - (Optional, Bool, ForceNew) Specifies whether Role-based access control is enabled.
  After the RBAC permission is enabled, the user's use of resources under the namespace will be controlled by the RBAC
  permission. Changing this will create a new CCI namespace resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Namespace ID.

* `created_at` - The time when the namespace was created, in UTC format, e.g., **2021-09-27T01:30:39Z**.

* `status` - Namespace status.

## Import

CCI Namespaces can be imported using their `name`, e.g.,

```bash
$ terraform import huaweicloud_cci_namespace.test terraform-test
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 3 minutes.
