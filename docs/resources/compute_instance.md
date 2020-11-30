---
subcategory: "Elastic Cloud Server (ECS)"
---

# huaweicloud\_compute\_instance

Manages a ECS VM instance resource within HuaweiCloud.
This is an alternative to `huaweicloud_compute_instance_v2`

## Example Usage

### Basic Instance

```hcl
data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_compute_flavors" "myflavor" {
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "basic" {
  name              = "basic"
  image_id          = data.huaweicloud_images_image.myimage.id
  flavor_id         = data.huaweicloud_compute_flavors.myflavor.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.mynet.id
  }
}
```

### Instance With Associated Eip

```hcl
resource "huaweicloud_compute_instance" "myinstance" {
  name              = "myinstance"
  image_id          = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id         = "s6.small.1"
  key_pair          = "my_key_pair_name"
  security_groups   = ["default"]
  availability_zone = "cn-north-4a"

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }
}

resource "huaweicloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_compute_eip_associate" "associated" {
  public_ip   = huaweicloud_vpc_eip.myeip.address
  instance_id = huaweicloud_compute_instance.myinstance.id
}
```

### Instance With Attached Volume

```hcl
resource "huaweicloud_evs_volume" "myvolume" {
  name              = "myvolume"
  availability_zone = "cn-north-4a"
  volume_type       = "SAS"
  size              = 10
}

resource "huaweicloud_compute_instance" "myinstance" {
  name              = "myinstance"
  image_id          = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id         = "s6.small.1"
  key_pair          = "my_key_pair_name"
  security_groups   = ["default"]
  availability_zone = "cn-north-4a"

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }
}

resource "huaweicloud_compute_volume_attach" "attached" {
  instance_id = huaweicloud_compute_instance.myinstance.id
  volume_id   = huaweicloud_evs_volume.myvolume.id
}
```

### Instance With Multiple Data Disks

It's possible to specify multiple `data_disks` entries to create an instance
with multiple data disks, but we can't ensure the volume attached order. So it's
recommended to use `Instance With Attached Volume` above.

```hcl
resource "huaweicloud_compute_instance" "multi-disk" {
  name              = "multi-net"
  image_id          = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id         = "s6.small.1"
  key_pair          = "my_key_pair_name"
  security_groups   = ["default"]
  availability_zone = "cn-north-4a"

  system_disk_type = "SAS"
  system_disk_size = 40

  data_disks {
    type = "SATA"
    size = "10"
  }
  data_disks {
    type = "SAS"
    size = "20"
  }

  delete_disks_on_termination = true

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }
}
```

### Instance With Multiple Networks

```hcl
resource "huaweicloud_compute_instance" "multi-net" {
  name              = "multi-net"
  image_id          = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id         = "s6.small.1"
  key_pair          = "my_key_pair_name"
  security_groups   = ["default"]
  availability_zone = "cn-north-4a"

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }

  network {
    uuid = "3c4a0d74-24b9-46cf-9d7f-8b7a4dc2f65c"
  }
}
```

### Instance with User Data (cloud-init)

```hcl
resource "huaweicloud_compute_instance" "myinstance" {
  name              = "instance"
  image_id          = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id         = "s6.small.1"
  key_pair          = "my_key_pair_name"
  security_groups   = ["default"]
  availability_zone = "az"
  user_data         = "#cloud-config\nhostname: instance_1.example.com\nfqdn: instance_1.example.com"

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }
}
```

`user_data` can come from a variety of sources: inline, read in from the `file`
function, or the `template_cloudinit_config` resource.

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the instance. If omitted,
	the provider-level region will be used. Changing this creates a new server.

* `name` - (Required, String) A unique name for the resource.

* `image_id` - (Optional, String) Required if `image_name` is empty. The image ID of
    the desired image for the server. Changing this creates a new server.

* `image_name` - (Optional, String) Required if `image_id` is empty. The name of the
    desired image for the server. Changing this creates a new server.

* `flavor_id` - (Optional, String) Required if `flavor_name` is empty. The flavor ID of
    the desired flavor for the server. Changing this resizes the existing server.

* `flavor_name` - (Optional, String) Required if `flavor_id` is empty. The name of the
    desired flavor for the server. Changing this resizes the existing server.

* `user_data` - (Optional, String, ForceNew) The user data to provide when launching the instance.
    Changing this creates a new server.

* `security_groups` - (Optional, String) An array of one or more security group names
    to associate with the server. Changing this results in adding/removing
    security groups from the existing server.

