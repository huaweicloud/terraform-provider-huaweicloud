---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_image_snapshot"
description: |-
  Manages a CCI v2 image snapshot resource within HuaweiCloud.
---

# huaweicloud_cciv2_image_snapshot

Manages a CCI v2 image snapshot resource within HuaweiCloud.

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the CCI image snapshot.

* `annotations` - (Optional, Map, NonUpdatable) Specifies the annotations of the CCI image snapshot.

* `labels` - (Optional, Map, NonUpdatable) Specifies the annotations of the CCI image snapshot.

* `finalizers` - (Optional, List, NonUpdatable) Specifies the finalizers of the CCI image snapshot.

* `building_config` - (Optional, List, NonUpdatable) Specifies the building config of the CCI image snapshot.
  The [building_config](#block--building_config) structure is documented below.

* `image_snapshot_size` - (Optional, Int, NonUpdatable) Specifies the size of the CCI image snapshot.

* `images` - (Optional, List, NonUpdatable) The images list of references to images to make image snapshot.
  The [images](#block--images) structure is documented below.

* `registries` - (Optional, List, NonUpdatable) Specifies the registries list.
  The [registries](#block--registries) structure is documented below.

* `ttl_days_after_created` - (Optional, Int, NonUpdatable) The TTL days after created.

<a name="block--building_config"></a>
The `building_config` block supports:

* `auto_create_eip` - (Optional, Bool, NonUpdatable) Specifies whether to auto create EIP.

* `auto_create_eip_attribute` - (Optional, List, NonUpdatable) Specifies whether to auto create EIP.
  The [auto_create_eip_attribute](#block--building_config--auto_create_eip_attribute) structure is documented below.

* `eip_id` - (Optional, String, NonUpdatable) Specifies the EIP ID.

* `namespace` - (Optional, String, NonUpdatable) Specifies the namespace.

<a name="block--building_config--auto_create_eip_attribute"></a>
The `auto_create_eip_attribute` block supports:

* `bandwidth_charge_mode` - (Optional, String, NonUpdatable) Specifies the bandwidth charge mode of EIP.

* `bandwidth_id` - (Optional, String, NonUpdatable) Specifies the ID of EIP.

* `bandwidth_size` - (Optional, Int, NonUpdatable) Specifies the bandwidth size of EIP.

* `ip_version` - (Optional, Int, NonUpdatable) Specifies the IP version used by pod.

* `type` - (Optional, String, NonUpdatable) Specifies the type of EIP.

<a name="block--images"></a>
The `images` block supports:

* `image` - (Optional, String, NonUpdatable) Specifies the name of the reference.

<a name="block--registries"></a>
The `registries` block supports:

* `image_pull_secret` - (Optional, String, NonUpdatable) Specifies the image pull secret.

* `insecure_skip_verify` - (Optional, Bool, NonUpdatable) Specifies whether to allow connections to SSL sites without certs.

* `plain_http` - (Optional, Bool, NonUpdatable) Specifies whether the server uses http protocol.

* `server` - (Optional, String, NonUpdatable) Specifies the image repository server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI image snapshot.

* `creation_timestamp` - The creation timestamp of the CCI image snapshot.

* `kind` - The kind of the CCI image snapshot.

* `resource_version` - The resource version of the CCI image snapshot.

* `status` - The status.
  The [status](#attrblock--status) structure is documented below.

* `uid` - The uid of the CCI image snapshot.

<a name="attrblock--status"></a>
The `status` block supports:

* `expire_date_time` - The expire date time.

* `images` - The list of images.
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

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The CCI v2 ImageSnapshot can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cciv2_image_snapshot.test <id>
```
