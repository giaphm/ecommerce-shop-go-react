import {
    DefaultApi as UsersDefaultApi,
    ApiClient as UsersApiClient,
    UserSignIn,
    UserSignUp,
    UpdatedUserInformation,
    UpdatedUserPassword,
} from './clients/users/src'

import 'firebase/auth';

import {Auth, setApiClientsAuth} from "./auth";

export const Shopkeeper = 'shopkeeper';
export const User = 'user';

let usersClient;

let getUserRole;

let getUsers;

let getUserBalance;

let getCurrentUser;

let loginUser;

let signupUser;

let updateUserInformation;

let updateUserPassword;

let getTestUsers;

let addTestUser;

// in client-side rendering
if (typeof window == "object") {
    const serverSettings = {
      hostname: window.location.hostname,
    };
    usersClient = new UsersApiClient()
    usersClient.basePath = usersClient.getBasePathFromSettings(0, serverSettings);
    let usersAPI = new UsersDefaultApi(usersClient)
    
    if (process.env.NODE_ENV === 'development') {
        usersClient.basePath = "http://localhost:3004/api"
    }
    
    getUserRole = function() {
        return localStorage.getItem('role')
    }

    getUsers = function(callback) {
        return usersAPI.getUsers((error, data, response) => {
            console.log('error', error)
            console.log('data', data)
            console.log('response', response)
            if (!error) {
                callback(data)
                return
            }
            console.error(error)
        })
    }
    
    getUserBalance = function(callback) {
        return usersAPI.getCurrentUser((error, data, response) => {
            if (!error) {
                callback(data.balance)
                console.log('data', data)
                console.log('response', response)
                return
            }
            console.error(error)
        })
    }

    getCurrentUser = function(callback) {
        return usersAPI.getCurrentUser((error, data, response) => {
            if (!error) {
                callback(data)
                console.log('data', data)
                console.log('response', response)
                return
            }
            console.error(error)
        })
    }
    
    loginUser = function(email, password, callback) {

        console.log("email", email)
        console.log("password", password)
    
        const userSignIn = new UserSignIn(email, password)
    
        usersAPI.signIn(userSignIn, (error, data, response) => {
            console.log('error', error)
            console.log('data', data)
            console.log('response', response)
            if (!error) {
                console.log("Calling signin to users service successfully!")

                const user = {
                    email,
                    password,
                };
                Auth.login(user) // set to local storage for Mock Auth
                    .then(() => Auth.waitForAuthReady())
                    .then(() => {
                        return Auth.getJwtToken(false)
                    })
                    .then(token => setApiClientsAuth(token))
                    .then(() => callback(response))
                return
            }
            console.error(error)
            callback(response)
        })
    }
    
    signupUser = function(displayName, email, password, role, callback) {
    
        const userSignUp = new UserSignUp(displayName, email, password, role)
    
        usersAPI.signUp(userSignUp, (error, data, response) => {
            if (!error) {
                console.log("Calling to signup user successfully!")
                console.log('data', data)
                console.log('response', response)
                callback(response)
                return
            }
            console.error(error)
            callback(response)
        })
    }

    updateUserInformation = function(userUuid, newDisplayName, newEmail, callback) {
        
        const newUpdatedUserInformation = new UpdatedUserInformation(userUuid, newDisplayName, newEmail)

        usersAPI.updateUserInformation(newUpdatedUserInformation, (error, data, response) => {
            console.log('data', data)
            console.log('response', response)
            if (!error) {
                console.log("Calling to update user information successfully!")
                callback(response)
                return
            }
            console.error(error)
            callback(response)
        })
    }
    
    updateUserPassword = function(userUuid, newPassword, callback) {
        
        const newUpdatedUserPassword = new UpdatedUserPassword(userUuid, newPassword)

        usersAPI.updateUserPassword(newUpdatedUserPassword, (error, data, response) => {
            console.log('data', data)
            console.log('response', response)
            if (!error) {
                console.log("Calling to update user password successfully!")
                callback(response)
                return
            }
            console.error(error)
            callback(response)
        })
    }
    
}

export { usersClient };

export { getUserRole };

export { getUsers };

export { getUserBalance };

export { getCurrentUser };

export { loginUser };

export { signupUser };

export { updateUserInformation };

export { updateUserPassword };

export { getTestUsers };

export { addTestUser };