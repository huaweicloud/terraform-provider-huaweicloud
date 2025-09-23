---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_instant_tasks"
description: |-
  Use this data source to get the list of GaussDB MySQL instant tasks.
---

# huaweicloud_gaussdb_mysql_instant_tasks

Use this data source to get the list of GaussDB MySQL instant tasks.

## Example Usage

```hcl
data "huaweicloud_gaussdb_mysql_instant_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) Specifies the task execution status. Value options:
  + **Running**: The task is being executed.
  + **Completed**: The task is successfully executed.
  + **Failed**: The task failed to be executed.
  + **Pending**: The task is not executed.

* `job_id` - (Optional, String) Specifies the task ID.

* `job_name` - (Optional, String) Specifies the task name. Value options:
  + **CreateGaussDBforMySQLInstance**: Creating a DB instance.
  + **RestoreGaussDBforMySQLNewInstance**: Restoring data to a new DB instance.
  + **AddGaussDBforMySQLNodes**: Adding nodes.
  + **DeleteGaussDBforMySQLNode**: Deleting nodes.
  + **RebootGaussDBforMySQLInstance**: Rebooting a DB instance.
  + **ModifyGaussDBforMySQLPort**: Changing a database port.
  + **ModifyGaussDBforMySQLSecurityGroup**: Changing a security group.
  + **ResizeGaussDBforMySQLFlavor**: Changing instance specifications.
  + **SwitchoverGaussDBforMySQLMasterNode**: Promoting a read replica to primary.
  + **GaussDBforMySQLBindEIP**: Binding an EIP.
  + **GaussDBforMySQLUnbindEIP**: Unbinding an EIP.
  + **RenameGaussDBforMySQLInstance**: Changing a DB instance name.
  + **DeleteGaussDBforMySQLInstance**: Deleting a DB instance.
  + **UpgradeGaussDBforMySQLDatabaseVersion**: Upgrading an instance version.
  + **EnlargeGaussDBforMySQLProxy**: Adding nodes for a database proxy.
  + **OpenGaussDBforMySQLProxy**: Enabling database proxy.
  + **CloseGaussDBforMySQLProxy**: Disabling database proxy.
  + **GaussdbforMySQLModifyProxyIp**: Changing the IP address of a database proxy.
  + **ScaleGaussDBforMySQLProxy**: Changing the node specifications of a database proxy.
  + **GaussDBforMySQLModifyInstanceMetricExtend**: Enabling or disabling Monitoring by Seconds.
  + **GaussDBforMySQLModifyInstanceDataVip**: Changing the private IP address.
  + **GaussDBforMySQLSwitchSSL**: Enabling or disabling SSL.
  + **GaussDBforMySQLModifyProxyConsist**: Changing the proxy consistency.
  + **GaussDBforMySQLModifyProxyWeight**: Changing the read weights of nodes.

* `start_time` - (Optional, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Optional, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - Indicates the task details.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `instance_status` - Indicates the iInstance status. The value can be:
  + **createfail**: The instance failed to be created.
  + **creating**: The instance is being created.
  + **normal**: The instance is running properly.
  + **abnormal**: The instance is abnormal.
  + **deleted**: The instance has been deleted.

* `job_id` - Indicates the task ID.

* `job_name` - Indicates the task name. The value can be:
  + **CreateGaussDBforMySQLInstance**: Creating a DB instance.
  + **RestoreGaussDBforMySQLNewInstance**: Restoring data to a new DB instance.
  + **AddGaussDBforMySQLNodes**: Adding nodes.
  + **DeleteGaussDBforMySQLNode**: Deleting nodes.
  + **RebootGaussDBforMySQLInstance**: Rebooting a DB instance.
  + **ModifyGaussDBforMySQLPort**: Changing a database port.
  + **ModifyGaussDBforMySQLSecurityGroup**: Changing a security group.
  + **ResizeGaussDBforMySQLFlavor**: Changing instance specifications.
  + **SwitchoverGaussDBforMySQLMasterNode**: Promoting a read replica to primary.
  + **GaussDBforMySQLBindEIP**: Binding an EIP.
  + **GaussDBforMySQLUnbindEIP**: Unbinding an EIP.
  + **RenameGaussDBforMySQLInstance**: Changing a DB instance name.
  + **DeleteGaussDBforMySQLInstance**: Deleting a DB instance.
  + **UpgradeGaussDBforMySQLDatabaseVersion**: Upgrading an instance version.
  + **EnlargeGaussDBforMySQLProxy**: Adding nodes for a database proxy.
  + **OpenGaussDBforMySQLProxy**: Enabling database proxy.
  + **CloseGaussDBforMySQLProxy**: Disabling database proxy.
  + **GaussdbforMySQLModifyProxyIp**: Changing the IP address of a database proxy.
  + **ScaleGaussDBforMySQLProxy**: Changing the node specifications of a database proxy.
  + **GaussDBforMySQLModifyInstanceMetricExtend**: Enabling or disabling Monitoring by Seconds.
  + **GaussDBforMySQLModifyInstanceDataVip**: Changing the private IP address.
  + **GaussDBforMySQLSwitchSSL**: Enabling or disabling SSL.
  + **GaussDBforMySQLModifyProxyConsist**: Changing the proxy consistency.
  + **GaussDBforMySQLModifyProxyWeight**: Changing the read weights of nodes.

* `status` - Indicates the task execution status.
  The value can be:
  + **Pending**: The task is delayed and not executed.
  + **Running**: The task is being executed.
  + **Completed**: The task is successfully executed.
  + **Failed**: The task failed to be executed.

* `process` - Indicates the task progress.

* `fail_reason` - Indicates the task failure cause.

* `order_id` - Indicates the order ID.

* `created_time` - Indicates the task creation time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `ended_time` - Indicates the task end time in the **yyyy-mm-ddThh:mm:ssZ** format.
