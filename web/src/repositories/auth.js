import { 
    getAuth,
    signInWithEmailAndPassword,
    onAuthStateChanged,
} from "firebase/auth";
import {sign} from "jsonwebtoken";

import { usersClient } from "./users";
import {checkoutsClient} from "./checkouts";
import {ordersClient} from "./orders";
import {productsClient} from "./products";

class FirebaseAuth {
    login(user) {
        return signInWithEmailAndPassword(getAuth(), user.email, user.password)
    }

    waitForAuthReady() {
        return new Promise((resolve) => {
            onAuthStateChanged(getAuth(), function () {
                    resolve()
                });
        })
    }

    getJwtToken(required) {
        return new Promise((resolve, reject) => {
            if (!getAuth().currentUser) {
                if (required) {
                    reject('no user found')
                } else {
                    resolve(null)
                }
                return
            }

            getAuth().currentUser.getIdToken(false)
                .then(function (idToken) {
                    resolve(idToken)
                })
                .catch(function (error) {
                    reject(error)
                });
        })
    }

    currentUser() {
        if (!getAuth().currentUser) {
            return null
        }

        return getAuth().currentUser
    }

    logout() {
        return new Promise(resolve => {
            if (!getAuth().currentUser) {
                resolve()
                return
            }

            return getAuth().signOut()
        })
    }


    isLoggedIn() {
        return getAuth().currentUser != null
    }
}

class MockAuth {
    
    // auth login (mock) is to set local storage
    login(user) {
        return new Promise((resolve, reject) => {
            setTimeout(function () {
                // let found = getTestUsers().filter(u => u.email === email && u.password === password);

                // console.log("found", found)

                console.log("user", user)

                if (user) {
                    localStorage.setItem('_mock_user', JSON.stringify(user));
                    resolve()
                } else {
                    reject('invalid email or password')
                }
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
                'name': user.displayName,
            }
            console.log("claims", claims)
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
