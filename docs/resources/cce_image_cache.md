---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_image_cache"
description: |-
  Manages a CCE image cache resource within huaweicloud.
---

# huaweicloud_cce_image_cache

Manages a CCE image cache resource within huaweicloud.

## Example Usage

### Basic Usage

```hcl
variable "name" {}

resource "huaweicloud_cce_image_cache" "test" {
  name   = var.name
  images = ["busybox:latest"]

  building_config {
    cluster            = huaweicloud_cce_cluster.test.id
    image_pull_secrets = ["default:default-secret"]
  }

  image_cache_size = 20
  retention_days   = 7
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE image cache resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new image cache resource.

* `name` - (Required, String) Specifies the image cache name.

* `images` - (Required, List) Specifies the list of container images in an image cache.

* `building_config` - (Required, List) Specifies configuration for creating an image cache.
  The [structure](#cce_image_cache_building_config) is documented below.

* `image_cache_size` - (Optional, Int) Specifies the size of the disk for an image cache, in GiB.
  The cached objects are the decompressed image files. The image cache size should be greater than or
  equal to three times the total size of all container images in the image cache.

* `retention_days` - (Optional, Int) Specifies the validity period of an image cache.
  If this parameter is not specified or the value is 0, the image cache is permanently valid.
  After the validity period expires, the image cache will expire automatically and be deleted.

<a name="cce_image_cache_building_config"></a>
The `building_config` block supports:

* `cluster` - (Required, String) Specifies the ID of a CCE Autopilot cluster where a temporary pod
  is started for creating an image cache.

* `image_pull_secrets` - (Optional, List) Specifies the list of access credentials for downloading
  the images to be cached. If no access credential is specified or no valid credential is available,
  only public images can be downloaded.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the image cache resource.
  
* `created_at` - The time when the image cache was created.

* `status` - The status of the image cache.

## Import

The image cache can be imported using the image cache ID, e.g.

```bash
 $ terraform import huaweicloud_cce_image_cache.my_cache <cache_id>
```
