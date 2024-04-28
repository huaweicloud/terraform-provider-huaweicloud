---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_service_permissions"
description: ""
---

# huaweicloud_vpcep_service_permissions

Use this data source to get VPC endpoint service permissions.

## Example Usage

```hcl
variable service_id {}
variable permission {}

data "huaweicloud_vpcep_service_permissions" "permissions" {
  service_id = var.service_id
  permission = var.permission
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the VPC endpoint service. If omitted, the
  provider-level region will be used.

* `service_id` - (Required, String) Specifies the ID of VPC endpoint service.

* `permission` - (Optional, String) Specifies the account or organization to access the VPC endpoint service.
  The permission format is **iam:domain::domain_id** or **organizations:orgPath::org_path**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permissions` - The list of VPC endpoint service permissions.

The `permissions` block supports:

* `permission_id` - The ID of permission.

* `permission` - The account or organization to access the VPC endpoint service.

* `permission_type` - The permission type of the VPC endpoint service. The value can be **domainId** or **orgPath**.

* `created_at` - The creation time of VPC endpoint permission.
