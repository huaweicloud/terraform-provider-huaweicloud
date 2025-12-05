---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_image_metadata"
description: |-
  Manages an IMS image metadata resource within HuaweiCloud.
---

# huaweicloud_ims_image_metadata

Manages an IMS image metadata resource within HuaweiCloud.

-> 1. This resource is only the metadata for creating the image, and the actual image file corresponding to the image
  does not exist.<br/>2. HuaweiCloud no longer provides Windows operating system type images, and this resource does not
  support the creation of Windows image metadata.

## Example Usage

```hcl
variable "name" {}
variable "os_version" {}

resource "huaweicloud_ims_image_metadata" "test" {
  name         = var.name
  __os_version = var.os_version

  tags = [
    "test_key   = test_value",
    "image_test = image_value"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `__os_version` - (Optional, String, NonUpdatable) Specifies the specific version of the operating system for the
  specified image. If omitted, the default setting is **Other Linux (64 bit)**, and there is no guarantee that this
  image will successfully create a virtual machine or that virtual machines created through this image will function
  properly. For the range of values, please refer to the document link
  [reference](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0910.html#section0).

* `visibility` - (Optional, String, NonUpdatable) Specifies whether other tenants are visible.
  The default value is **private**. When creating image metadata, the valid value can only be **private**.

* `name` - (Optional, String, NonUpdatable) Specifies the name of the image metadata.
  If omitted, it defaults to empty, but creating a virtual opportunity using this image fails. The length ranges from
  `1` to `255` characters. Please refer to the document link for the description of the `name` parameter
  [reference](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0901.html).

* `protected` - (Optional, Bool, NonUpdatable) Specifies whether the image is protected, and the protected image cannot
  be deleted. The default value is **false**.

* `container_format` - (Optional, String, NonUpdatable) Specifies the container format.
  The default value is **bare**.

* `disk_format` - (Optional, String, NonUpdatable) Specifies image format.
  Support **zvhd2**, **vhd**, **raw**, **qcow2**, and **iso**. The default value for non **iso** formats is **zvhd2**.

* `tags` - (Optional, List, NonUpdatable) Specifies image label list. The length ranges from `1` to `255` bits.
  Default is empty. The assignment method for the key in the tag is **"key=value"**.

* `min_ram` - (Optional, Int, NonUpdatable) Specifies the minimum memory required for image operation, in MB.
  Default is `0`. The parameter values are based on the specifications limitations of the cloud server.

* `min_disk` - (Optional, Int, NonUpdatable) Specifies the minimum disk required for image operation, in GB.
  + Linux operating system values range from `10` to `1,024` GB.
  + Windows operating system values range from `20` to `1,024` GB.

  It must be greater than the capacity of the image system disk, otherwise creating a cloud server may fail.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the image metadata.

* `file` - The image file download and upload links.

* `self` - The image link information.

* `schema` - The image view.

* `status` - The image status.  
  The valid values are as follows:
  + **queued**: Indicates that the image metadata has been successfully created and is waiting to upload the image file.
  + **saving**: Indicates that the image is uploading files to the backend storage.
  + **deleted**: Indicates that the image has been deleted.
  + **killed**: Indicates an image upload error.
  + **active**: Indicates that the image can be used normally.

* `max_ram` - The maximum memory supported by the image, measured in MB. The value can refer to the specifications of
  the cloud server and is generally not set.

* `updated_at` - The last update time, in UTC format.

* `__os_bit` - The number of bits in the operating system is usually set to `32` or `64`.

* `__description` - The image description information. Please refer to the document link for parameter specifications
  [reference](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0901.html).

* `__isregistered` - Is it a registered image with a value of **true** or **false**.

* `__platform` - The classification of image platforms, includes **Windows**, **Ubuntu**, **Red Hat**, **SUSE**,
  **CentOS**, **Debian**, **OpenSUSE**, **Oracle Linux**, **Fedora**, **Other**, **CoreOS**, and **Euler OS**.

* `__os_type` - The operating system type, valid values are **Linux**, **Windows**, and **Other**.

* `virtual_env_type` - The image usage environment type.
  + If it is a cloud server image, the value is **FusionCompute**.
  + If it is a data disk image, the value is **DataImage**.
  + If it is a bare metal server image, the value is **Ironic**.
  + If it is an ISO image, the value is **IsoImage**.

* `__image_source_type` - The image backend storage type, currently supports **uds**.

* `__imagetype` - The image type. Currently supporting **gold**, **private**, **shared**, and **market**.

* `created_at` - The creation time, in UTC format.

* `__originalimagename` - The father image ID. Public image or private image created through files, value is empty.

* `__backup_id` - The backup ID. If the image is not created by backup, the value is empty.

* `__productcode` - The product ID of the market image.

* `__image_size` - The size of the image file, in bytes.

* `__data_origin` - The image source. If the image is a public image, value is empty.

* `__lazyloading` - Whether the image supports lazy loading. The value can be **true**, **false**, **True**,
  or **False**.

* `active_at` - The time when the image status became active.

* `__image_displayname` - The name for external display.

* `__os_feature_list` - The additional attributes of the image. The value is a list (in JSON format) of advanced
  features supported by the image.

* `__support_kvm` - Whether the image supports KVM. If yes, the value is **true**.

* `__support_xen` - Whether the image supports Xen. If yes, the value is **true**.

* `__support_largememory` - Whether the image supports large-memory ECSs. If yes, the value is **true**.

* `__support_diskintensive` - Whether the image supports disk-intensive ECSs. If yes, the value is **true**.

* `__support_highperformance` - Whether the image supports high-performance ECSs. If yes, the value is **true**.

* `__support_xen_gpu_type` - Whether the image supports GPU-accelerated ECSs on the Xen platform.
  Please refer to the document link for its value
  [reference](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0910.html#ims_03_0910__table65768383152758).

  This attribute cannot co-exist with `__support_xen` and `__support_kvm`.

* `__support_kvm_gpu_type` - Whether the image supports GPU-accelerated ECSs on the KVM platform.
  Please refer to the document link for its value
  [reference](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0910.html#ims_03_0910__table282523154017).

  This attribute cannot co-exist with `__support_xen` and `__support_kvm`.

* `__support_xen_hana` - Whether the image supports HANA ECSs on the Xen platform. If yes, the value is **true**.
  This attribute cannot co-exist with `__support_xen` and `__support_kvm`.

* `__support_kvm_infiniband` - Whether the image supports ECSs with InfiniBand NICs on the KVM platform. If yes, the
  value is **true**.
  This attribute cannot co-exist with `__support_xen`.

* `__system_support_market` - Whether an image can be published in KooGallery.
  The valid values are as follows:
  + **true**: Support.
  + **false**: Not support.

* `__is_offshelved` - Whether the KooGallery image has been taken offline.
  The valid values are as follows:
  + **true**: Removed.
  + **false**: Not taken down.

* `enterprise_project_id` - The enterprise project ID that the image belongs to.
  `0` or empty indicates that the image belongs to the default enterprise project.

* `__root_origin` - Indicates that the current image source is imported from an external source.
  The value is *file**.

* `__sequence_num` - Indicates the system disk slot position of the cloud server corresponding to the current image.

* `__support_fc_inject` - Whether the image supports password/private key injection using Cloud-Init.
  If the value is **true**, password/private key injection using Cloud-Init is not supported.

  This feature field only applies to ECS system disk images and does not apply to other types of images.

* `hw_firmware_type` - The ECS boot mode.
  The valid values are as follows:
  + **bios**: Indicates the BIOS boot mode.
  + **uefi**: Indicates the UEFI boot mode.

* `hw_vif_multiqueue_enabled` - Whether the image supports NIC multi-queue. The value can be **true** or **false**.

* `__support_arm` - Whether the image uses the Arm architecture. The value can be **true** or **false**.

* `__support_agent_list` - The agents configured for the image.
  The valid values are as follows:
  + **hss**: Server security.
  + **ces**: The host monitoring agent is configured for the image.

  If it is empty, the HSS or host monitoring agents are not configured for the image.

* `__system__cmkid` - The ID of the key used to encrypt the image.

* `__account_code` - The charging identifier for the image.

* `__support_amd` - Whether the image uses AMD's x86 architecture. The value can be **true** or **false**.

* `__support_kvm_hi1822_hisriov` - Whether SR-IOV is supported. If supported, the value is **true**.

* `__support_kvm_hi1822_hivirtionet` - Whether Virtio-Net is supported. If supported, the value is **true**.

* `os_shutdown_timeout` - The timeout duration for a graceful shutdown.
  The value is an integer ranging from `60` to `300`, in seconds. The default value is `60`.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 10 minutes.

## Import

The IMS image metadata resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ims_image_metadata.test <id>
```
