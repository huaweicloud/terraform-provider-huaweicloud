---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_resource_pool"
description: ""
---

# huaweicloud_modelarts_resource_pool

Manages a ModelArts dedicated resource pool resource within HuaweiCloud.

-> It is recommended to use `huaweicloud_modelartsv2_resource_pool` resource to create resource pool.

~> If you want to expand hyper instance nodes, the provider version must be `1.75.5` or later.

## Example Usage

### create a basic resource pool

```hcl
variable "modelarts_network_id" {}
variable "cce_cluster_id" {}
variable "login_key_pair_name" {}

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "demo"
  description = "This is a demo"
  scope       = ["Train", "Infer", "Notebook"]
  network_id  = var.modelarts_network_id

  resources {
    flavor_id = "modelarts.vm.cpu.16u64g.d"
    count     = 1
    # max_count = 1

    extend_params = jsonencode({
      "dockerBaseSize": "50",
      "runtime": "docker"
    })
    data_volumes {
      volume_type   = "SSD"
      size          = "500Gi"
      extend_params = jsonencode({
        "billing": "on",
        "volumeGroup": "vguser0"
      })
    }
    volume_group_configs {
      lvm_config {
        lv_type = "linear"
      }
      types = [
        "local"
      ]
      volume_group = "vgpaas"
    }
    volume_group_configs {
      lvm_config {
        lv_type = "linear"
        path    = "/data"
      }
      volume_group = "vguser0"
    }
  }

  clusters {
    provider_id = var.cce_cluster_id
  }

  user_login {
    key_pair_name = var.login_key_pair_name
  }
}
```

