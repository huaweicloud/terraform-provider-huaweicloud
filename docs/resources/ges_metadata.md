---
subcategory: "Graph Engine Service (GES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ges_metadata"
description: ""
---

# huaweicloud_ges_metadata

Manages a GES metadata resource within HuaweiCloud.  

-> Only 50 metadata resources can be created.

## Example Usage

```hcl
variable "obs_path" {}

resource "huaweicloud_ges_metadata" "test" {
  name          = "demo"
  description   = "This is a demo"
  metadata_path = var.obs_path
  ges_metadata {
    labels {
      name = "user"
      properties = [{
        "dataType"    = "char"
        "name"        = "sex"
        "cardinality" = "single"
        },
        {
          "dataType"      = "enum"
          "name"          = "country"
          "cardinality"   = "single"
          "typeNameCount" = "3"
          "typeName1"     = "US"
          "typeName2"     = "EN"
          "typeName3"     = "CN"
        },
        {
          "dataType"    = "char array"
          "name"        = "firstName"
          "cardinality" = "single"
          "maxDataSize" = "20"
        },
        {
          "dataType"    = "string"
          "name"        = "lastName"
          "cardinality" = "single"
        },
        {
          "dataType"    = "long"
          "name"        = "children"
          "cardinality" = "set"
        },
        {
          "dataType"    = "long"
          "name"        = "friends"
          "cardinality" = "list"
        },
        {
          "dataType"      = "enum"
          "name"          = "cards"
          "cardinality"   = "list"
          "typeNameCount" = "3"
          "typeName1"     = "card_1"
          "typeName2"     = "card_2"
          "typeName3"     = "card_3"
        }
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Metadata name.  
  The name contains 1 to 64 characters consisting of only letters, digits, and underscores (_).
  Changing this parameter will create a new resource.

* `metadata_path` - (Required, String) OBS Path for storing the metadata.  

* `description` - (Required, String, ForceNew) Metadata description.  

  Changing this parameter will create a new resource.

* `ges_metadata` - (Required, List) Object for storing metadata message information.
  The [Metadata](#GesMetadata_Metadata) structure is documented below.

* `encryption` - (Optional, List) The configuration of data encryption.
  The graph instance is not encrypted by default.
  The [Encryption](#GesMetadata_Encryption) structure is documented below.

<a name="GesMetadata_Metadata"></a>
The `Metadata` block supports:

* `labels` - (Optional, List) Label list.  
  For details, see [data formats](https://support.huaweicloud.com/intl/en-us/usermanual-ges/ges_01_0153.html).
  The [Labels](#GesMetadata_MetadataLabels) structure is documented below.

<a name="GesMetadata_MetadataLabels"></a>
The `MetadataLabels` block supports:

* `name` - (Optional, String) Name of a label.

* `properties` - (Optional, List) The list of label properties. A property refers to the data format of a single
  property and contains some fields.
  For details, see [data formats](https://support.huaweicloud.com/intl/en-us/usermanual-ges/ges_01_0153.html).

<a name="GesMetadata_Encryption"></a>
The `Encryption` block supports:

* `enable` - (Optional, Bool) Whether to enable data encryption The value can be true or false.
  The default value is false.  

* `master_key_id` - (Optional, String) ID of the customer master key created by DEW in the project where
  the graph is created.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Status of a metadata. **200** is available.

## Import

The ges metadata can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ges_metadata.test 55b32ad9-1aba-407d-86cf-85f4f765d37a
```
