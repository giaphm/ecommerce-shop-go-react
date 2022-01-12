import { DefaultApi as OrdersDefaultApi, ApiClient as OrdersApiClient, NewOrder, NewOrderItem } from './clients/orders/src';

let ordersClient;

let getOrder;

let getOrders;

let getUserOrders;

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
  
  getUserOrders = function(userUuid, callback) {
  
    ordersAPI.getUserOrders(userUuid, (error, data, response) => {
      if (!error) {
        callback(data)
        console.log("data", data)
        console.log("response", response)
        return
      }
      console.error(error)
    })
  }
  
  createOrder = function(userUuid, orderItems, totalPrice, callback) {

    const newOrderItems = [];

    orderItems.map(orderItem => {
      const newOrderItem = new NewOrderItem(orderItem.productUuid, orderItem.quantity);
      newOrderItems.push(newOrderItem);
    })
  
    const newOrder = new NewOrder(userUuid, newOrderItems, totalPrice)
  
    ordersAPI.createOrder(newOrder, (error, data, response) => {
      if (!error) {
        console.log("data", data)
        console.log("response", response)
        callback(response)
        return
      }
      console.error(error)
    })
  }
  
  cancelOrder = function(orderUuid, callback) {
  
    ordersAPI.cancelOrder(orderUuid, (error, data, response) => {
      if (!error) {
        console.log("data", data)
        console.log("response", response)
        callback(response)
        return
      }
      console.error(error)
    })
  }
}

export { ordersClient };

export { getOrder };

export { getOrders };

export { getUserOrders };

export { createOrder };

export { cancelOrder };
