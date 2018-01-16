
resource "huaweicloud_elb_loadbalancer" "loadbalancer_1" {
  name               = "test-elb-name"
  description        = "test elb description"
  type               = "Internal"
  security_group_id  = "${var.security_group_id}"
  vpc_id             = "${var.vpc_id}"
  vip_subnet_id      = "${var.vip_subnet_id}"
  vip_address        = "${var.vip_address}"
  tenantid           = "${var.tenantid}"
  az                 = "${var.available_zone}"
  admin_state_up     = 1
}


resource "huaweicloud_elb_listener" "listener_1" {
  name              = "listener_1"
  protocol          = "TCP"
  port              = 22
  backend_protocol  = "TCP"
  backend_port      = 80
  lb_algorithm      = "roundrobin"
  loadbalancer_id   = "${huaweicloud_elb_loadbalancer.loadbalancer_1.id}"
  depends_on        = ["huaweicloud_elb_loadbalancer.loadbalancer_1"]
}

resource "huaweicloud_elb_healthcheck" "cz_health_1" {
  listener_id            = "${huaweicloud_elb_listener.listener_1.id}"
  healthcheck_protocol   = "HTTP"
  healthy_threshold      = 5
  healthcheck_timeout    = 25
  healthcheck_interval   = 3
  healthcheck_uri        = "/"
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
  depends_on =["huaweicloud_elb_listener.listener_1"]
}

resource "huaweicloud_elb_backendecs" "cz_backend_3" {
  private_address = "${var.vm_private_address}"
  listener_id = "${huaweicloud_elb_listener.listener_1.id}"
  server_id = "${var.vm_id}"
  depends_on =["huaweicloud_elb_listener.listener_1"]
}

