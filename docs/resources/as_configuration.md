---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_configuration"
description: |-
  Manages an AS configuration resource within HuaweiCloud.
---

# huaweicloud_as_configuration

Manages an AS configuration resource within HuaweiCloud.

## Example Usage

### Basic AS Configuration

```hcl
variable "flavor_id" {}
variable "image_id" {}
variable "ssh_key" {}
variable "security_group_id" {}

resource "huaweicloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor             = var.flavor_id
    image              = var.image_id
    key_name           = var.ssh_key
    security_group_ids = [var.security_group_id]

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
  }
}
```

### AS Configuration With Encrypted Data Disk

```hcl
variable "flavor_id" {}
variable "image_id" {}
variable "ssh_key" {}
variable "kms_id" {}
variable "security_group_id" {}

resource "huaweicloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor             = var.flavor_id
    image              = var.image_id
    key_name           = var.ssh_key
    security_group_ids = [var.security_group_id]

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }

    disk {
      size        = 100
      volume_type = "SSD"
      disk_type   = "DATA"
      kms_id      = var.kms_id
    }
  }
}
```

### AS Configuration With User Data and Metadata

```hcl
variable "flavor_id" {}
variable "image_id" {}
variable "ssh_key" {}
variable "security_group_id" {}

resource "huaweicloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor             = var.flavor_id
    image              = var.image_id
    key_name           = var.ssh_key
    security_group_ids = [var.security_group_id]
    user_data          = file("userdata.txt")

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }

    metadata  = {
      some_key = "some_value"
    }
  }
}
```

### AS Configuration uses password authentication for Linux ECS

```hcl
variable "flavor_id" {}
variable "ecs_image_id" {}
variable "security_group_id" {}

resource "huaweicloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor             = var.flavor_id
    image              = var.ecs_image_id
    security_group_ids = [var.security_group_id]

    user_data = <<EOT
#! /bin/bash
echo 'root:$6$V6azyeLwcD3CHlpY$BN3VVq18fmCkj66B4zdHLWevqcxlig' | chpasswd -e
EOT

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
  }
}
```

### AS Configuration uses password authentication for Windows ECS

```hcl
variable "flavor_id" {}
variable "windows_image_id" {}
variable "security_group_id" {}
variable "admin_pass" {}

resource "huaweicloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor             = var.flavor_id
    image              = var.windows_image_id
    security_group_ids = [var.security_group_id]
    admin_pass         = var.admin_pass

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
  }
}
```

### AS Configuration uses the existing instance specifications as the template

