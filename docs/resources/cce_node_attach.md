---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud_cce_node_attach

Add a node from an existing ecs server to a CCE cluster.

## Basic Usage

```hcl
resource "huaweicloud_cce_node_attach" "test" {
  cluster_id = var.cluster_id
  server_id  = var.server_id
  key_pair   = var.keypair_name
  os         = "EulerOS 2.5"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the CCE node attach resource. If omitted, the
  provider-level region will be used. Changing this creates a new CCE node attach resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the ID of the cluster. Changing this parameter will create a new
  resource.

* `name` - (Optional, String) Specifies the Node Name.

* `server_id` - (Required, String, ForceNew) Specifies the ecs server ID. Changing this parameter will create a new
  resource.

* `os` - (Required, String) Specifies the operating System of the node. Changing this parameter will reset the node.
  + For VM nodes, clusters of v1.13 and later support *EulerOS 2.5* and *CentOS 7.6*.

* `key_pair` - (Optional, String) Specifies the key pair name when logging in to select the key pair mode.
  This parameter and `password` are alternative. Changing this parameter will reset the node.

* `password` - (Optional, String) Specifies the root password when logging in to select the password mode.
  This parameter can be plain or salted and is alternative to `key_pair`.
  Changing this parameter will reset the node.

* `max_pods` - (Optional, Int, ForceNew) Specifies the the maximum number of instances a node is allowed to create.
  Changing this parameter will create a new resource.

* `docker_base_size` - (Optional, Int, ForceNew) Specifies the available disk space of a single docker container on the
  node in device mapper mode. Changing this parameter will create a new resource.

* `lvm_config` - (Optional, String, ForceNew) Specifies the docker data disk configurations. The following is an
  example:

```hcl
  lvm_config = "dockerThinpool=vgpaas/90%VG;kubernetesLV=vgpaas/10%VG"
```

Changing this parameter will create a new resource.

* `preinstall` - (Optional, String, ForceNew) Specifies the script required before installation. The input value can be
  a Base64 encoded string or not. Changing this parameter will create a new resource.

* `postinstall` - (Optional, String, ForceNew) Specifies the script required after installation. The input value can be
  a Base64 encoded string or not. Changing this parameter will create a new resource.

* `labels` - (Optional, Map, ForceNew) Specifies the tags of a Kubernetes node, key/value pair format.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the tags of a VM node, key/value pair format.

* `taints` - (Optional, List, ForceNew) Specifies the taints configuration of the nodes to set anti-affinity.
  Changing this parameter will create a new resource. Each taint contains the following parameters:

  + `key` - (Required, String, ForceNew) A key must contain 1 to 63 characters starting with a letter or digit.
    Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed. A DNS subdomain name can be used
    as the prefix of a key. Changing this parameter will create a new resource.
  + `value` - (Required, String, ForceNew) A value must start with a letter or digit and can contain a maximum of 63
    characters, including letters, digits, hyphens (-), underscores (_), and periods (.). Changing this parameter will
    create a new resource.
  + `effect` - (Required, String, ForceNew) Available options are NoSchedule, PreferNoSchedule, and NoExecute.
    Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `status` - Node status information.
* `private_ip` - Private IP of the CCE node.
* `public_ip` - Public IP of the CCE node.
* `flavor_id` - The flavor ID of the CCE node.
* `availability_zone` - The name of the available partition (AZ).
* `root_volume` - The system disk related configuration.
* `data_volumes` - The data disks related configuration.
* `runtime` - The runtime of the CCE node.
* `ecs_group_id` - The Ecs group ID.
* `subnet_id` - The ID of the subnet to which the NIC belongs.
* `charging_mode` - The charging mode of the CCE node. Valid values are *prePaid* and *postPaid*.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minute.
* `update` - Default is 20 minute.
* `delete` - Default is 20 minute.
