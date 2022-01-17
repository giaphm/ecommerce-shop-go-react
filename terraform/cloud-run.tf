module cloud_run_checkouts_http {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "checkouts"
  protocol = "http"
  auth     = false

  envs = [
    {
      name  = "ORDERS_GRPC_ADDR"
      value = module.cloud_run_orders_grpc.endpoint
    },
    {
      name  = "PRODUCTS_GRPC_ADDR"
      value = module.cloud_run_products_grpc.endpoint
    },
    {
      name  = "USERS_GRPC_ADDR"
      value = module.cloud_run_users_grpc.endpoint
    }
  ]
}
module cloud_run_orders_grpc {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "orders"
  protocol = "grpc"
}

module cloud_run_orders_http {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "orders"
  protocol = "http"
  auth     = false

  envs = [
    {
      name  = "ORDERS_GRPC_ADDR"
      value = module.cloud_run_orders_grpc.endpoint
    }
  ]
}

module cloud_run_products_grpc {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "products"
  protocol = "grpc"
}

module cloud_run_products_http {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "products"
  protocol = "http"
  auth     = false

  envs = [
    {
      name  = "PRODUCTS_GRPC_ADDR"
      value = module.cloud_run_products_grpc.endpoint
    }
  ]
}

module cloud_run_users_grpc {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "users"
  protocol = "grpc"
}

module cloud_run_users_http {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images

  name     = "users"
  protocol = "http"
  auth     = false

  envs = [
    {
      name  = "USERS_GRPC_ADDR"
      value = module.cloud_run_users_grpc.endpoint
    }
  ]
}
