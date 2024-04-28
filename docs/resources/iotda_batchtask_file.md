---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_batchtask_file"
description: ""
---

# huaweicloud_iotda_batchtask_file

!> **WARNING:** It has been deprecated.

Manages an IoTDA batch task file within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_iotda_batchtask_file" "test" {
  content = "./BatchCreateDevices_Template.xlsx"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA batch task file resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `content` - (Required, String, ForceNew) Specifies the path to the batch task file to be uploaded.
  Currently, only the **xlsx/xls** file format is supported, and the maximum number of lines in the file is `30000`.
  The file name cannot be duplicated. Changing this parameter will create a new resource.
  Please following [reference](https://support.huaweicloud.com/intl/en-us/usermanual-iothub/iot_01_0032.html),
  download the template file and fill it out.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `name` - The resource name.

* `created_at` - The time of file upload. The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.

## Import

The batch task file can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_batchtask_file.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `content`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_iotda_batchtask_file" "test" { 
  ...
  
  lifecycle {
    ignore_changes = [
      content,
    ]
  }
}
```
