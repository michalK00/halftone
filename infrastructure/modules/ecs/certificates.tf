resource "tls_private_key" "key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "cert" {
  private_key_pem = tls_private_key.key.private_key_pem

  subject {
    common_name  = aws_lb.main.dns_name
    organization = "Example Organization"
  }

  validity_period_hours = 8760 # Valid for 1 year

  dns_names = [
    aws_lb.main.dns_name
  ]

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]

  depends_on = [aws_lb.main]
}

resource "aws_acm_certificate" "self_signed" {
  private_key      = tls_private_key.key.private_key_pem
  certificate_body = tls_self_signed_cert.cert.cert_pem

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "self-signed-certificate"
  }
}
