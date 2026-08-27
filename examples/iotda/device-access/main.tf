# Create an IoTDA resource space
resource "huaweicloud_iotda_space" "test" {
  name = var.space_name
}

# Create an IoTDA product under the space
resource "huaweicloud_iotda_product" "test" {
  name        = var.product_name
  device_type = var.product_device_type
  protocol    = var.product_protocol
  space_id    = huaweicloud_iotda_space.test.id
  data_type   = var.product_data_type

  services {
    id   = var.product_service_id
    type = var.product_service_type
  }
}

# Create an IoTDA device under the product
resource "huaweicloud_iotda_device" "test" {
  node_id    = var.device_node_id
  name       = var.device_name
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
  secret     = var.device_secret
}
