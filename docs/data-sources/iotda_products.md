---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_products"
description: ""
---

# huaweicloud_iotda_products

Use this data source to get the list of IoTDA products within HuaweiCloud.

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
variable "product_id" {}

data "huaweicloud_iotda_products" "test" {
  product_id = var.product_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the products.
  If omitted, the provider-level region will be used.

* `product_id` - (Optional, String) Specifies the ID of the product to be queried.

* `product_name` - (Optional, String) Specifies the name of the product to be queried.

* `space_id` - (Optional, String) Specifies the space ID of the products to be queried.
  If omitted, query all products under the current instance.

* `space_name` - (Optional, String) Specifies the space name of the products to be queried.

* `device_type` - (Optional, String) Specifies the device type of the products to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `products` - All products that match the filter parameters.
  The [products](#iotda_products) structure is documented below.

<a name="iotda_products"></a>
The `products` block supports:

* `space_id` - The space ID to which the product belongs.

* `space_name` - The space name to which the product belongs.

* `id` - The product ID.

* `name` - The product name.

* `device_type` - The device type of the product.

* `protocol_type` - The protocol type used by devices under the product.
  The value can be **MQTT**, **CoAP**, **HTTP**, **HTTPS**, **Modbus**, **ONVIF**, **OPC-UA**, **OPC-DA**, or **Other**.

* `data_type` - The format of data reported by devices under the product.
  The value can be **json** or **binary**.

* `manufacturer_name` - The manufacturer name.

* `industry` - The industry to which the devices under the product belongs.

* `description` - The description of the product.

* `created_at` - The creation time of the product. The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.
