---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_cluster_protect_protection_items"
description: |-
  Use this data source to get the list of HSS cluster protect all protection items within HuaweiCloud.
---

# huaweicloud_hss_cluster_protect_protection_items

Use this data source to get the list of HSS cluster protect protection items within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_cluster_protect_protection_items" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `vuls` - The list of vulnerability information.

* `baselines` - The list of baseline information.

  The [baselines](#baselines_struct) structure is documented below.

* `malwares` - The list of malicious file information.

  The [malwares](#malwares_struct) structure is documented below.

* `images` - The list of image information.

  The [images](#images_struct) structure is documented below.

* `clusters` - The list of cluster information.

  The [clusters](#clusters_struct) structure is documented below.

<a name="baselines_struct"></a>
The `baselines` block supports:

* `baseline_desc` - The baseline description.

* `baseline_index` - The baseline ID.

* `baseline_type` - The baseline type.

<a name="malwares_struct"></a>
The `malwares` block supports:

* `malware_type` - The malicious file type.

<a name="images_struct"></a>
The `images` block supports:

* `image_name` - The image name.

* `image_version` - The image version.

* `id` - The ID.

<a name="clusters_struct"></a>
The `clusters` block supports:

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `cluster_ns` - The cluster namespace list.

* `cluster_labels` - The cluster label list.

* `protect_status` - The protection status.  
  The valid values are as follows:
  + **unprotected**: Unprotected.
  + **plugin error**: Plugin error.
  + **protected with policy**: Protected with policy.
  + **deploy policy failed**: Deploy policy failed.
  + **protected without policy**: Protected without policy.
  + **uninstall failed**: Uninstall failed.
  + **uninstall**: Uninstall.
  + **installing**: Installing.
