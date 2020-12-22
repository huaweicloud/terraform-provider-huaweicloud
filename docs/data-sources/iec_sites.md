# huaweicloud\_iec\_sites

Use this data source to get the available of HuaweiCloud IEC sites.

## Example Usage

### Basic IEC Sites

```hcl
data "huaweicloud_iec_sites" "iec_sites" {}
```

### IEC Sites with sites block
```hcl
data "huaweicloud_iec_sites" "iec_sites" {
  sites {
    name = "iec-site-1"
  }
  sites {
    name = "iec-site-2"
  }
}
```

## Argument Reference

The following arguments are supported:
 
* `area` - (Optional, String) Specifies the area of the iec instance located.

* `province` - (Optional, String) Specifies the province of the iec instance 
    located.

* `city` - (Optional, String) Specifies the city of the iec instance located. 

* `operator` - (Optional, String) Specifies the operator supported of the iec 
    instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the iec service site.

* `status` - The current status of the iec service site.

* `sites` - An array of one or more iec service sites.
    The images object structure is documented below.

The `sites` block supports:

* `id` - The id of the iec service site.
* `name` - The name of the iec service site.
* `area` - The area of the iec service site located.
* `province` - The province of the iec service site located.
* `city` - The city of the iec service site located.
* `operator` - The operator information of the iec service site.
* `status` - The current status of the iec service site.
