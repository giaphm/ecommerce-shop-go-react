import { DefaultApi as OrdersDefaultApi, ApiClient as OrdersApiClient, CreateOrder } from './clients/orders/src';

let ordersClient;

let getOrder;

let getOrders;

let createOrder;

let cancelOrder;

// in client-side rendering
if (typeof window == "object") {
  const serverSettings = {
    hostname: window.location.hostname,
  };
  ordersClient = new OrdersApiClient()
  ordersClient.basePath = ordersClient.getBasePathFromSettings(0, serverSettings);
  let ordersAPI = new OrdersDefaultApi(ordersClient)
  
  if (process.env.NODE_ENV === 'development') {
    ordersClient.basePath = "http://localhost:3002/api"
  }
  
  getOrder = function(orderUuid, callback) {
  
    ordersAPI.getOrder(orderUuid, (error, data, response) => {
      if (!error) {
        callback(data)
        console.log("data", data)
        console.log("response", response)
        return
      }
      console.error(error)
    })
  }
  
  getOrders = function(callback) {
  
    ordersAPI.getOrders((error, data, response) => {
      if (!error) {
        callback(data)
        console.log("data", data)
        console.log("response", response)
        return
      }
      console.error(error)
    })
  }
  
  createOrder = function(productUuids, totalPrice) {
  
    const createOrder = new CreateOrder(productUuids, totalPrice)
  
    ordersAPI.createOrder(createOrder, (error, data, response) => {
      if (!error) {
        console.log("data", data)
        console.log("response", response)
        return
      }
      console.error(error)
    })
  }
  
  cancelOrder = function(orderUuid) {
  
    ordersAPI.cancelOrder(orderUuid, (error, data, response) => {
      if (!error) {
        console.log("data", data)
        console.log("response", response)
        return
      }
      console.error(error)
    })
  }
}

export { ordersClient };

export { getOrder };

export { getOrders };

export { createOrder };

export { cancelOrder };
