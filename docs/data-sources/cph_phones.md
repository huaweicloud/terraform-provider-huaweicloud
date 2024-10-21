---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_phones"
description: |-
  Use this data source to get the list of CPH phones.
---

# huaweicloud_cph_phones

Use this data source to get the list of CPH phones.

## Example Usage

```hcl
data "huaweicloud_cph_phones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `phone_name` - (Optional, String) Specifies the cloud phone name and support fuzzy query.

* `server_id` - (Optional, String) Specifies the cloud phone server ID.

* `status` - (Optional, Int) Specifies the cloud phone status.
  + **1**: Creating
  + **2**: Running
  + **3**: Resetting
  + **4**: Restarting
  + **6**: Freeze
  + **7**: Shutting down
  + **8**: Shut down
  + **-5**: Reset failed
  + **-6**: Restart failed
  + **-7**: Mobile phone abnormality
  + **-8**: Creation failed
  + **-9**: Shutdown failed

* `type` - (Optional, String) Specifies the cloud phone type.
  + **0**: Ordinary cloud phone
  + **1**: Trial cloud phone

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `phones` - The cloud phone list.

  The [phones](#phones_struct) structure is documented below.

<a name="phones_struct"></a>
The `phones` block supports:

* `phone_id` - The cloud phone ID.

* `phone_model_name` - The cloud phone flavor name.

* `status` - The cloud phone status.

* `create_time` - The create time.

* `update_time` - The update time.

* `phone_name` - The cloud phone name.

* `server_id` - The cloud phone server ID.

* `image_id` - The cloud phone image ID.

* `vnc_enable` - Whether to enable the VNC service on the cloud phone.
  + **true**: enable
  + **false**: disable

* `type` - The cloud phone type.

* `metadata` - The order and product related information.

  The [metadata](#phones_metadata_struct) structure is documented below.

* `volume_mode` - Whether the physical disk of the mobile phone is independent.
  + **0**: Not independent
  + **1**: Independent

* `availability_zone` - The availability zone where the cloud mobile server is located.

* `image_version` - The image version.

* `imei` - The IMEI of the phone.

* `traffic_type` - The phone routing type.
  + **direct**: default route
  + **routing**: routing to the encoding container

<a name="phones_metadata_struct"></a>
The `metadata` block supports:

* `order_id` - The order ID.

* `product_id` - The product ID.