```hcl
variable "instance_id" {}
variable "ssh_key" {}
variable "security_group_id" {}

resource "huaweicloud_as_configuration" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    instance_id        = var.instance_id
    key_name           = var.ssh_key
    security_group_ids = [var.security_group_id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the AS configuration.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `scaling_configuration_name` - (Required, String, ForceNew) Specifies the AS configuration name.
  The name contains only letters, digits, underscores (_), and hyphens (-), and cannot exceed `64` characters.
  Changing this will create a new resource.

* `instance_config` - (Required, List, ForceNew) Specifies the information about instance configuration.
  The [instance_config](#as_instance_config) structure is documented below.
  Changing this will create a new resource.

<a name="as_instance_config"></a>
The `instance_config` block supports:

* `instance_id` - (Optional, String, ForceNew) Specifies the ECS instance ID when using its specification
  as the template to create AS configurations. In this case, `flavor`, `image`, `disk`, `security_group_ids`, `tenancy`
  and `dedicated_host_id` arguments do not take effect.
  If this argument is not specified, `flavor`, `image`, and `disk` arguments are mandatory.
  Changing this will create a new resource.

* `flavor` - (Optional, String, ForceNew) Specifies the ECS flavor name. A maximum of `10` flavors can be selected.
  Use a comma (,) to separate multiple flavor names. Changing this will create a new resource.

* `image` - (Optional, String, ForceNew) Specifies the ECS image ID. Changing this will create a new resource.

* `disk` - (Optional, List, ForceNew) Specifies the disk group information. System disks are mandatory and
  data disks are optional. The [disk](#instance_config_disk_object) structure is documented below.
  Changing this will create a new resource.

* `security_group_ids` - (Required, List, ForceNew) Specifies an array of one or more security group IDs.
  Changing this will create a new resource.

* `charging_mode` - (Optional, String, ForceNew) Specifies a billing mode for an ECS.
  The value can be `postPaid` and `spot`. The default value is `postPaid`.
  Changing this will create a new resource.

* `flavor_priority_policy` - (Optional, String, ForceNew) Specifies the priority policy used when there are multiple flavors
  and instances to be created using an AS configuration. The value can be `PICK_FIRST` and `COST_FIRST`.

  + **PICK_FIRST** (default): When an ECS is added for capacity expansion, the target flavor is determined in the order
    in the flavor list.
  + **COST_FIRST**: When an ECS is added for capacity expansion, the target flavor is determined for minimal expenses.

  Changing this will create a new resource.

* `ecs_group_id` - (Optional, String, ForceNew) Specifies the ECS group ID. Changing this will create a new resource.

  -> To ensure service reliability, an ECS group allows ECSs within in the group to be automatically allocated to
  different hosts.

* `tenancy` - (Optional, String, ForceNew) Configure this field to **dedicated** to create ECS instances on DeHs.
  Before configuring this field, prepare DeHs. Changing this will create a new resource.

* `dedicated_host_id` - (Optional, String, ForceNew) Specifies the ID of the DEH.
  Changing this will create a new resource.

  -> This parameter is valid only when `tenancy` is set to **dedicated**.
  <br/>If this parameter is specified, ECSs will be created on a specified DeH.
  <br/>If this parameter is not specified, the system automatically selects the DeH with the maximum available memory
  size from the DeHs that meet specifications requirements to create the ECSs, thereby balancing load of the DeHs.

* `public_ip` - (Optional, List, ForceNew) Specifies the EIP of the ECS instance.
  The [public_ip](#instance_config_public_ip_object) structure is documented below.
  Changing this will create a new resource.

* `key_name` - (Optional, String, ForceNew) Specifies the name of the SSH key pair used to log in to the instance.
  Changing this will create a new resource.

* `user_data` - (Optional, String, ForceNew) Specifies the user data to be injected during the ECS creation process.
  Changing this will create a new resource. For more information, see
  [Passing User Data to ECSs](https://support.huaweicloud.com/intl/en-us/usermanual-ecs/en-us_topic_0032380449.html).

  -> 1. The content to be injected must be encoded with base64. The maximum size of the content to be injected
  (before encoding) is `32` KB.
  <br/>2. If `key_name` is not specified, the data injected by `user_data` is the password of user `root` for logging in
  to the ECS by default.
  <br/>3. If both `key_name` and `user_data` are specified, `user_data` only injects user data.
  <br/>4. This parameter is mandatory when you create a Linux ECS using the password authentication mode. Its value is
  the initial user `root` password.
  <br/>5. When the value of this field is used as a password, the recommended complexity for the password is as follows:
  (1) The value ranges from `8` to `26` characters. (2) The value contains at least three of the following character
  types: uppercase letters, lowercase letters, digits, and special characters `!@$%^-_=+[{}]:,./?`.

* `admin_pass` - (Optional, String, ForceNew) Specifies the initial login password of the administrator account for
  logging in to an ECS using password authentication. The Windows administrator is `Administrator`.

  -> Password complexity requirements:
  <br/>1. Consists of `8` to `26` characters.
  <br/>2. Contains at least three of the following character types: uppercase letters, lowercase letters, digits, and
  special characters `!@$%^-_=+[{}]:,./?`.
  <br/>3. The password cannot contain the username or the username in reversed order.
  <br/>4. The Windows ECS password cannot contain the username, the username in reversed order, or more than two
  consecutive characters in the username.

~> Field `admin_pass` is used for Windows system password authentication, and `user_data` is used for Linux system
password authentication.

* `metadata` - (Optional, Map, ForceNew) Specifies the key/value pairs to make available from within the instance.
  Changing this will create a new resource.

* `personality` - (Optional, List, ForceNew) Specifies the customize personality of an instance by defining one or
  more files and their contents. The [personality](#instance_config_personality_object) structure is documented below.
  Changing this will create a new resource.

<a name="instance_config_disk_object"></a>
The `disk` block supports:

* `size` - (Required, Int, ForceNew) Specifies the disk size. The unit is GB.
  The system disk size ranges from `1` to `1024`, and not less than the minimum value of the system disk in the
  instance image. The data disk size ranges from `10` to `32,768`.
  Changing this will create a new resource.

* `volume_type` - (Required, String, ForceNew) Specifies the disk type. Changing this will create a new resource.
  Available options are:
  + **SSD**: The ultra-high I/O type.
  + **SAS**: The high I/O EVS type.
  + **SATA**: The common I/O type.
  + **GPSSD**: The general purpose SSD type.
  + **ESSD**: The extreme SSD type.
  + **GPSSD2**: The general purpose SSD V2 type.
  + **ESSD2**: The extreme SSD V2 type.

  -> Different ECS flavors support different disk types. For details about disk types, see
  [Disk Types and Performance](https://support.huaweicloud.com/intl/en-us/productdesc-evs/en-us_topic_0014580744.html).

* `disk_type` - (Required, String, ForceNew) Specifies whether the disk is a system disk or a data disk.
  Option **DATA** indicates a data disk, option **SYS** indicates a system disk.
  Changing this will create a new resource.

* `kms_id` - (Optional, String, ForceNew) Specifies the encryption KMS ID of the **DATA** disk.
  Changing this will create a new resource.

* `iops` - (Optional, Int, ForceNew) Specifies the IOPS configured for an EVS disk.
  Changing this will create a new resource.

  -> This parameter is mandatory only when `volume_type` is set to **GPSSD2** or **ESSD2**.
  <br/>For details about IOPS of GPSSD2 and ESSD2 EVS disks, see
  [Disk Types and Performance](https://support.huaweicloud.com/intl/en-us/productdesc-evs/en-us_topic_0014580744.html).
  <br/>Only pay-per-use billing is supported currently.

* `throughput` - (Optional, Int, ForceNew) Specifies the throughput of an EVS disk. The unit is MiB/s.
  Changing this will create a new resource.

  -> This parameter is mandatory only when `volume_type` is set to **GPSSD2** and cannot be configured
  when `volume_type` is set to other values.
  <br/>For details about the throughput range of GPSSD2 EVS disks, see
  [Disk Types and Performance](https://support.huaweicloud.com/intl/en-us/productdesc-evs/en-us_topic_0014580744.html).
  <br/>Only pay-per-use billing is supported currently.

* `dedicated_storage_id` - (Optional, String, ForceNew) Specifies a DSS device ID for creating an ECS disk.

  -> Specify DSS devices for all disks in an AS configuration or not. If DSS devices are specified, all the
  data stores must belong to the same AZ, and the disk types supported by a DSS device for a disk must be
  the same as the `volume_type` value.

* `data_disk_image_id` - (Optional, String, ForceNew) Specifies the ID of a data disk image used to export data disks of
  an ECS.

* `snapshot_id` - (Optional, String, ForceNew) Specifies the disk backup snapshot ID for restoring the system disk and
  data disks using a full-ECS backup when a full-ECS image is used.

  -> You can obtain the disk backup snapshot ID using the full-ECS backup ID in
  [Querying a Single Backup](https://support.huaweicloud.com/intl/en-us/api-csbs/en-us_topic_0059304234.html).
  <br/>Each disk in an AS configuration must correspond to a disk backup in the full-ECS backup by `snapshot_id`.

<a name="instance_config_public_ip_object"></a>
The `public_ip` block supports:

* `eip` - (Required, List, ForceNew) Specifies the EIP configuration that will be automatically assigned to the instance.
  The object structure is documented below. Changing this will create a new resource.

The `eip` block supports:

* `ip_type` - (Required, String, ForceNew) Specifies the EIP type. Possible values are **5_bgp** (dynamic BGP)
  and **5_sbgp** (static BGP). Changing this will create a new resource.

* `bandwidth` - (Required, List, ForceNew) Specifies the bandwidth information. The object structure is documented below.
  Changing this will create a new resource.

The `bandwidth` block supports:

* `share_type` - (Required, String, ForceNew) Specifies the bandwidth sharing type.
  The value can be **PER** (exclusive bandwidth) or **WHOLE** (shared bandwidth).
  Changing this will create a new resource.

  -> If `share_type` is set to **PER**, the parameter `charging_mode` and `size` are mandatory, the parameter `id`
  is invalid.
  If `share_type` is set to **WHOLE**, the parameter `id` is mandatory, the parameter `charging_mode` and `size`
  are invalid.

* `charging_mode` - (Optional, String, ForceNew) Specifies the bandwidth billing type.
  Changing this creates a new resource. The valid values are as follows:
  + **bandwidth**: Billing by bandwidth.
  + **traffic**: Billing by traffic.

* `size` - (Optional, Int, ForceNew) Specifies the bandwidth (Mbit/s). The value range for bandwidth billed by bandwidth
  is `1` to `2,000` and that for bandwidth billed by traffic is `1` to `300`.
  Changing this creates a new resource.

* `id` - (Optional, String, ForceNew) Specifies the ID of the shared bandwidth.
  Changing this will create a new resource.

<a name="instance_config_personality_object"></a>
The `personality` block supports:

* `path` - (Required, String, ForceNew) Specifies the path of the injected file. Changing this creates a new resource.
  + For Linux OSs, specify the path, for example, **/etc/foo.txt**, for storing the injected file.
  + For Windows, the injected file is automatically stored in the root directory of drive `C`. You only need to specify
    the file name, for example, **foo**. The file name contains only letters and digits.

* `content` - (Required, String, ForceNew) Specifies the content of the injected file, which must be encoded with base64.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The AS configuration status, the value can be **Bound** or **Unbound**.

* `create_time` - The creation time of the AS configuration, in UTC format.

* `instance_config` - The list of information about instance configurations.
  The [instance_config](#as_instance_config_attr) structure is documented below.

<a name="as_instance_config_attr"></a>
The `instance_config` block supports:

* `key_fingerprint` - The fingerprint of the SSH key pair used to log in to the instance.

## Import

AS configurations can be imported by their `id`, e.g.

```bash
$ terraform import huaweicloud_as_configuration.test <id>
```

Note that the imported state may not be identical to your resource definition, due to `instance_config.0.instance_id`,
`instance_config.0.admin_pass`, `instance_config.0.user_data`, and `instance_config.0.metadata` are missing from the
API response. You can ignore changes after importing an AS configuration as below.

```hcl
resource "huaweicloud_as_configuration" "test" {
  ...

  lifecycle {
    ignore_changes = [
      instance_config.0.instance_id,
      instance_config.0.admin_pass,
      instance_config.0.user_data,
      instance_config.0.metadata,
    ]
  }
}
```
