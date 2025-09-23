---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_access_credential"
description: -|
  Manages an access credential resource within HuaweiCloud.
---

# huaweicloud_iotda_access_credential

Manages an access credential resource within HuaweiCloud.

-> 1.This resource is only a one-time action resource for doing API action. Deleting this resource will not clear
  the corresponding request record, but will only remove the resource information from the tfstate file.
  <br>2. The client will only retain one record. If the resource is redeployed, it will only reset the access
  credential, rendering the previous credential invalid.

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
resource "huaweicloud_iotda_access_credential" "test" {
  type             = "AMQP"
  force_disconnect = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the access credential resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the type of the access credential.
  The valid values are as follows:
  + **AMQP** (Default value).
  + **MQTT**

  Changing this parameter will create a new resource.

* `force_disconnect` - (Optional, Bool, ForceNew) Specifies whether to disconnect AMQP or MQTT connection when
  creating access credential.
  The valid values are as follows:
  + **false** (Default value).
  + **true**

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `access_key` - The access name.

* `access_code` - The access credential.
