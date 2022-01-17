resource "null_resource" "enable_firestore" {
  provisioner "local-exec" {
    command = "make firestore"
  }

  depends_on = [google_firebase_project_location.default]
}

resource "google_firestore_index" "orders_user_time" {
  collection = "orders"

  fields {
    field_path = "Uuid"
    order      = "ASCENDING"
  }

  fields {
    field_path = "UserUuid"
    order      = "ASCENDING"
  }

  fields {
    field_path = "ProposedTime"
    order      = "ASCENDING"
  }

  fields {
    field_path = "ExpiresAt"
    order      = "ASCENDING"
  }

  fields {
    field_path = "__name__"
    order      = "ASCENDING"
  }

  depends_on = [null_resource.enable_firestore]
}
