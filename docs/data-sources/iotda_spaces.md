---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_spaces"
description: ""
---

# huaweicloud_iotda_spaces

Use this data source to get the list of IoTDA spaces within HuaweiCloud.

-> When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify the IoTDA service
  endpoint in `provider` block.
  You can login to the IoTDA console, choose the instance **Overview** and click **Access Details**
  to view the HTTPS application access address. An example of the access address might be
  **9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com**, then you need to configure the
  `provider` block as follows:

  ```hcl
  provider "huaweicloud" {
    endpoints = {
      iotda = "https://9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com"
    }
  }
  ```

## Example Usage

```hcl
variable "space_id" {}

data "huaweicloud_iotda_spaces" "test" {
  space_id = var.space_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the spaces.
  If omitted, the provider-level region will be used.

* `space_id` - (Optional, String) Specifies the ID of the space to be queried.

* `space_name` - (Optional, String) Specifies the name of the space to be queried.

* `is_default` - (Optional, String) Specifies whether to query the default space.
  The valid values are as follows:
  + **true**: Query the default space.
  + **false**: Query all non default spaces.
  If omitted, query all spaces under the current instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `spaces` - All spaces that match the filter parameters.
  The [spaces](#iotda_spaces) structure is documented below.

<a name="iotda_spaces"></a>
The `spaces` block supports:

* `id` - The space ID.

* `name` - The space name.

* `created_at` - The creation time of the space. The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.

* `is_default` - Is it the default space.
