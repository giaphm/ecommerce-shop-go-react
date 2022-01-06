import { useState } from 'react'
import type { AppProps } from 'next/app'
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from "@mui/material/CssBaseline";  
import theme from '../src/theme';
import "../styles/globals.css"

import CurrentUserAppCtx, { CurrentUser } from '../store/current-user-context';


function MyApp({ Component, pageProps }: AppProps) {

  const [currentUser, setCurrentUser] = useState<CurrentUser>({
    uuid: null,
    email: null,
    displayName: null,
    role: null,
  })

  
  const fetchCurrentUser = (currentUser: CurrentUser) => {
      setCurrentUser({
        uuid: currentUser["uuid"],
        email: currentUser["email"],
        displayName: currentUser["displayName"],
        role: currentUser["role"],
      })
  }

  const removeCurrentUser = () => {
    setCurrentUser({
      uuid: null,
      email: null,
      displayName: null,
      role: null,
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
