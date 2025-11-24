---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_images"
description: |-
  Use this data source to get the list of HSS container images within HuaweiCloud.
---

# huaweicloud_hss_container_images

Use this data source to get the list of HSS container images within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_images" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `keyword` - (Optional, String) Specifies the search keyword.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of container images.

* `data_list` - The list of container images.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `image_id` - The image unique identifier.

* `image_name` - The image name.

* `image_version` - The image version.

* `create_time` - The image creation time (local save time).
