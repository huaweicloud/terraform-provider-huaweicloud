---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flink_job_import"
description: |-
  Use this resource to import Flink jobs within HuaweiCloud.
---

# huaweicloud_dli_flink_job_import

Use this resource to import Flink jobs within HuaweiCloud.

-> This resource is a one-time action resource for importing Flink jobs. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "obs_path" {}

resource "huaweicloud_dli_flink_job_import" "test" {
  obs_path = var.obs_path
  is_cover = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies region where the imported flink jobs are located.
  If omitted, the provider-level region will be used.  
  Changing this will create a new resource.

* `obs_path` - (Required, String, NonUpdatable) Specifies the OBS path for the imported job zip file.  
  This parameter supports entering OBS folder name and file name. Such as:
  + **folder**: **bucket_name/dir1/dir2**
  + **file**: **bucket_name/dir1/dir2/job.zip**

* `is_cover` - (Optional, Bool, NonUpdatable) Specifies whether to overwrite existing jobs with the same name.  
  Defaults to **false**.

  ~> If the value of `obs_path` is a folder name, and multiple files in this folder are exported by the same job.
  Then the API will return an error named **The job name already exists** although `is_cover` set to `true`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
