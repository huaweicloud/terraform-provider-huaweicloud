---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_download_batch_create_template"
description: |-
  Manages a resource to download the batch import task template within HuaweiCloud.
---

# huaweicloud_drs_download_batch_create_template

Manages a resource to download the batch import task template within HuaweiCloud.

-> 1. This resource is a one-time action resource used to download the batch import task template. Deleting this
  resource will not clear the corresponding request record, but will only remove the resource information from the
  tf state file.
  <br/>2. Executing this resource will generate a file with the suffix **.zip** in the current working directory.

## Example Usage

```hcl
resource "huaweicloud_drs_download_batch_create_template" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `engine_type` - (Optional, String, NonUpdatable) Specifies the database engine.  
  The valid values are as follows:
  + **postgresql**

* `template_file_name` - (Optional, String, NonUpdatable) Specifies the name of the download task template file.
  If omitted, the default file name `drs-batch-create-template.zip` will be used.
  If the file name does not end in **.zip**, **.zip** will be automatically appended.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
