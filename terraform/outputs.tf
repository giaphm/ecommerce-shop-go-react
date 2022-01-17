output "checkouts_http_url" {
  value = module.cloud_run_checkouts_http.url
}

output "orders_grpc_url" {
  value = module.cloud_run_orders_grpc.url
}

output "orders_http_url" {
  value = module.cloud_run_orders_http.url
}

output "products_grpc_url" {
  value = module.cloud_run_products_grpc.url
}

output "products_http_url" {
  value = module.cloud_run_products_http.url
}

output "users_grpc_url" {
  value = module.cloud_run_users_grpc.url
}

output "users_http_url" {
  value = module.cloud_run_users_http.url
}

output "repo_url" {
  value = google_sourcerepo_repository.ecommerce_shop_go_react.url
}
