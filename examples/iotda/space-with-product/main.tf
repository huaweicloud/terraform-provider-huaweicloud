resource "huaweicloud_iotda_space" "test" {
  name = var.space_name
}

resource "huaweicloud_iotda_product" "test" {
  name              = var.product_name
  space_id          = huaweicloud_iotda_space.test.id
  device_type       = "TemperatureSensor"
  protocol          = "MQTT"
  data_type         = "json"
  manufacturer_name = "demo_manufacturer"
  industry          = "smart_home"
  description       = "Created-by-Terraform-for-IoTDA-best-practice-example"

  services {
    id   = "temperature_service"
    type = "sensor"

    properties {
      name        = "temperature"
      type        = "decimal"
      method      = "RW"
      min         = "-40"
      max         = "85"
      step        = 0.1
      unit        = "C"
      description = "The-current-temperature"
    }

    commands {
      name = "reset"

      paras {
        name = "delay"
        type = "int"
        min  = "0"
        max  = "60"
      }

      responses {
        name = "result"
        type = "string"
      }
    }
  }
}