* `availability_zone` - (Required, String, ForceNew) The availability zone in which to create
    the server. Please following [reference](https://developer.huaweicloud.com/endpoint)
    for the values. Changing this creates a new server.

* `network` - (Required, List) An array of one or more networks to attach to the
    instance. The network object structure is documented below. Changing this
    creates a new server.

* `admin_pass` - (Optional, String) The administrative password to assign to the server.
    Changing this changes the root password on the existing server.

* `key_pair` - (Optional, String, ForceNew) The name of a key pair to put on the server. The key
    pair must already be created and associated with the tenant's account.
    Changing this creates a new server.

* `system_disk_type` - (Optional, String) The system disk type of the server. Defaults to `GPSSD`. For details about disk types,
	see [Disk Types and Disk Performance](https://support.huaweicloud.com/en-us/productdesc-evs/en-us_topic_0014580744.html)
    Changing this creates a new server. Available options are:
	* `SSD`: ultra-high I/O disk type.
	* `GPSSD`: general purpose SSD disk type.
	* `SAS`: high I/O disk type.

* `system_disk_size` - (Optional, String) The system disk size in GB, The value range is 1 to 1024. Changing this parameter will update the disk. 
    You can extend the disk by setting this parameter to a new value, which must be between current size and the max size(1024). 
    Shrinking the disk is not supported.

* `data_disks` - (Optional, String) An array of one or more data disks to attach to the
    instance. The data_disks object structure is documented below. Changing this
    creates a new server.

* `tags` - (Optional, Map) Tags key/value pairs to associate with the instance.

* `scheduler_hints` - (Optional, List) Provide the scheduler with hints on how
    the instance should be launched. The available hints are described below.

* `stop_before_destroy` - (Optional, Bool) Whether to try stop instance gracefully
    before destroying it, thus giving chance for guest OS daemons to stop correctly.
    If instance doesn't stop within timeout, it will be destroyed anyway.

* `enterprise_project_id` - (Optional, String) The enterprise project id. Changing this creates a new server.

* `delete_disks_on_termination` - (Optional, Bool) Delete the data disks upon termination of the instance. Defaults to false. Changing this creates a new server.

* `charging_mode` - (Optional, String, ForceNew) The charging mode of the instance. Valid options are: prePaid and postPaid, defaults to postPaid. Changing this creates a new server.

* `period_unit` - (Optional, String, ForceNew) The charging period unit of the instance. Valid options are: month and year, defaults to month. Changing this creates a new server.

* `period` - (Optional, Int, ForceNew) The charging period of the instance. Changing this creates a new server.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto renew is enabled. Changing this creates a new server.

* `user_id` - (Optional, String, ForceNew) User ID, required when using key_pair in prePaid charging mode. Changing this creates a new server.


The `network` block supports:

* `uuid` - (Required, String, ForceNew) The network UUID to
    attach to the server. Changing this creates a new server.

* `fixed_ip_v4` - (Optional, String, ForceNew) Specifies a fixed IPv4 address to be used on this
    network. Changing this creates a new server.

* `access_network` - (Optional, Bool) Specifies if this network should be used for
    provisioning access. Accepts true or false. Defaults to false.

The `scheduler_hints` block supports:

* `group` - (Optional, String, ForceNew) A UUID of a Server Group. The instance will be placed
	into that group.

* `tenancy` - (Optional, String, ForceNew) The tenancy specifies whether the ECS is to be created on a Dedicated Host
	(DeH) or in a shared pool.

* `deh_id` - (Optional, String, ForceNew) The ID of DeH. This parameter takes effect only when the value
	of tenancy is dedicated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `access_ip_v4` - The first detected Fixed IPv4 address _or_ the
    Floating IP.
* `network/fixed_ip_v4` - The Fixed IPv4 address of the Instance on that network.
* `network/mac` - The MAC address of the NIC on that network.
* `network/port` - The port ID corresponding to the IP address on that network.
* `volume_attached/volume_id` - The volume id on that attachment.
* `volume_attached/pci_address` - The volume pci address on that attachment.
* `volume_attached/boot_index` - The volume boot index on that attachment.
* `volume_attached/size` - The volume size on that attachment.
* `system_disk_id` - The system disk voume ID.


## Import

Instances can be imported by their `id`. For example,
```
terraform import huaweicloud_compute_instance.my_instance b11b407c-e604-4e8d-8bc4-92398320b847
```
Note that the imported state may not be identical to your resource definition, which 
could be because of a different network interface attachment order, missing ephemeral
disk configuration, or some other reason. It is generally recommended running 
`terraform plan` after importing an instance. You can then decide if changes should
be applied to the instance, or the resource definition should be updated to align
with the instance. 

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 30 minute.
- `update` - Default is 30 minute.
- `delete` - Default is 30 minute.

