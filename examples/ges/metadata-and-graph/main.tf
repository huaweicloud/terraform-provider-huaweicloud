resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.obs_bucket_name
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ges_metadata" "test" {
  depends_on = [
    huaweicloud_obs_bucket.test,
  ]

  name          = var.metadata_name
  description   = var.metadata_description
  metadata_path = "${huaweicloud_obs_bucket.test.bucket}/schema.xml"

  ges_metadata {
    labels {
      name       = "user"
      properties = [
        {
          "dataType"    = "string"
          "name"        = "name"
          "cardinality" = "single"
        },
        {
          "dataType"    = "long"
          "name"        = "age"
          "cardinality" = "single"
        },
      ]
    }

    labels {
      name       = "follows"
      properties = [
        {
          "dataType"    = "long"
          "name"        = "since"
          "cardinality" = "single"
        },
      ]
    }
  }
}

resource "huaweicloud_ges_graph" "test" {
  depends_on = [
    huaweicloud_ges_metadata.test,
  ]

  name                  = var.graph_name
  graph_size_type_index = var.graph_size_type_index
  cpu_arch              = "x86_64"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  crypt_algorithm       = "generalCipher"
  enable_https          = false
  tags                  = var.graph_tags
}
