---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_devserver_images"
description: |-
  Use this data source to get the list of ModelArts DevServer images.
---

# huaweicloud_modelarts_devserver_images

Use this data source to get the list of ModelArts DevServer images.

## Example Usage

### Query all DevServer images without any filter

```hcl
data "huaweicloud_modelarts_devserver_images" "test" {}
```

### Query the DevServer images using server type filter

```hcl
data "huaweicloud_modelarts_devserver_images" "test" {
  server_type = "ECS"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DevServer images are located.  
  If omitted, the provider-level region will be used.

* `server_type` - (Optional, String) Specifies the type of the server to be queried.  
  The valid values are as follows:
  + **BMS**: Bare metal server.
  + **ECS**: Elastic cloud server.
  + **HPS**: Hypernode server.

* `resource_flavor_name` - (Optional, String) Specifies the name of the resource flavor to be queried.

  -> You can use data source `huaweicloud_modelarts_devserver_flavors` to get the resource flavor names.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `images` - The list of DevServer images that matched filter parameters.  
  The [images](#modelarts_devserver_images_attr) structure is documented below.

<a name="modelarts_devserver_images_attr"></a>
The `images` block supports:

* `id` - The ID of the image.

* `name` - The name of the image.

* `server_type` - The type of the server.
  + **BMS**: Bare metal server.
  + **ECS**: Elastic cloud server.
  + **HPS**: Hypernode server.

* `arch` - The architecture of the image.
  + **ARM**
  + **X86**

* `status` - The status of the image.
  + **ACTIVE**
  + **INACTIVE**
