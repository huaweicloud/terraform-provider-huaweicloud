---
subcategory: "Cloud Connect (CC)"
---

# huaweicloud_cc_bandwidth_package

Manages a bandwidth package resource of Cloud Connect within HuaweiCloud.  

## Example Usage

```hcl
resource "huaweicloud_cc_bandwidth_package" "test" {
  name           = "demo"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 5
  description    = "This is a demo"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The bandwidth package name.  
  The name can contain a maximum of 64 characters.

* `local_area_id` - (Required, String, ForceNew) The local area ID.  
  Valid values are **Chinese-Mainland**, **Asia-Pacific**, **Africa**, **Western-Latin-America**,
   **Eastern-Latin-America** and **Northern-Latin-America**.

  Changing this parameter will create a new resource.

* `remote_area_id` - (Required, String, ForceNew) The remote area ID.  
  Valid values are **Chinese-Mainland**, **Asia-Pacific**, **Africa**, **Western-Latin-America**,
   **Eastern-Latin-America** and **Northern-Latin-America**.

  Changing this parameter will create a new resource.

* `charge_mode` - (Required, String, ForceNew) Billing option of the bandwidth package.  
  Valid value is **bandwidth**.

  Changing this parameter will create a new resource.

* `billing_mode` - (Required, String) Billing mode of the bandwidth package.  
  The options are as follows:
    + **3**: pay-per-use for the Chinese Mainland website.
    + **4**: pay-per-use for the International website.
    + **5**: 95th percentile bandwidth billing for the Chinese Mainland website.
    + **6**: 95th percentile bandwidth billing for the International website.

  -> This argument can only be modified to **5** and **6**.

* `bandwidth` - (Required, Int) Bandwidth in the bandwidth package.  

* `project_id` - (Optional, String, ForceNew) Project ID.
  If omitted, the provider-level project ID will be used.
  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description about the bandwidth package.  
  The description can contain a maximum of 85 characters.

* `enterprise_project_id` - (Optional, String) ID of the enterprise project that the bandwidth package
  belongs to. Value 0 indicates the default enterprise project.

* `resource_id` - (Optional, String) ID of the resource that the bandwidth package is bound to.  

* `resource_type` - (Optional, String) Type of the resource that the bandwidth package is bound to.  
   Valid value is **cloud_connection**.

* `tags` - (Optional, Map) The key/value pairs to associate with the bandwidth package.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Bandwidth package status.  
  The valid value are as follows:
    + **ACTIVE**: Bandwidth packages are available.

## Import

The bandwidth package can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cc_bandwidth_package.test 0ce123456a00f2591fabc00385ff1234
```
