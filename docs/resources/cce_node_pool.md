---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud_cce_node_pool

Add a node pool to a container cluster.

## Example Usage

```hcl
variable "cluster_id" {}
variable "key_pair" {}
variable "availability_zone" {}

resource "huaweicloud_cce_node_pool" "node_pool" {
  cluster_id               = var.cluster_id
  name                     = "testpool"
  os                       = "EulerOS 2.5"
  initial_node_count       = 2
  flavor_id                = "s3.large.4"
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair
  scall_enable             = true
  min_node_count           = 1
  max_node_count           = 10
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the CCE pool resource. If omitted, the
  provider-level region will be used. Changing this creates a new CCE node pool resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the cluster ID.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the node pool name.

* `initial_node_count` - (Required, Int) Specifies the initial number of expected nodes in the node pool.
  This parameter can be also used to manually scale the node count afterwards.

* `flavor_id` - (Required, String, ForceNew) Specifies the flavor ID. Changing this parameter will create a new
  resource.

* `type` - (Optional, String, ForceNew) Specifies the node pool type. Possible values are: **vm** and **ElasticBMS**.

* `availability_zone` - (Optional, String, ForceNew) Specifies the name of the available partition (AZ). Default value
  is random to create nodes in a random AZ in the node pool. Changing this parameter will create a new resource.

* `os` - (Optional, String, ForceNew) Specifies the operating system of the node.
  Changing this parameter will create a new resource.

* `key_pair` - (Optional, String, ForceNew) Specifies the key pair name when logging in to select the key pair mode.
  This parameter and `password` are alternative. Changing this parameter will create a new resource.

* `password` - (Optional, String, ForceNew) Specifies the root password when logging in to select the password mode.
  This parameter can be plain or salted and is alternative to `key_pair`.
  Changing this parameter will create a new resource.

* `subnet_id` - (Optional, String, ForceNew) Specifies the ID of the subnet to which the NIC belongs.
  Changing this parameter will create a new resource.

* `max_pods` - (Optional, Int, ForceNew) Specifies the maximum number of instances a node is allowed to create.
  Changing this parameter will create a new resource.

* `preinstall` - (Optional, String, ForceNew) Specifies the script to be executed before installation.
  The input value can be a Base64 encoded string or not. Changing this parameter will create a new resource.

* `postinstall` - (Optional, String, ForceNew) Specifies the script to be executed after installation.
  The input value can be a Base64 encoded string or not. Changing this parameter will create a new resource.

* `extend_param` - (Optional, Map, ForceNew) Specifies the extended parameter.
  Changing this parameter will create a new resource.
  The available keys are as follows:
  + **agency_name**: The agency name to provide temporary credentials for CCE node to access other cloud services.
  + **alpha.cce/NodeImageID**: The custom image ID used to create the BMS nodes.
  + **dockerBaseSize**: The available disk space of a single docker container on the node in device mapper mode.
  + **DockerLVMConfigOverride**: Specifies the data disk configurations of Docker.
  
  The following is an example default configuration:

```hcl
extend_param = {
  DockerLVMConfigOverride = "dockerThinpool=vgpaas/90%VG;kubernetesLV=vgpaas/10%VG;diskType=evs;lvType=linear"
}
```

* `scall_enable` - (Optional, Bool) Specifies whether to enable auto scaling.
  If Autoscaler is enabled, install the autoscaler add-on to use the auto scaling feature.

* `min_node_count` - (Optional, Int) Specifies the minimum number of nodes allowed if auto scaling is enabled.

* `max_node_count` - (Optional, Int) Specifies the maximum number of nodes allowed if auto scaling is enabled.

* `scale_down_cooldown_time` - (Optional, Int) Specifies the time interval between two scaling operations, in minutes.

* `priority` - (Optional, Int) Specifies the weight of the node pool.
  A node pool with a higher weight has a higher priority during scaling.

* `pod_security_groups` - (Optional, List, ForceNew) Specifies the list of security group IDs for the pod.
  Only supported in CCE Turbo clusters of v1.19 and above. Changing this parameter will create a new resource.

* `labels` - (Optional, Map) Specifies the tags of a Kubernetes node, key/value pair format.

* `tags` - (Optional, Map) Specifies the tags of a VM node, key/value pair format.

* `root_volume` - (Required, List, ForceNew) Specifies the configuration of the system disk.
  The structure is described below. Changing this parameter will create a new resource.

* `data_volumes` - (Required, List, ForceNew) Specifies the configuration of the data disks.
  The structure is described below. Changing this parameter will create a new resource.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the CCE node pool. Valid values are
  *prePaid* and *postPaid*, defaults to *postPaid*. Changing this parameter will create a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the CCE node pool.
  Valid values are *month* and *year*. This parameter is mandatory if `charging_mode` is set to *prePaid*.
  Changing this parameter will create a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the CCE node pool. If `period_unit` is set to
  *month*, the value ranges from 1 to 9. If `period_unit` is set to *year*, the value ranges from 1 to 3. This parameter
  is mandatory if `charging_mode` is set to *prePaid*. Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto renew is enabled. Valid values are "true" and "false".
  Changing this parameter will create a new resource.

* `taints` - (Optional, List) Specifies the taints configuration of the nodes to set anti-affinity.
  The structure is described below.

The `root_volume` block supports:

* `size` - (Required, Int, ForceNew) Specifies the disk size in GB. Changing this parameter will create a new resource.

* `volumetype` - (Required, String, ForceNew) Specifies the disk type. Changing this parameter will create a new resource.

* `extend_params` - (Optional, Map, ForceNew) Specifies the disk expansion parameters.
  Changing this parameter will create a new resource.

The `data_volumes` block supports:

* `size` - (Required, Int, ForceNew) Specifies the disk size in GB. Changing this parameter will create a new resource.

* `volumetype` - (Required, String, ForceNew) Specifies the disk type. Changing this parameter will create a new resource.

* `extend_params` - (Optional, Map, ForceNew) Specifies the disk expansion parameters.
  Changing this parameter will create a new resource.

* `kms_key_id` - (Optional, String, ForceNew) Specifies the KMS key ID. This is used to encrypt the volume.
  Changing this parameter will create a new resource.

  -> You need to create an agency (EVSAccessKMS) when disk encryption is used in the current project for the first time ever.

The `taints` block supports:

* `key` - (Required, String) A key must contain 1 to 63 characters starting with a letter or digit. Only letters,
  digits, hyphens (-), underscores (_), and periods (.) are allowed. A DNS subdomain name can be used as the
  prefix of a key.

* `value` - (Required, String) A value must start with a letter or digit and can contain a maximum of 63 characters,
  including letters, digits, hyphens (-), underscores (_), and periods (.).

* `effect` - (Required, String) Available options are NoSchedule, PreferNoSchedule, and NoExecute.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - Node status information.

* `billing_mode` - Billing mode of a node.

* `current_node_count` - The current number of the nodes.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minute.
* `delete` - Default is 20 minute.

## Import

CCE node pool can be imported using the cluster ID and node pool ID separated by a slash, e.g.:

```
$ terraform import huaweicloud_cce_node_pool.my_node_pool 5c20fdad-7288-11eb-b817-0255ac10158b/e9287dff-7288-11eb-b817-0255ac10158b
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`password`, `subnet_id`, `preinstall`, `posteinstall`, `taints`, `initial_node_count` and `pod_security_groups`.
It is generally recommended running `terraform plan` after importing a node pool.
You can then decide if changes should be applied to the node pool, or the resource
definition should be updated to align with the node pool. Also you can ignore changes as below.

```
resource "huaweicloud_cce_node_pool" "my_node_pool" {
    ...

  lifecycle {
    ignore_changes = [
      password, subnet_id,
    ]
  }
}
```
