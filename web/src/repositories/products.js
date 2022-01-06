import { DefaultApi as ProductsDefaultApi, ApiClient as ProductsApiClient, NewProduct, Product, UpdatedProduct } from './clients/products/src';


let productsClient;

let getProduct;

let getProducts;

let getShopkeeperProducts;

let addProduct;

let updateProduct;

let deleteProduct;

// in client-side rendering
if (typeof window == "object") {
  const serverSettings = {
    hostname: window.location.hostname,
  };
  productsClient = new ProductsApiClient()
  productsClient.basePath = productsClient.getBasePathFromSettings(0, serverSettings);
  let productsAPI = new ProductsDefaultApi(productsClient)
  
  if (process.env.NODE_ENV === 'development') {
    productsClient.basePath = "http://localhost:3003/api"
  }
  
  getProduct = function(productUuid, callback) {
  
    productsAPI.getProduct(productUuid, (error, data, response) => {
      if (!error) {
        callback(data)
        console.log("data", data)
        console.log("response", response)
        return
      }
      console.error(error)
    })
  }
  
  getProducts = function(callback) {
  
    productsAPI.getProducts((error, data, response) => {
      console.log("error", error)
      console.log("data", data)
      console.log("response", response)
      if (error) {
        console.error(error)
      } else {
        callback(data)
      }
    })
  }
  
  getShopkeeperProducts = function(callback) {
  
    productsAPI.getShopkeeperProducts((error, data, response) => {
      if (!error) {
        callback(data)
        console.log("data", data)
        console.log("response", response)
        return
      }
      console.error(error)
    })
  }
  
  addProduct = function(category, title, image, description, price, quantity) {
  
    const newProduct = new NewProduct(category, title, image, description, price, quantity)
  
    productsAPI.addProduct(newProduct, (error, data, response) => {
      if (!error) {
        console.log("data", data)
        console.log("response", response)
        return
      }
      console.error(error)
    })
  }
  
  updateProduct = function(productUuid, userUuid, category, title, description, image, price, quantity){
  
    const newUpdateProduct = new UpdatedProduct(productUuid, userUuid, category, title, image, description, price, quantity)
  
    productsAPI.updateProduct(productUuid, newUpdateProduct, (error, data, response) => {
      if (!error) {
        console.log("data", data)
        console.log("response", response)
        return
      }
      console.error(error)
    })
  }
  
  deleteProduct = function(productUuid, callback) {
  
    productsAPI.deleteProduct(productUuid, (error, data, response) => {
      if (!error) {
        console.log("data", data)
        console.log("response", response)
        callback()
        return
      }
      console.error(error)
    })
  }
}

export { productsClient };

export { getProduct };

export { getProducts };

export { getShopkeeperProducts };

export { addProduct };

export { updateProduct };

export { deleteProduct };