### create a lite resource pool

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "cce_cluster_id" {}
variable "login_user_password" {}
variable "security_group_ids" {
  type = list(string)
} 

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = "demo"
  prefix      = "test-prefix"
  description = "This is a demo"
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id

  clusters {
    provider_id = var.cce_cluster_id
  }

  user_login {
    password = var.login_user_password
  }

  resources {
    flavor_id          = "modelarts.vm.cpu.8ud"
    count              = 1
    node_pool          = "test-name"
    vpc_id             = var.vpc_id
    subnet_id          = var.subnet_id
    security_group_ids = var.security_group_ids

    labels = {
      aaa = "111"
      bbb = "222"
    }

    tags = {
      key   = "terraform"
      owner = "value"
    }

    taints {
      key    = "key"
      value  = "value"
      effect = "NoSchedule"
    }
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
```

### Prepaid resource pool expansion node (expansion node billing period is one month and auto-renew is disabled)

```hcl
variable "resource_pool_name" {}
variable "subnet_id" {}
variable "vpc_id" {}
variable "description" {}
variable "login_user_password" {}
variable "resource_image_id" {}
variable "resource_flavor_id" {}
variable "resource_provider_id" {}

resource "huaweicloud_modelarts_resource_pool" "test" {
  name        = var.resource_pool_name
  subnet_id   = var.resource_pool_name
  vpc_id      = var.subnet_id
  description = var.description

  metadata {
    annotations = jsonencode({
      "os.modelarts/billing.mode" = "1"
      "os.modelarts/period.type"  = "2"
      "os.modelarts/period.num"   = "1"
      "os.modelarts/auto.renew"   = "0"
      "os.modelarts/auto.pay"     = "1"
    })
  }

  resources {
    flavor_id = var.resource_flavor_id
    count     = 2
    max_count = 2

    os {
      image_id = var.resource_image_id
    }
  }

  clusters {
    provider_id = var.resource_provider_id
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 2
  auto_renew    = false

  user_login {
    password = var.login_user_password
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the dedicated resource pool.  
  The name can contain `4` to `32` characters, only lowercase letters, digits and hyphens (-) are allowed.
  The name must start with a lowercase letter and end with a lowercase letter or digit.

  Changing this parameter will create a new resource.

* `resources` - (Required, List) Specifies the list of resource specifications in the resource pool.  
  Including resource flavors and the number of resources of the corresponding flavors.
  The [resources](#ModelartsResourcePool_ResourceFlavor) structure is documented below.

-> If there are billing nodes in the resource pool, you cannot scale down or delete the node pool through this resource.

* `scope` - (Required, List) Specifies the list of job types supported by the resource pool. It is mandatory when
  `network_id` is specified and can not be specified when `vpc_id` is specified. The options are as follows:
  + **Train**: training job.
  + **Infer**: inference job.
  + **Notebook**: Notebook job.

* `network_id` - (Optional, String, ForceNew) Specifies the ModelArt network ID of the resource pool.

  Changing this parameter will create a new resource.

* `vpc_id` - (Optional, String, ForceNew) Specifies the VPC ID.

  Changing this parameter will create a new resource.

-> **NOTE:** Exactly one of `network_id`, `vpc_id` should be specified.

* `subnet_id` - (Optional, String, ForceNew) Specifies the network ID of a subnet. It is mandatory when `vpc_id` is
  specified.
  
  Changing this parameter will create a new resource.

* `clusters` - (Optional, List, ForceNew) Specifies the list of the CCE clusters. It is mandatory when `vpc_id` is
  specified and can not be specified when `network_id` is specified.
  The [clusters](#ModelartsResourcePool_Clusters) structure is documented below.

  Changing this parameter will create a new resource.

* `user_login` - (Optional, List, ForceNew) Specifies the user login info of the resource pool. It is mandatory when
  `vpc_id` is specified and can not be specified when `network_id` is specified.
  The [user_login](#ModelartsResourcePool_User_login) structure is documented below.

  Changing this parameter will create a new resource.

* `workspace_id` - (Optional, String, ForceNew) Specifies the workspace ID, which defaults to **0**.

  Changing this parameter will create a new resource.

* `prefix` - (Optional, String, ForceNew) Specifies the prefix of the user-defined node name of the resource pool.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the dedicated resource pool.  
  This parameter and `os.modelarts/description` in `metadata.annotations` are set at the same time, the former will
  take precedence.

* `metadata` - (Optional, List) Specifies the metadata of the dedicated resource pool.
  The [metadata](#ModelartsResourcePool_Metadata) structure is documented below.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the resource pool.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.
  Changing this parameter will create a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the resource pool.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the resource pool.
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Valid values are **true** and **false**.

-> 1. `charging_mode`, `period_unit`, `period`, `auto_renew` are mandatory when `vpc_id` is specified.
   <br>2. When creating a resource pool, if the `charging_mode`, `period_unit`, `period` and `auto_renew` parameters
   and subscription parameters in the `metadata.annotations` are specified at the same time, the former will take precedence.

<a name="ModelartsResourcePool_ResourceFlavor"></a>
The `resources` block supports:

* `flavor_id` - (Required, String) Specifies the resource flavor ID.  

* `count` - (Required, Int) Specifies the number of resources of the corresponding flavors.  
  This parameter must be an integer multiple of `resources.creating_step.step`.  
  If you want to expand the nodes, you only need to increase the number of nodes to be expanded
  based on the current parameter value.

  -> If there are billing nodes under the resource pool, use `huaweicloud_modelartsv2_node_batch_delete` to delete
     the on-demand nodes and use `huaweicloud_modelartsv2_node_batch_delete` to unsubscribe the billing nodes.

* `node_pool` - (Optional, String) Specifies the name of resource pool nodes. It can contain `1` to `50`
  characters, and should start with a letter and ending with a letter or digit, only lowercase letters, digits,
  hyphens (-) are allowed, and cannot end with a hyphen (-).

* `max_count` - (Optional, Int) Specifies the max number of resources of the corresponding flavors.
  This parameter must be an integer multiple of `resources.creating_step.step`.

* `vpc_id` - (Optional, String) Specifies the VPC ID. It is mandatory when `resources.subnet_id`,
  `resources.security_group_ids` is specified, and can not be specified when `network_id` is specified.

* `subnet_id` - (Optional, String) Specifies the network ID of a subnet. It is mandatory when
  `resources.security_group_ids`is specified, and can not be specified when `network_id` is specified.

* `security_group_ids` - (Optional, List) Specifies the security group IDs. It can not be specified when `network_id` is
  specified.

* `azs` - (Optional, List) Specifies the AZs for resource pool nodes.
  The [azs](#ModelartsResourcePool_Resources_azs) structure is documented below.

* `labels` - (Optional, Map) Specifies the labels of resource pool nodes.

* `taints` - (Optional, List) Specifies the taints added to nodes. It can not be specified when `network_id` is specified.
  The [taints](#ModelartsResourcePool_Resources_taints) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the resource pool. It can not be specified
  when `network_id` is specified.

* `extend_params` - (Optional, String) Specifies the extend parameters of the resource pool.

* `root_volume` - (Optional, List) Specifies the root volume of the resource pool nodes.  
  The [root_volume](#ModelartsResourcePool_Resources_root_volume) structure is documented below.

* `data_volumes` - (Optional, List) Specifies the data volumes of the resource pool nodes.  
  The [data_volumes](#ModelartsResourcePool_Resources_data_volumes) structure is documented below.

* `volume_group_configs` - (Optional, List) Specifies the extend configurations of the volume groups.  
  The [volume_group_configs](#ModelartsResourcePool_Resources_volume_group_configs) structure is documented below.

* `os` - (Optional, List) Specifies the image information for the specified OS.  
  The [os](#ModelartsResourcePool_Resources_os_info) structure is documented below.

* `driver` - (Optional, List) Specifies the driver information.  
  The [driver](#ModelartsResourcePool_Resources_driver) structure is documented below.

* `creating_step` - (Optional, List) Specifies the creation step configuration of the resource pool nodes.  
  The [creating_step](#ModelartsResourcePool_Resources_creating_step) structure is documented below.

<a name="ModelartsResourcePool_Resources_azs"></a>
The `azs` block supports:

* `az` - (Optional, String) Specifies the AZ name.

* `count` - (Optional, Int) Specifies the number of nodes.

<a name="ModelartsResourcePool_Resources_taints"></a>
The `taints` block supports:

* `key` - (Required, String) Specifies the key of the taint.

* `value` - (Optional, String) Specifies the value of the taint.

* `effect` - (Required, String) Specifies the effect of the taint. Value options: **NoSchedule**, **PreferNoSchedule**,
  **NoExecute**.

<a name="ModelartsResourcePool_Resources_root_volume"></a>
The `root_volume` block supports:

* `volume_type` - (Required, String) Specifies the type of the root volume.  
  The valid values are as follows:
  + **SSD**
  + **GPSSD**
  + **SAS**

* `size` - (Required, String) Specifies the size of the root volume, e.g. **100Gi**.

<a name="ModelartsResourcePool_Resources_data_volumes"></a>
The `data_volumes` block supports:

* `volume_type` - (Required, String) Specifies the type of the data volume.  
  The valid values are as follows:
  + **SSD**
  + **GPSSD**
  + **SAS**

* `size` - (Required, String) Specifies the size of the data volume, e.g. **100Gi**.

* `extend_params` - (Optional, String) Specifies the extend parameters of the data volume.

* `count` - (Optional, Int) Specifies the count of the current data volume configuration.

<a name="ModelartsResourcePool_Resources_volume_group_configs"></a>
The `volume_group_configs` block supports:

* `volume_group` - (Required, String) Specifies the name of the volume group.  
  The valid values are as follows:
  + **vgpaas**
  + **default**
  + **vguser{num}**
  + **vg-everest-localvolume-persistent**
  + **vg-everest-localvolume-ephemeral**

* `docker_thin_pool` - (Optional, Int) Specifies the percentage of container volumes to data volumes on resource pool nodes.

* `lvm_config` - (Optional, List) Specifies the configuration of the LVM management.  
  The [lvm_config](#ModelartsResourcePool_Resources_volume_group_configs_lvm_config) structure is documented below.

* `types` - (Optional, List) Specifies the storage types of the volume group.  
  The valid values for the list elements are as follows:
  + **volume**
  + **local**

<a name="ModelartsResourcePool_Resources_volume_group_configs_lvm_config"></a>
The `lvm_config` block supports:

* `lv_type` - (Required, String) Specifies the LVM write mode.  
  The valid values are as follows:
  + **linear**
  + **striped**

* `path` - (Optional, String) Specifies the volume mount path.

<a name="ModelartsResourcePool_Resources_os_info"></a>
The `os` block supports:

* `name` - (Optional, String) Specifies the OS name of the image.

* `image_id` - (Optional, String) Specifies the image ID.

* `image_type` - (Optional, String) Specifies the image type.

<a name="ModelartsResourcePool_Resources_driver"></a>
The `driver` block supports:

* `version` - (Optional, String) Specifies the driver version.

<a name="ModelartsResourcePool_Resources_creating_step"></a>
The `creating_step` block supports:

* `step` - (Required, Int) Specifies the creation step of the resource pool nodes.

* `type` - (Required, String) Specifies the type of the resource pool nodes.
  The valid values are as follows:
  + **hyperinstance**

<a name="ModelartsResourcePool_Clusters"></a>
The `clusters` block supports:

* `provider_id` - (Required, String, ForceNew) Specifies the ID of the CCE cluster.

<a name="ModelartsResourcePool_User_login"></a>
The `user_login` block supports:

* `password` - (Optional, String, ForceNew) Specifies the password of the login user. The value needs to be salted,
  encrypted and base64 encoded. Default user is **root**.

  Changing this parameter will create a new resource.

* `key_pair_name` - (Optional, String, ForceNew) Specifies key pair name of the login user.

  Changing this parameter will create a new resource.

  -> **NOTE:** Exactly one of `password`, `key_pair_name` should be specified.

<a name="ModelartsResourcePool_Metadata"></a>
The `metadata` block supports:

* `annotations` - (Optional, String) Specifies the annotations of the resource pool, in JSON format.  
  For details, please refer to the [document](https://support.huaweicloud.com/intl/en-us/api-modelarts/CreatePool.html#EN-US_TOPIC_0000001868289874__request_PoolAnnotationsCreation).

  -> 1. This parameter only affects the nodes to be expanded, and will be applied to all nodes that are expanded under
     the `resources` parameter.
     <br>2. The parameter is not allowed to modify the resource pool billing-related parameters.
     <br>3. The `os.modelarts/description` cannot be set at the same time as the external `description` parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the resource pool.

* `resource_pool_id` - Indicates the resource ID of the resource pool.

* `clusters` - Indicates the list of the CCE clusters.
  The [clusters](#ModelartsResourcePool_Clusters) structure is documented below.

<a name="ModelartsResourcePool_Clusters"></a>
The `clusters` block supports:

* `name` - Indicates the name of the CCE cluster.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 90 minutes.
* `update` - Default is 90 minutes.
* `delete` - Default is 30 minutes.

## Import

The ModelArts resource pool can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_resource_pool.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `period_unit`, `period`, `auto_renew`,
`user_login`, `metadata`. It is generally recommended running `terraform plan` after importing a ModelArts
resource pool
You can then decide if changes should be applied to the ModelArts resource pool, or the resource definition should be
updated to align with the ModelArts resource pool. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_modelarts_resource_pool" "resource_pool" {
  ...

  lifecycle {
    ignore_changes = [
      period_unit, period, auto_renew, user_login, metadata
    ]
  }
}
```
