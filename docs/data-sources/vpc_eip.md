---
subcategory: "Elastic IP (EIP)"
---

# huaweicloud\_vpc\_eip

Use this data source to get the details of an available Eip.

## Example Usage

```hcl

variable "port_id" {}

data "huaweicloud_vpc_eip" "by_port_id" {
  port_id = var.port_id
}

```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the eip. If omitted, the provider-level region will be used.

* `public_ip` - (Optional, String) The public ip address of the eip.

* `port_id` - (Optional, String) The port id of the eip.

* `enterprise_project_id` - (Optional, String) The enterprise project id of the eip.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - The status of the eip.

* `type` - The type of the eip.

* `private_ip` - The private ip of the eip.

* `bandwidth_id` - The bandwidth id of the eip.

* `bandwidth_size` - The bandwidth size of the eip.

* `bandwidth_share_type` - The bandwidth share type of the eip.
