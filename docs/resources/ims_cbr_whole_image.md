---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_cbr_whole_image"
description: |-
  Manages an IMS whole image resource created from CBR backup within HuaweiCloud.
---

# huaweicloud_ims_cbr_whole_image

Manages an IMS whole image resource created from CBR backup within HuaweiCloud.

-> After deleting the image, the CBR backup that has not been deleted will be retained and continue to be charged.
   If you need to delete it later, you can delete the corresponding CBR backup in the CBR backup console.

## Example Usage

### Creating an IMS whole image from CBR backup

```hcl
variable "name" {}
variable "backup_id" {}

resource "huaweicloud_ims_cbr_whole_image" "test" {
  name      = var.name
  backup_id = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the image.
  The valid length is limited from `1` to `128` characters.
  The first and last letters of the name cannot be spaces.
  The name can contain uppercase letters, lowercase letters, numbers, spaces, chinese, and special characters (-._).

* `backup_id` - (Required, String, ForceNew) Specifies the CBR instance backup ID used to create the image.
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the image.

* `max_ram` - (Optional, Int) Specifies the maximum memory of the image, in MB unit.

* `min_ram` - (Optional, Int) Specifies the minimum memory of the image, in MB unit.
  The default value is `0`, indicating that the memory is not restricted.

* `is_delete_backup` - (Optional, Bool) Specifies whether to delete the associated CBR backup when deleting image.
  The value can be **true** or **false**.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the image.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the IMS image belongs.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the image.

* `file` - The image file download and upload links.

* `self` - The image link information.

* `schema` - The image view.

* `status` - The image status.  
  The valid values are as follows:
  + **queued**: Indicates that the image has been successfully created and is waiting to upload the image file.
  + **saving**: Indicates that the image is uploading files to the backend storage.
  + **deleted**: Indicates that the image has been deleted.
  + **killed**: Indicates an image upload error.
  + **active**: Indicates that the image can be used normally.

* `visibility` - Whether other tenants are visible.  
  The valid values are as follows:
  + **private**: Indicates private image.
  + **public**: Indicates public image.
  + **shared**: Indicates shared image.

* `protected` - Whether the image is protected, and the protected image cannot be deleted. The value can be **true**
  or **false**.

* `container_format` - The container format.

* `updated_at` - The last update time, in UTC format.

* `__os_bit` - The number of bits in the operating system is usually set to `32` or `64`.

* `os_version` - The operating system version of the image.

* `disk_format` - The image format. The value can be **zvhd2**, **vhd**, **zvhd**, **raw**, or **qcow2**.

* `__isregistered` - Is it a registered image with a value of **true** or **false**.

* `__platform` - The classification of image platforms, includes **Windows**, **Ubuntu**, **Red Hat**, **SUSE**,
  **CentOS**, **Debian**, **OpenSUSE**, **Oracle Linux**, **Fedora**, **Other**, **CoreOS**, and **Euler OS**.

* `__os_type` - The operating system type, valid values are **Linux**, **Windows**, and **Other**.

* `min_disk` - The minimum disk space required to run an image, in GB unit.

* `virtual_env_type` - The image usage environment type.
  + If it is a cloud server image, the value is **FusionCompute**.
  + If it is a data disk image, the value is **DataImage**.
  + If it is a bare metal server image, the value is **Ironic**.
  + If it is an ISO image, the value is **IsoImage**.

* `__image_source_type` - The image backend storage type, currently supports **uds**.

* `__imagetype` - The image type. Currently supporting **gold**, **private**, **shared**, and **market**.

* `created_at` - The creation time, in UTC format.

* `__originalimagename` - The father image ID. Public image or private image created through files, value is empty.

* `__productcode` - The product ID of the market image.

* `__image_size` - The size of the image file, in bytes unit.

* `data_origin` - The image source. The format is **server_backup,backup_id**.

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
  value is **true**. This attribute cannot co-exist with `__support_xen`.

* `__system_support_market` - Whether an image can be published in KooGallery.
  The valid values are as follows:
  + **true**: Support.
  + **false**: Not support.

* `__is_offshelved` - Whether the KooGallery image has been taken offline.
  The valid values are as follows:
  + **true**: Removed.
  + **false**: Not taken down.

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

* `create` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

The IMS CBR whole image resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ims_cbr_whole_image.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `is_delete_backup`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the image. Also, you can ignore
changes as below.

```
resource "huaweicloud_ims_cbr_whole_image" "test" {
    ...
  lifecycle {
    ignore_changes = [
      is_delete_backup,
    ]
  }
}
```
