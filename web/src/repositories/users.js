import { DefaultApi as UsersDefaultApi, ApiClient as UsersApiClient, UserSignIn, UserSignUp } from './clients/users/src'

import 'firebase/auth';

import {Auth, setApiClientsAuth} from "./auth";

export const Shopkeeper = 'shopkeeper';
export const User = 'user';

let usersClient;

let getUserRole;

let getUsers;

let getUserBalance;

let loginUser;

let signupUser;

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
    
    loginUser = function(email, password) {
    
        const userSignIn = new UserSignIn(email, password)
    
        usersAPI.signIn(userSignIn, (error, data, response) => {
            if (!error) {
                console.log("Calling signin to users service successfully!")
                console.log('data', data)
                console.log('response', response)
                return
            }
            console.error(error)
        })
    
        return Auth.login(email, password)
            .then(function () {
                return Auth.waitForAuthReady()
            })
            .then(function () {
                return Auth.getJwtToken(false)
            })
            .then(token => {
                setApiClientsAuth(token)
            })
            .then(function () {
                return new Promise(((resolve, reject) => {
                    usersAPI.getCurrentUser((error, data) => {
                        if (!error) {
                            resolve(data)
                            return
                        }
                        reject(error)
                    })
                }))
            })
            .then(data => {
                console.log(data)
                localStorage.setItem('role', data.role)
                return new Promise(((resolve, reject) => {
                    if(data) {
                        resolve(data)
                    }
                    else {
                        reject(data)
                    }
                }))
            })
    }
    
    signupUser = function(displayName, email, password, role) {
    
        const userSignUp = new UserSignUp(displayName, email, password, role)
    
        usersAPI.signUp(userSignUp, (error, data, response) => {
            if (!error) {
                console.log("Calling to signup user successfully!")
                console.log('data', data)
                console.log('response', response)
                return
            }
            console.error(error)
        })
    
        return Auth.signup(displayName, email, password, role)
            .then(() => Auth.waitForAuthReady())
            .then(() => console.log("Sign up successfully!"))
    }
    
    const testUsers = [
        {
            'uuid': '0',
            'email': 'shopkeeper1@gmail.com',
            'password': '123456',
            'role': 'shopkeeper',
            'name': 'Raheem Arnold',
        },
        {
            'uuid': '1',
            'email': 'user1@gmail.com',
            'password': '123456',
            'role': 'user',
            'name': 'Mariusz Pudzianowski',
        },
        {
            'uuid': '2',
            'email': 'user2@gmail.com',
            'password': '123456',
            'role': 'user',
            'name': 'Arnold Schwarzenegger',
        },
    ]
    
    getTestUsers = function() {
        return testUsers
    }
    
    addTestUser = function(newTestUser) {
        testUsers.push(newTestUser)
    }    
}

export { usersClient };

export { getUserRole };

export { getUsers };

export { getUserBalance };

export { loginUser };

export { signupUser };

export { getTestUsers };

export { addTestUser };