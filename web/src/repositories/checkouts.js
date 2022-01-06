import { DefaultApi as CheckoutsDefaultApi, ApiClient as CheckoutsApiClient, Checkout } from './clients/checkouts/src';

let checkoutsClient;

let createCheckout;

let getCheckouts;

// in client-side rendering
if (typeof window == "object") {
  const serverSettings = {
    hostname: window.location.hostname,
  };
  checkoutsClient = new CheckoutsApiClient()
  checkoutsClient.basePath = checkoutsClient.getBasePathFromSettings(0, serverSettings);
  let checkoutsAPI = new CheckoutsDefaultApi(checkoutsClient)
  
  if (process.env.NODE_ENV === 'development') {
    checkoutsClient.basePath = "http://localhost:3001/api"
  }
  
  createCheckout = function(orderUuid, notes, proposedTime){
  
    const newCheckout = new Checkout(orderUuid, notes, proposedTime)
  
    checkoutsAPI.createCheckout(newCheckout, (error, data, response) => {
      if (!error){
        console.log("Calling createCheckout to checkout service successfully!")
        console.log('data', data)
        console.log('response', response)
        return
      }
      console.error(err)
    })
  }
  
  getCheckouts = function(callback){
  
    checkoutsAPI.getCheckouts((error, data, response) => {
      if (!error){
        callback(data)
        console.log("Calling getCheckouts to checkout service successfully!")
        console.log('data', data)
        console.log('response', response)
        return
      }
      console.error(err)
    })
  }
}

export { checkoutsClient };

export { createCheckout }

export { getCheckouts };