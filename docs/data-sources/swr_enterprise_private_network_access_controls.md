---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_private_network_access_controls"
description: |-
  Use this data source to get the list of SWR enterprise private network access control rule list.
---

# huaweicloud_swr_enterprise_private_network_access_controls

Use this data source to get the list of SWR enterprise private network access control rule list.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_private_network_access_controls" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `internal_endpoints` - Indicates the private network access control list.

  The [internal_endpoints](#internal_endpoints_struct) structure is documented below.

<a name="internal_endpoints_struct"></a>
The `internal_endpoints` block supports:

* `id` - Indicates the private network access rule ID.

* `vpcep_endpoint_id` - Indicates the VPCEP endpoint ID.

* `endpoint_ip` - Indicates the endpoint IP.

* `project_id` - Indicates the project ID to which the VPC belongs.

* `project_name` - Indicates the project name to which the VPC belongs.

* `vpc_id` - Indicates the VPC ID.

* `vpc_name` - Indicates the VPC name.

* `vpc_cidr` - Indicates the VPC CIDR block.

* `subnet_id` - Indicates the subnet ID.

* `subnet_name` - Indicates the subnet name.

* `subnet_cidr` - Indicates the subnet CIDR block.

* `status` - Indicates the access control rule status.

* `status_text` - Indicates the access control rule status text.

* `description` - Indicates description of the access control rule.

* `created_at` - Indicates creation time of the access control rule.
