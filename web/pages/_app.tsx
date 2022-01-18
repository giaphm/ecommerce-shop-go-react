import { useState, useEffect } from 'react'
import type { AppProps } from 'next/app'
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from "@mui/material/CssBaseline";

import theme, { darkTheme } from '../src/theme';
import "../styles/globals.css";

import { Auth, setApiClientsAuth } from "../src/repositories/auth";
import { loadFirebaseConfig } from "../src/firebase";

import CurrentUserAppCtx, { CurrentUser } from '../store/current-user-context';


function MyApp({ Component, pageProps }: AppProps) {

  const [currentUser, setCurrentUser] = useState<CurrentUser>({
    uuid: null,
    email: null,
    displayName: null,
    role: null,
    balance: 0,
  })
  
  // useEffect(() => {
  //   if(typeof window !== "undefined") {
  //     console.log(typeof window)
  //     loadFirebaseConfig()
  //   }
  // }, [])
  
  const fetchCurrentUser = (currentUser: CurrentUser) => {
      setCurrentUser({
        uuid: currentUser["uuid"],
        email: currentUser["email"],
        displayName: currentUser["displayName"],
        role: currentUser["role"],
        balance: currentUser["balance"],
      })
  }

  const removeCurrentUser = () => {
    setCurrentUser({
      uuid: null,
      email: null,
      displayName: null,
      role: null,
      balance: 0,
    })
  }

  const value = {...currentUser, fetchCurrentUser, removeCurrentUser}

  return (
    <CurrentUserAppCtx.Provider value={value}>
      <ThemeProvider theme={theme}>
        {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
        <CssBaseline />
        <Component {...pageProps} />
      </ThemeProvider>
    </CurrentUserAppCtx.Provider>
  )
}

export default MyApp
