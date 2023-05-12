---
subcategory: "Data Warehouse Service (DWS)"
---

# huaweicloud_dws_cluster

Manages a cluster in the Data Warehouse Service.

## Example Usage

```hcl
variable "availability_zone" {}
variable "dws_cluster_name" {}
variable "user_name" {}
variable "user_pwd" {}
variable "vpc_id" {}
variable "network_id" {}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "sg_dws"
  description = "terraform security group"
}

resource "huaweicloud_dws_cluster" "cluster" {
  name              = var.dws_cluster_name
  node_type         = "dws.m3.xlarge"
  number_of_node    = 3
  availability_zone = var.availability_zone
  user_name         = var.user_name
  user_pwd          = var.user_pwd
  vpc_id            = var.vpc_id
  network_id        = var.network_id
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the cluster resource. If omitted, the
  provider-level region will be used. Changing this creates a new cluster resource.

* `name` - (Required, String, ForceNew) Cluster name, which must be unique and contains 4 to 64 characters, which
  consist of letters, digits, hyphens(-), or underscores(_) only and must start with a letter.
  Changing this creates a new cluster resource.

* `node_type` - (Required, String, ForceNew) Node type.
  Changing this creates a new cluster resource.

* `number_of_node` - (Required, Int) Number of nodes in a cluster. The value ranges from 3 to 32. When expanding,
  add at least 3 nodes.

* `user_name` - (Required, String, ForceNew) Administrator username for logging in to a data warehouse cluster The
  administrator username must:  Consist of lowercase letters, digits, or underscores. Start with a lowercase letter or
  an underscore. Contain 1 to 63 characters. Cannot be a keyword of the DWS database.
  Changing this creates a new cluster resource.

* `user_pwd` - (Required, String) Administrator password for logging in to a data warehouse cluster A password
  must conform to the following rules:  Contains 8 to 32 characters. Cannot be the same as the username or the username
  written in reverse order. Contains three types of the following:
  Lowercase letters Uppercase letters Digits Special characters
  ~!@#%^&*()-_=+|[{}];:,<.>/?

* `vpc_id` - (Required, String, ForceNew) VPC ID, which is used for configuring cluster network.
  Changing this creates a new cluster resource.

* `network_id` - (Required, String, ForceNew) Network ID, which is used for configuring cluster network.
  Changing this creates a new cluster resource.

* `security_group_id` - (Required, String, ForceNew) ID of a security group. The ID is used for configuring cluster
  network.
  Changing this creates a new cluster resource.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the dws cluster,
  Value 0 indicates the default enterprise project.
  Changing this creates a new cluster resource.

* `availability_zone` - (Optional, String, ForceNew) AZ in a cluster.
  Changing this creates a new cluster resource.

* `port` - (Optional, Int, ForceNew) Service port of a cluster (8000 to 10000). The default value is 8000.
  Changing this creates a new cluster resource.

* `number_of_cn` - (Optional, Int, ForceNew) The number of CN. If you use a large-scale cluster, deploy multiple CNs.
  Changing this creates a new cluster resource.

* `tags` - (Optional, Map) The key/value pairs to associate with the cluster.

* `public_ip` - (Optional, List, ForceNew) A nested object resource Structure is documented below.
  Changing this creates a new cluster resource.

The `public_ip` block supports:

* `eip_id` - (Optional, String, ForceNew) EIP ID.
  Changing this creates a new cluster resource.

* `public_bind_type` - (Optional, String, ForceNew) Binding type of an EIP. The value can be either of the following:
  **auto_assign**, **not_use**, **bind_existing**. The default value is **not_use**.
  Changing this creates a new cluster resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The cluster ID

* `endpoints` - View the private network connection information about the cluster. Structure is documented below.

* `public_endpoints` - Public network connection information about the cluster. If the value is not specified, the
  public network connection information is not used by default Structure is documented below.

* `recent_event` - The recent event number.

* `status` - Cluster status, which can be one of the following:  CREATING AVAILABLE UNAVAILABLE CREATION FAILED.

* `sub_status` - Sub-status of clusters in the AVAILABLE state. The value can be one of the following:  NORMAL READONLY
  REDISTRIBUTING REDISTRIBUTION-FAILURE UNBALANCED UNBALANCED | READONLY DEGRADED DEGRADED | READONLY DEGRADED |
  UNBALANCED UNBALANCED | REDISTRIBUTING UNBALANCED | REDISTRIBUTION-FAILURE READONLY | REDISTRIBUTION-FAILURE
  UNBALANCED | READONLY | REDISTRIBUTION-FAILURE DEGRADED | REDISTRIBUTION-FAILURE DEGRADED | UNBALANCED |
  REDISTRIBUTION-FAILURE DEGRADED | UNBALANCED | READONLY | REDISTRIBUTION-FAILURE DEGRADED | UNBALANCED | READONLY

* `task_status` - Cluster management task. The value can be one of the following:
  RESTORING SNAPSHOTTING GROWING REBOOTING SETTING_CONFIGURATION CONFIGURING_EXT_DATASOURCE DELETING_EXT_DATASOURCE
  REBOOT_FAILURE RESIZE_FAILURE

* `created` - Cluster creation time. The format is ISO8601:YYYY-MM-DDThh:mm:ssZ

* `updated` - Last modification time of a cluster. The format is ISO8601:YYYY-MM-DDThh:mm:ssZ

* `version` - Data warehouse version.

* `private_ip` - List of private network IP address.

The `endpoints` block contains:

* `connect_info` - (Optional, String) Private network connection information.

* `jdbc_url` - (Optional, String)
  JDBC URL. The following is the default format:
  jdbc:postgresql://< connect_info>/<YOUR_DATABASE_NAME>

The `public_endpoints` block contains:

* `jdbc_url` - (Optional, String)
  JDBC URL. The following is the default format:
  jdbc:postgresql://< public_connect_info>/<YOUR_DATABASE_NAME>

* `public_connect_info` - (Optional, String)
  Public network connection information.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

Cluster can be imported using the following format:

```
$ terraform import huaweicloud_dws_cluster.test 47ad727e-9dcc-4833-bde0-bb298607c719
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `user_pwd`, `number_of_cn`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```
resource "huaweicloud_dws_cluster" "test" {
    ...

  lifecycle {
    ignore_changes = [
      user_pwd, number_of_cn,
    ]
  }
}
```
