---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cci_namespaces"
description: ""
---

# huaweicloud_cci_namespaces

Use this data source to obtain CCI namespaces within HuaweiCloud.

## Example Usage

### Get the specified namespace details

```hcl
variable "namespace_name" {}

data "huaweicloud_cci_namespaces" "test" {
  name = var.namespace_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCI namespace list.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the CCI namespace type.
  The valid values are **general-computing** and **gpu-accelerated**.

* `name` - (Optional, String) Specifies th name of the specified CCI namespace.
  This parameter can contain a maximum of 63 characters, which may consist of lowercase letters, digits and hyphens,
  and must start and end with lowercase letters and digits.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID in UUID format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Data source ID.

* `namespaces` - All CCI namespaces that meet the query parameters.

The `namespaces` block supports:

* `id` - The CCI namespace ID in UUID format.

* `type` - The CCI namespace type.

* `name` - The CCI namespace name.

* `auto_expend_enabled` - Whether elastic scheduling is enabled.

* `enterprise_project_id` - The enterprise project ID in UUID format.

* `warmup_pool_size` - The size of IP pool to warm-up.

* `recycling_interval` - The IP address recycling interval in hour.
  The idle IP resources from the elastic expansion of the IP resource pool can be recycled within this time.

* `container_network_enabled` - Whether container network is enabled.

* `rbac_enabled` - Whether Role-based access control is enabled.
  After the RBAC permission is enabled, the user's use of resources under the namespace will be controlled by the RBAC
  permission.

* `created_at` - The time when the namespace was created in UTC format, such as **2021-09-27T01:30:39Z**.

* `status` - The CCI namespace status.

* `network` - The network information of the CCI namespace. The structure is documented below.

The `network` block supports:

* `name` - The CCI network name.

* `security_group_id` - The default security group ID in UUID format.

* `vpc` - The network information of the VPC under the CCI network. The structure is documented below.

The `vpc` block supports:

* `id` - The VPC ID in UUID format.

* `subnet_id` - The VPC subnet ID in UUID format.

* `subnet_cidr` - The subnet CIDR block.

* `network_id` - The network ID of the VPC subnet in UUID format.
