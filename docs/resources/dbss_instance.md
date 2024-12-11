---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_instance"
description: |-
  Manages a DBSS instance resource within HuaweiCloud.
---

# huaweicloud_dbss_instance

Manages a DBSS instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "product_spec_desc" {}

resource "huaweicloud_dbss_instance" "test" {
  name               = var.name
  flavor             = "c3ne.xlarge.4"
  availability_zone  = var.availability_zone
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  security_group_id  = var.security_group_id
  product_spec_desc  = var.product_spec_desc
  charging_mode      = "prePaid"
  period_unit        = "month"
  period             = 1
  resource_spec_code = "dbss.bypassaudit.low"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the instance name. The name can contain `1` to `64` characters.
  Only letters, digits, underscores (_), and hyphens (-) are allowed.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone to which the instance belongs.
  Primary and secondary AZs are separated by commas. Example: **cn-north-4a,cn-north-4b**.

  Changing this parameter will create a new resource.

* `flavor` - (Required, String, ForceNew) Specifies the flavor. Possible values are:
  + **c3ne.xlarge.4**: Basic version.
  + **c3ne.2xlarge.4**: Professional version.
  + **c6.4xlarge.4**: Advanced version.

  Changing this parameter will create a new resource.

* `resource_spec_code` - (Required, String, ForceNew) Specifies the resource specifications. Possible values are:
  + **dbss.bypassaudit.low**: Basic version.
  + **dbss.bypassaudit.medium**: Professional version.
  + **dbss.bypassaudit.high**: Advanced version.

  Changing this parameter will create a new resource.

* `product_spec_desc` - (Required, String, ForceNew) Specifies the product specification description in
  JSON string format: `{"specDesc":{"zh-cn":{"key1":"value1"},"en-us":{"key1":"value1"}}}`

  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID.

  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet ID of the NIC.

  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the ID of the security group.

* `charging_mode` - (Required, String, ForceNew) Specifies the billing mode.  
  The valid value is **prePaid**, which means the yearly/monthly billing mode.

  Changing this parameter will create a new resource.

* `period_unit` - (Required, String, ForceNew) Specifies the charging period unit.  
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `period` - (Required, Int, ForceNew) Specifies the charging period.  
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto-renew is enabled. Valid values are **true** and
  **false**. Defaults to **false**.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) Specifies the Enterprise project ID.
  For enterprise users, if omitted, default enterprise project will be used.

* `ip_address` - (Optional, String, ForceNew) Specifies the IP address.
  If the value of this parameter is left blank or is set to an empty string, the IP address is automatically assigned.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the instance.

* `action` - (Optional, String) Specifies operation the DBSS instance status.
  The valid values are as follows:
  + **start**: Indicates enable the DBSS instance.
  + **stop**: Indicates disable the DBSS instance.
  + **reboot**: Indicates restart the DBSS instance.

  -> 1.After a DBSS instance created, the default status is **start**.
  <br/>2.The same operation cannot be performed repeatedly.
  <br/>3.The **stop** or **reboot** operation can only be performed when the DBSS instance has no tasks
  and is in **start** status.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `connect_ip` - The connection address.

* `connect_ipv6` - The IPv6 address.

* `created_at` - The creation time

* `expired_at` - The expired time

* `port_id` - The ID of the port that the EIP is bound to.

* `status` - The instance status. The value can be:
  + **SHUTOFF**: disabled;
  + **ACTIVE**: operations allowed;
  + **DELETING**: no operations allowed;
  + **BUILD**: no operations allowed;
  + **DELETED**: not displayed;
  + **ERROR**: only deletion allowed;
  + **HAWAIT**: waiting for the standby to be created, no operations allowed;
  + **FROZEN**: only renewal, binding, and unbinding allowed;
  + **UPGRADING**: no operations allowed;

* `instance_id` - The ID of the audit instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 20 minutes.

## Import

The DBSS instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dbss_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `charging_mode`, `enterprise_project_id`, `flavor`, `period`,
`period_unit`, `product_spec_desc`, `tags` and `action`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to align
with the instance. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_dbss_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      charging_mode, enterprise_project_id, flavor, period, period_unit, product_spec_desc, tags, action,
    ]
  }
}
```
