---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_image_snapshot"
description: |-
  Manages a CCI v2 ImageSnapshot resource within HuaweiCloud.
---

# huaweicloud_cciv2_image_snapshot

Manages a CCI v2 ImageSnapshot resource within HuaweiCloud.

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the CCI Image Snapshot.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI Image Snapshot.

* `building_config` - (Optional, List) Specifies the building config of the CCI Image Snapshot.
  The [building_config](#block--building_config) structure is documented below.

* `image_snapshot_size` - (Optional, Float) Specifies the size of the CCI Image Snapshot.

* `images` - (Optional, List) The images list of references to images to make image snapshot.
  The [images](#block--images) structure is documented below.

* `labels` - (Optional, Map) Specifies the annotations of the CCI Image Snapshot.

* `registries` - (Optional, List) Specifies the registries list.
  The [registries](#block--registries) structure is documented below.

* `ttl_days_after_created` - (Optional, Int) The TTL days after created.

<a name="block--building_config"></a>
The `building_config` block supports:

* `auto_create_eip` - (Optional, Bool) Specifies whether to auto create EIP.

* `auto_create_eip_attribute` - (Optional, List) Specifies whether to auto create EIP.
  The [auto_create_eip_attribute](#block--building_config--auto_create_eip_attribute) structure is documented below.

* `eip_id` - (Optional, String) Specifies the EIP ID.

* `namespace` - (Optional, String) Specifies the namespace.

<a name="block--building_config--auto_create_eip_attribute"></a>
The `auto_create_eip_attribute` block supports:

* `bandwidth_charge_mode` - (Optional, String) Specifies the bandwidth charge mode of EIP.

* `bandwidth_id` - (Optional, String) Specifies the ID of EIP.

* `bandwidth_size` - (Optional, Int) Specifies the bandwidth size of EIP.

* `ip_version` - (Optional, Int) Specifies the IP version used by pod.

* `type` - (Optional, String) Specifies the type of EIP.

<a name="block--images"></a>
The `images` block supports:

<a name="block--registries"></a>
The `registries` block supports:

* `image_pull_secret` - (Optional, String) Specifies the image pull secret.

* `insecure_skip_verify` - (Optional, Bool) Specifies whether to allow connections to SSL sites without certs.

* `plain_http` - (Optional, Bool) Specifies whether the server uses http protocol.

* `server` - (Optional, String) Specifies the image repository server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI Image Snapshot.

* `creation_timestamp` - The creation timestamp of the CCI Image Snapshot.

* `images` - The images list of references to images to make image snapshot.
  The [images](#attrblock--images) structure is documented below.

* `kind` - The kind of the CCI Image Snapshot.

* `resource_version` - The resource version of the CCI Image Snapshot.

* `status` - The status.
  The [status](#attrblock--status) structure is documented below.

* `uid` - The uid of the CCI Image Snapshot.

<a name="attrblock--images"></a>
The `images` block supports:

* `image` - Specifies the image name.

<a name="attrblock--status"></a>
The `status` block supports:

* `expire_date_time` - The expire date time.

* `images` - The status.
  The [images](#attrblock--status--images) structure is documented below.

* `last_updated_time` - The last updated time.

* `message` - The message.

* `phase` - The phase.

* `reason` - The reason.

* `snapshot_id` - The snapshot ID.

* `snapshot_name` - The snapshot name.

<a name="attrblock--status--images"></a>
The `images` block supports:

* `digest` - The image digest.

* `image` - The image name.

* `size_bytes` - The size of the image in bytes.

## Import

The CCI v2 ImageSnapshot can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cciv2_image_snapshot.test <id>
```
