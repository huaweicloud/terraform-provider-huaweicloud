---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_custom_authentications"
description: |-
  Use this data source to get the list of IoTDA custom authentications.
---

# huaweicloud_iotda_custom_authentications

Use this data source to get the list of IoTDA custom authentications.

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
data "huaweicloud_iotda_custom_authentications" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `authorizer_name` - (Optional, String) Specifies the name of the custom authentication.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `authorizers` - The list of the custom authentications.
  The [authorizers](#iotda_authorizers) structure is documented below.

<a name="iotda_authorizers"></a>
The `authorizers` block supports:

* `authorizer_id` - The ID of the custom authentication.

* `authorizer_name` - The name of the custom authentication.

* `func_name` - The name of the function associated with the custom authentication.

* `func_urn` -  The URN of the function associated with the custom authentication.

* `signing_enable` - Whether to enable signature authentication.

* `default_authorizer` - Whether the custom authentication is the default authentication mode.

* `status` - Whether to enable the custom authentication mode.

* `cache_enable` - Whether to enable the cache function.

* `create_time` - The creation time of the custom authentication.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

* `update_time` - The latest update time of the custom authentication.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.
