---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_cluster_ports"
description: |-
  Use this data source to get the list of CPCS cluster ports.
---

# huaweicloud_cpcs_cluster_ports

Use this data source to get the list of CPCS cluster ports.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cpcs_cluster_ports" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `result` - The list of the ports.
  The [result](#ports_result_struct) structure is documented below.

<a name="ports_result_struct"></a>
The `result` block supports:

* `id` - The UUID.

* `cluster_id` - The cluster ID.

* `elb_id` - The ELB ID.

* `elb_ip` - The ELB IP address.

* `mode` - The port mode.
  The valid values are as follows:
  + **PROXY**: Indicates proxy mode port.
  + **TUNNEL**: Indicates tunnel mode custom port.
  + **TUNNEL_FIXED**: Indicates fixed tunnel port in tunnel mode.

* `listener_port` - The ELB listerner port.

* `listener_id` - The ELB listerner ID.

* `server_group_port` - The business port of the backend server bound to the backend service group.

* `server_group_id` - The backend service group ID.

* `project_id` - The project ID.

* `validate_time` - The final verification time.

* `wrong` - Whether the resources are abnormal.

* `wrong_msg` - The resource anomaly information.
