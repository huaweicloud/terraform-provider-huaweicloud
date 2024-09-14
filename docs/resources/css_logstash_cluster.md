---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_cluster"
description: ""
---

# huaweicloud_css_logstash_cluster

Manages CSS logstash cluster resource within HuaweiCloud

## Example Usage

### create a logstash cluster

```hcl
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}

resource "huaweicloud_css_logstash_cluster" "test" {
  name           = "tf_test_cluster"
  engine_version = "7.10.0"

  node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = var.availability_zone
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the logstash cluster resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new cluster resource.

* `name` - (Required, String) Specifies the cluster name. It contains `4` to `32` characters.
  Only letters, digits, hyphens (-), and underscores (_) are allowed. The value must start with a letter.

* `engine_version` - (Required, String, ForceNew) Specifies the engine version.
  [For details](https://support.huaweicloud.com/intl/en-us/bulletin-css/css_05_0001.html)
  Changing this parameter will create a new resource.

* `node_config` - (Required, List) Specifies the config of data node.
  The [node_config](#Css_node_config) structure is documented below.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone name.
  Separate multiple AZs with commas (,), for example, az1,az2. AZs must be unique. The number of nodes must be greater
  than or equal to the number of AZs. If the number of nodes is a multiple of the number of AZs, the nodes are evenly
  distributed to each AZ. If the number of nodes is not a multiple of the number of AZs, the absolute difference
  between node quantity in any two AZs is **1** at most.
  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet ID.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the security group ID.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the logstash cluster.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the CSS logstash cluster,
  The value `0` indicates the default enterprise project.

* `charging_mode` - (Optional, String) Specifies the charging mode of the CSS logstash cluster.
  The valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

* `period_unit` - (Optional, String) Specifies the charging period unit of the instance.
  The valid values are **month** and **year**.

* `period` - (Optional, Int) Specifies the charging period of the instance.
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.

  -> **NOTE:** `charging_mode`, `period_unit`, `period` can only be updated when changing
  from **postPaid** to **prePaid** billing mode.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.
  The valid values are **true** and **false**, defaults to **false**.

* `routes` - (Optional, List) Specifies the list of route objects.
  The [routes](#Css_route) structure is documented below.

<a name="Css_node_config"></a>
The `node_config` block supports:

* `flavor` - (Required, String, ForceNew) Specifies the flavor name. The value options are as follows:
  + **ess.spec-4u8g**: The value range of the flavor is `40` GB to `1,500` GB.
  + **ess.spec-4u16g**: The value range of the flavor is `40` GB to `1,600` GB.
  + **ess.spec-4u32g**: The value range of the flavor is `40` GB to `2,560` GB.
  + **ess.spec-8u16g**: The value range of the flavor is `80` GB to `1,600` GB.
  + **ess.spec-8u32g**: The value range of the flavor is `80` GB to `3,200` GB.
  + **ess.spec-8u64g**: The value range of the flavor is `80` GB to `5,120` GB.
  + **ess.spec-16u32g**: The value range of the flavor is `100` GB to `3,200` GB.
  + **ess.spec-16u64g**: The value range of the flavor is `100` GB to `6,400` GB.
  + **ess.spec-32u64g**: The value range of the flavor is `100` GB to `10,240` GB.
  + **ess.spec-32u128g**: The value range of the flavor is `100` GB to `10,240` GB.
  Changing this parameter will create a new resource.

* `instance_number` - (Required, Int) Specifies the number of cluster instances. The value range is `1` to `32`.

* `volume` - (Optional, List, ForceNew) Specifies the information about the volume.
  The [volume](#Css_volume) structure is documented below. Changing this parameter will create a new resource.

<a name="Css_volume"></a>
The `volume` block supports:

* `size` - (Required, Int, ForceNew) Specifies the volume size in GB, which must be a multiple of `10`.
  Changing this parameter will create a new resource.

* `volume_type` - (Required, String, ForceNew) Specifies the volume type. The value options are as follows:
  + **HIGH**: High I/O. The SAS disk is used.
  + **ULTRAHIGH**: Ultra-high I/O. The solid-state drive (SSD) is used.
  + **ESSD**: Extreme speed I/O. The SATA disk is used.
  Changing this parameter will create a new resource.

<a name="Css_route"></a>
The `routes` block supports:

* `ip_address` - (Required, String) Specifies the route ip address.

* `ip_net_mask` - (Required, String) Specifies the subnet mask of the route ip address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `nodes` - List of node objects. The [nodes](#Css_nodes_attr) structure is documented below.

* `engine_type` - The engine type.

* `created_at` - The creation time. The format is ISO8601: **CCYY-MM-DDThh:mm:ss**.

* `endpoint` - The IP address and port number.

* `status` - The CSS logstash cluster status. The value are as follows:
  + **100**: The operation, such as instance creation, is in progress.
  + **200**: The CSS logstash cluster is available.
  + **303**: The CSS logstash cluster is unavailable.

* `updated_at` - Time when a cluster is updated. The format is ISO8601: CCYY-MM-DDThh:mm:ss.

* `is_period` - Whether a cluster is billed on the yearly/monthly mode.

<a name="Css_nodes_attr"></a>
The `nodes` block supports:

* `id` - Instance ID.

* `name` - Instance name.

* `type` - Node type.

* `availability_zone` - The availability zone where the instance resides.

* `status` - Instance status.

* `spec_code` - Instance specification code.

* `ip` - Instance IP address.

* `resource_id` - The resource ID of this instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

CSS logstash cluster can be imported by `id`, e.g.

```bash
terraform import huaweicloud_css_logstash_cluster.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `charging_mode`, `period`, `period_unit`, `auto_renew`.
It is generally recommended running `terraform plan` after importing an cluster.
You can then decide if changes should be applied to the cluster, or the resource definition should be updated
to align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_css_logstash_cluster" "test" {
    ...

  lifecycle {
    ignore_changes = [
      charging_mode, period, period_unit, auto_renew,
    ]
  }
}
```
