---
subcategory: "Relational Database Service (RDS)"
---

# huaweicloud_rds_engine_versions

Use this data source to obtain all version information of the specified engine type of HuaweiCloud RDS.

## Example Usage

```hcl
data "huaweicloud_rds_engine_versions" "test" {
  type = "SQLServer"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the RDS engine versions.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the RDS engine type.
  The valid values are **MySQL**, **PostgreSQL** and **SQLServer**, default to **MySQL**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID in hashcode format.

* `versions` - List of RDS versions. Structure is documented below.

The `versions` block contains:

* `id` - Version ID.

* `name` - Version name.
