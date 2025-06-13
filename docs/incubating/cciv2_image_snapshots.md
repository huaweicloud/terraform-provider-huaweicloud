---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_image_snapshots"
description: |-
  Use this data source to get the list of CCI image snapshots within HuaweiCloud.
---

# huaweicloud_cciv2_image_snapshots

Use this data source to get the list of CCI image snapshots within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_cciv2_image_snapshots" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the CCI image snapshot.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `image_snapshots` - The image snapshots.
  The [image_snapshots](#image_snapshots) structure is documented below.

<a name="image_snapshots"></a>
The `image_snapshots` block supports:

* `annotations` - The annotations of the CCI image snapshot.

* `building_config` - The building config of the CCI image snapshot.
  The [building_config](#image_snapshots_building_config) structure is documented below.

* `creation_timestamp` - The creation timestamp of the CCI image snapshot.

* `finalizers` - The finalizers of the CCI image snapshot.

* `image_snapshot_size` - The size of the CCI image snapshot.

* `images` - The images list of references to images to make image snapshot.
  The [images](#image_snapshots_images) structure is documented below.

* `labels` - The labels of the CCI image snapshot.

* `registries` - The registries list.
  The [registries](#image_snapshots_registries) structure is documented below.

* `resource_version` - The resource version of the CCI image snapshot.

* `status` - The status.
  The [status](#image_snapshots_status) structure is documented below.

* `ttl_days_after_created` - The TTL days after created.

* `uid` - The uid of the CCI image snapshot.

<a name="image_snapshots_building_config"></a>
The `building_config` block supports:

* `auto_create_eip` - Whether to auto create EIP.

* `auto_create_eip_attribute` - The auto create EIP attribute.
  The [auto_create_eip_attribute](#image_snapshots_building_config_auto_create_eip_attribute) structure is documented below.

* `eip_id` - The EIP ID.
* `namespace` - The namespace.

<a name="image_snapshots_building_config_auto_create_eip_attribute"></a>
The `auto_create_eip_attribute` block supports:

* `bandwidth_charge_mode` - The bandwidth charge mode of EIP.

* `bandwidth_id` - The bandwidth ID of EIP.

* `bandwidth_size` - The bandwidth size of EIP.

* `ip_version` - The IP version used by pod.

* `type` - The type of EIP.

<a name="image_snapshots_images"></a>
The `images` block supports:

* `image` - The image name.

<a name="image_snapshots_registries"></a>
The `registries` block supports:

* `image_pull_secret` - The image pull secret.

* `insecure_skip_verify` - Whether to allow connections to SSL sites without certs.

* `plain_http` - Whether the server uses http protocol.

* `server` - The image repository server.

<a name="image_snapshots_status"></a>
The `status` block supports:

* `expire_date_time` - The expire date time.

* `images` - The list of images.
  The [images](#image_snapshots_status_images) structure is documented below.

* `last_updated_time` - The last updated time.

* `message` - The message.

* `phase` - The phase.

* `reason` - The reason.

* `snapshot_id` - The snapshot ID.

* `snapshot_name` - The snapshot name.

<a name="image_snapshots_status_images"></a>
The `images` block supports:

* `digest` - The image digest.

* `image` - The image name.

* `size_bytes` - The size of the image in bytes.
