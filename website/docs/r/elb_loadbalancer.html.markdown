---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancer"
sidebar_current: "docs-huaweicloud-resource-elb-loadbalancer"
description: |-
  Manages an elastic loadbalancer resource within huawei cloud.
---

# huaweicloud\_elb\_loadbalancer

Manages an elastic loadbalancer resource within huawei cloud.

## Example Usage

```hcl
resource "huaweicloud_elb_loadbalancer" "elb" {
  name           = "elb"
  type           = "External"
  description    = "test elb"
  vpc_id         = "e346dc4a-d9a6-46f4-90df-10153626076e"
  admin_state_up = 1
  bandwidth      = 5
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the load balancer name. The name is a string
    of 1 to 64 characters that consist of letters, digits, underscores (_),
    and hyphens (-).

* `description` - (Optional) Provides supplementary information about the
    listener. The value is a string of 0 to 128 characters and cannot be <>.

* `vpc_id` - (Required) Specifies the VPC ID.

* `bandwidth` - (Optional) Specifies the bandwidth (Mbit/s). This parameter
    is mandatory when type is set to External, and it is invalid when type
    is set to Internal. The value ranges from 1 to 300.

* `type` - (Required) Specifies the load balancer type. The value can be
    Internal or External.

* `admin_state_up` - (Required) Specifies the status of the load balancer.
    Value range: 0 or false: indicates that the load balancer is stopped. Only
    tenants are allowed to enter these two values. 1 or true: indicates that
    the load balancer is running properly. 2 or false: indicates that the load
    balancer is frozen. Only tenants are allowed to enter these two values.

* `vip_subnet_id` - (Optional) Specifies the ID of the private network
    to be added. This parameter is mandatory when type is set to Internal,
    and it is invalid when type is set to External.

* `az` - (Optional) Specifies the ID of the availability zone (AZ). This
    parameter is mandatory when type is set to Internal, and it is invalid
    when type is set to External.

* `charge_mode` - (Optional) This is a reserved field. If the system supports
    charging by traffic and this field is specified, then you are charged by
    traffic for elastic IP addresses. The value is traffic.

* `eip_type` - (Optional) This parameter is reserved.

* `security_group_id` - (Optional) Specifies the security group ID. The
    value is a string of 1 to 200 characters that consists of uppercase and
    lowercase letters, digits, and hyphens (-). This parameter is mandatory
    only when type is set to Internal.

* `vip_address` - (Optional) Specifies the IP address provided by ELB.
    When typeis set to External, the value of this parameter is the elastic
    IP address. When type is set to Internal, the value of this parameter is
    the private network IP address. You can select an existing elastic IP address
    and create a public network load balancer. When this parameter is configured,
    parameters bandwidth, charge_mode, and eip_type are invalid.

* `tenantid` - (Optional) Specifies the tenant ID. This parameter is mandatory
    only when type is set to Internal.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `bandwidth` - See Argument Reference above.
* `type` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `vip_subnet_id` - See Argument Reference above.
* `az` - See Argument Reference above.
* `charge_mode` - See Argument Reference above.
* `eip_type` - See Argument Reference above.
* `security_group_id` - See Argument Reference above.
* `vip_address` - See Argument Reference above.
* `tenantid` - See Argument Reference above.
* `update_time` - Specifies the time when information about the load balancer
    was updated.
* `create_time` - Specifies the time when the load balancer was created.
* `id` - Specifies the load balancer ID.
* `status` - Specifies the status of the load balancer. The value can be
    ACTIVE, PENDING_CREATE, or ERROR.
