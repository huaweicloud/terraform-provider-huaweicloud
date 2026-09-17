---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_download_file"
description: |-
  Manages a resource to download the exported DRS files within HuaweiCloud.
---

# huaweicloud_drs_download_file

Manages a resource to download the exported DRS files within HuaweiCloud.

-> This resource is a one-time action resource used to download the exported DRS files. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
resource "huaweicloud_drs_download_file" "test" {
  files = ["2e69ebb1-f06d-4eb2-a789-2883629jb204_postgresql_20250704024811.xlsx"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to download the DRS files.
  If omitted, the provider-level region will be used.

* `files` - (Required, List) Specifies the list of file names to be downloaded.
  The file will be saved to the current working directory using the first file name in the list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
