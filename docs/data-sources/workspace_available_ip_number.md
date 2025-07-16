---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_available_ip_number"
description: |-
  Use this data source to get the number of available IPs in a subnet within HuaweiCloud Workspace.
---

# huaweicloud_workspace_available_ip_number

Use this data source to get the number of available IPs in a subnet within HuaweiCloud Workspace.

## Example Usage

```hcl
variable "subnet_id" {}

data "huaweicloud_workspace_available_ip_number" "test" {
  subnet_id = var.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the subnet is located.  
  If omitted, the provider-level region will be used.

* `subnet_id` - (Required, String) Specifies the ID of the subnet to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `available_ip` - The number of available IPs in the subnet.
