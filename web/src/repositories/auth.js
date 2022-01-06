import firebase from "firebase/auth";
import {sign} from "jsonwebtoken";

import {getTestUsers, addTestUser, usersClient} from "./users";
import {checkoutsClient} from "./checkouts";
import {ordersClient} from "./orders";
import {productsClient} from "./products";

class FirebaseAuth {
    login(email, password) {
        return firebase.signInWithEmailAndPassword(firebase.getAuth(), email, password)
    }

    waitForAuthReady() {
        return new Promise((resolve) => {
            firebase
                .onAuthStateChanged(firebase.getAuth(), function () {
                    resolve()
                });
        })
    }

    getJwtToken(required) {
        return new Promise((resolve, reject) => {
            if (!firebase.getAuth().currentUser) {
                if (required) {
                    reject('no user found')
                } else {
                    resolve(null)
                }
                return
            }

            firebase.getAuth().currentUser.getIdToken(false)
                .then(function (idToken) {
                    resolve(idToken)
                })
                .catch(function (error) {
                    reject(error)
                });
        })
    }

    currentUser() {
        if (!firebase.getAuth().currentUser) {
            return null
        }

        return firebase.getAuth().currentUser
    }

    logout() {
        return new Promise(resolve => {
            if (!firebase.getAuth().currentUser) {
                resolve()
                return
            }

            return firebase.getAuth().signOut()
        })
    }


    isLoggedIn() {
        return firebase.getAuth().currentUser != null
    }
}

class MockAuth {
    login(email, password) {
        return new Promise((resolve, reject) => {
            setTimeout(function () {
                let found = getTestUsers().filter(u => u.email === email && u.password === password);

                if (found) {
                    localStorage.setItem('_mock_user', JSON.stringify(found[0]));
                    resolve()
                } else {
                    reject('invalid email or password')
                }
            }, 500) // simulate http request
        })
    }
    
    signup(displayName, email, password, role) {
        return new Promise((resolve, reject) => {
            setTimeout(function () {
                const numUsers = getTestUsers().length

                const newTestUser = {
                    'uuid': `${numUsers}`,
                    'email': email,
                    'password': password,
                    'role': role,
                    'name': displayName,
                }

                addTestUser(newTestUser);

                resolve()
                
            }, 500) // simulate http request
        })
    }


    waitForAuthReady() {
        return new Promise((resolve) => {
            setTimeout(resolve, 50)
        })
    }

    getJwtToken() {
        return new Promise((resolve) => {
            let user = this.currentUser()

            let claims = {
                'user_uuid': user.uuid,
                'email': user.email,
                'role': user.role,
                'name': user.name,
            }
            let token = sign(claims, 'mock_secret')
            resolve(token)
        })
    }

    currentUser() {
        let userStr = localStorage.getItem('_mock_user');
        if (!userStr) {
            return null
        }
    
        return JSON.parse(userStr)
    }

    logout() {
        return new Promise(resolve => {
            localStorage.setItem('_mock_user', '')

            setTimeout(resolve, 50)
        })
    }

    isLoggedIn() {
        return this.currentUser() !== null
    }
}

export function setApiClientsAuth(idToken) {
    console.log("idToken", idToken)
    checkoutsClient.authentications['bearerAuth'].accessToken = idToken
    ordersClient.authentications['bearerAuth'].accessToken = idToken
    productsClient.authentications['bearerAuth'].accessToken = idToken
    usersClient.authentications['bearerAuth'].accessToken = idToken
    console.log("productsClient", productsClient)
    console.log("usersClient", usersClient)
}

const MOCK_AUTH = process.env.NODE_ENV === 'development'
export let Auth

if (MOCK_AUTH) {
    Auth = new MockAuth()
} else {
    Auth = new FirebaseAuth()
}
