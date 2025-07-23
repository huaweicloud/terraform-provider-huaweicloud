---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_hda_latest_versions"
description: |-
  Use this data source to get HDA latest versions under a specified region of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_hda_latest_versions

Use this data source to get HDA latest versions under a specified region of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_app_hda_latest_versions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the HDA latest versions are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `hda_latest_versions` - The list of HDA latest versions that matched filter parameters.  
  The [hda_latest_versions](#workspace_app_hda_latest_versions) structure is documented below.

<a name="workspace_app_hda_latest_versions"></a>
The `hda_latest_versions` block supports:

* `latest_version` - The latest version of the HDA.

* `hda_type` - The type of the HDA.
  + **SBC**: Non-VDI SBC type
  + **OA_APP**: VDI non-GPU type
  + **PT_APP**: VDI GPU type
