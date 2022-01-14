// display profile
import React, { SyntheticEvent, useState, useEffect, Dispatch } from "react";
import Router from "next/router"
import {
  TextField,
  Button,
  Box,
  Typography,
  Divider,
  Snackbar,
  IconButton,
  Backdrop,
  CircularProgress,
} from "@mui/material";

import MuiAlert, { AlertProps } from '@mui/material/Alert';

import CloseIcon from '@mui/icons-material/Close';

import * as UsersAPI from "../src/repositories/users";
import { Auth, setApiClientsAuth } from "../src/repositories/auth";

import Layout from "../components/layout";
import CurrentUserAppCtx from "../store/current-user-context";

const Alert = React.forwardRef<HTMLDivElement, AlertProps>(function Alert(
  props,
  ref,
) {
  return <MuiAlert elevation={6} ref={ref} variant="standard" {...props} />;
});

function Profile() {
  const [email, setEmail] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [balance, setBalance] = useState("");
  const [role, setRole] = useState("");
  const [passwordUpdateInformation, setPasswordUpdateInformation] = useState("");
  const [passwordUpdatePassword, setPasswordUpdatePassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [newPasswordConfirm, setNewPasswordConfirm] = useState("");
  
  const [openBackdrop, setOpenBackdrop] = React.useState(false);

  const [
    showInvalidPasswordUpdateInformationNotification,
    setShowInvalidPasswordUpdateInformationNotification
  ] = useState(false);
  
  const [
    showInvalidPasswordUpdatePasswordNotification,
    setShowInvalidPasswordUpdatePasswordNotification
  ] = useState(false);
  
  const [
    showUnmatchedNewPasswordAndPasswordConfirmNotification,
    setShowUnmatchedNewPasswordAndPasswordConfirmNotification
  ] = useState(false);
  
  const currentUserAppCtx = React.useContext(CurrentUserAppCtx);
  
  console.log("currentUserAppCtx", currentUserAppCtx)

  
  console.log("email", email)
  console.log("displayName", displayName)
  console.log("balance", balance)
  console.log("role", role)
  console.log("passwordUpdateInformation", passwordUpdateInformation)
  console.log("passwordUpdatePassword", passwordUpdatePassword)
  console.log("newPassword", newPassword)
  console.log("newPasswordConfirm", newPasswordConfirm)

  React.useEffect(() => {
    console.log(typeof window);
    console.log("currentUserAppCtx", currentUserAppCtx);
    // const mockUserLoggedIn = JSON.parse(localStorage.getItem("_mock_user") || "{}")
    const isCurrentUserLoggedIn = Auth.isLoggedIn();
    console.log("isCurrentUserLoggedIn", isCurrentUserLoggedIn);
    if (isCurrentUserLoggedIn) {
      const currentUser = Auth.currentUser();
      console.log("currentUser", currentUser);
      // toast.message("Hey buddy!")
      Auth.waitForAuthReady()
        .then(() => {
          return Auth.getJwtToken(false)
        })
        .then((token: string) => setApiClientsAuth(token))
        .then(() => {        
          console.log("UsersAPI.usersClient", UsersAPI.usersClient);
          console.log("currentUser", currentUser);
          currentUserAppCtx!.fetchCurrentUser({
            uuid: currentUser["uuid"],
            email: currentUser["email"],
            displayName: currentUser["displayName"],
            role: currentUser["role"],
            balance: currentUser["balance"],
          });
          console.log("currentUserAppCtx", currentUserAppCtx);
    
          UsersAPI.getCurrentUser((user: any) => {
            console.log("UsersAPI.usersClient", UsersAPI.usersClient);
            console.log(
              "UsersAPI.usersClient.authentications",
              UsersAPI.usersClient.authentications
            );
            console.log("user", user);
            setEmail(user.email)
            setDisplayName(user.displayName)
            setBalance(user.balance)
            setRole(user.role)
            
        
            // console.log("setIsLoading(true);");
            // setIsLoading(true);
          });
        })
    } else if (!currentUserAppCtx!["uuid"]) {
      Router.push("/login");
    }
  }, []);


  const updateInfoSubmit = (e: SyntheticEvent) => {
    e.preventDefault();
    console.log("updateInfoSubmit")
    console.log("e", e)
    handleOpenBackdrop()

    UsersAPI.loginUser(currentUserAppCtx!.email, passwordUpdateInformation, (response: any) => {
      console.log(response);
      if(response.statusCode === 200){
        UsersAPI.updateUserInformation(currentUserAppCtx!.uuid, displayName, email, (response: any) => {
          console.log("response", response)
          // success
          if(response.statusCode === 204){
            console.log("update user information successfully")
            UsersAPI.loginUser(email, passwordUpdateInformation, (response: any) => {
              console.log(response);
              if(response.statusCode === 200){
                const updatedUser = {...response.body}
                console.log("updatedUser", updatedUser);
                currentUserAppCtx!.fetchCurrentUser({
                  uuid: updatedUser["uuid"],
                  email: updatedUser["email"],
                  displayName: updatedUser["displayName"],
                  role: updatedUser["role"],
                  balance: updatedUser["balance"],
                });
                setEmail(updatedUser.email)
                setDisplayName(updatedUser.displayName)
                setBalance(updatedUser.balance)
                setRole(updatedUser.role)
                handleCloseBackdrop()
              }
            })
          }
        })
      } else {
        setShowInvalidPasswordUpdateInformationNotification(true);
        handleCloseBackdrop()
      }
    })
  };

  const updatePasswordSubmit = (e: SyntheticEvent) => {
    e.preventDefault();
    console.log("updatePasswordSubmit")
    console.log("e", e)
    handleOpenBackdrop()

    if(newPassword !== newPasswordConfirm) {
      setTimeout(() => {
        console.log("New password confirm is invalid")
        setShowUnmatchedNewPasswordAndPasswordConfirmNotification(true);
        handleCloseBackdrop()
      }, 1000);
      return
    }

    UsersAPI.loginUser(email, passwordUpdatePassword, (response: any) => {
      console.log(response);
      if(response.statusCode === 200){
        UsersAPI.updateUserPassword(currentUserAppCtx!.uuid, newPassword, (response: any) => {
          console.log("response", response)
          // success
          if(response.statusCode === 204){
            console.log("update user password successfully")
            UsersAPI.loginUser(email, newPassword, (response: any) => {
              console.log(response);
              if(response.statusCode === 200){
                const updatedUser = {...response.body}
                console.log("updatedUser", updatedUser);
                currentUserAppCtx!.fetchCurrentUser({
                  uuid: updatedUser["uuid"],
                  email: updatedUser["email"],
                  displayName: updatedUser["displayName"],
                  role: updatedUser["role"],
                  balance: updatedUser["balance"],
                });
                setEmail(updatedUser.email)
                setDisplayName(updatedUser.displayName)
                setBalance(updatedUser.balance)
                setRole(updatedUser.role)
                handleCloseBackdrop()
              }
            })
          }
        })
      } else {
        setShowInvalidPasswordUpdatePasswordNotification(true);
        handleCloseBackdrop()
      }
    })
  };

  const handleOpenBackdrop = () => {
    setOpenBackdrop(true);
  }
  
  const handleCloseBackdrop = () => {
    setOpenBackdrop(false);
  }
  
  const handleCloseInvalidPasswordUpdateInformationNotification = () => {
    setShowInvalidPasswordUpdateInformationNotification(false);
  }

  const showInvalidPasswordUpdateInformationNotificationAction = (
      <IconButton
        size="small"
        aria-label="close"
        color="inherit"
        onClick={handleCloseInvalidPasswordUpdateInformationNotification}
      >
        <CloseIcon fontSize="small" />
      </IconButton>
  )
  

  const handleCloseInvalidPasswordUpdatePasswordNotification = () => {
    setShowInvalidPasswordUpdatePasswordNotification(false);
  }
  
  const showInvalidPasswordUpdatePasswordNotificationAction = (
      <IconButton
        size="small"
        aria-label="close"
        color="inherit"
        onClick={handleCloseInvalidPasswordUpdatePasswordNotification}
      >
        <CloseIcon fontSize="small" />
      </IconButton>
  )
  


  const handleCloseUnmatchedNewPasswordAndPasswordConfirmNotification = () => {
    setShowUnmatchedNewPasswordAndPasswordConfirmNotification(false);
  }

  const showUnmatchedNewPasswordAndPasswordConfirmNotificationAction = (
      <IconButton
        size="small"
        aria-label="close"
        color="inherit"
        onClick={handleCloseUnmatchedNewPasswordAndPasswordConfirmNotification}
      >
        <CloseIcon fontSize="small" />
      </IconButton>
  )

  return (
    <Layout role={currentUserAppCtx!.role}>
      <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
        open={showInvalidPasswordUpdateInformationNotification}
        onClose={handleCloseInvalidPasswordUpdateInformationNotification}
        autoHideDuration={5000}
        // message={showInvalidPasswordUpdateInformationNotification ?
        //   "Invalid password to update information" : ""
        // }
        // key={"top" + "center"} 
        action={showInvalidPasswordUpdateInformationNotificationAction}
      >
        <Alert onClose={handleCloseInvalidPasswordUpdateInformationNotification} severity="error" sx={{ width: '100%' }}>
          Invalid password to update information!
        </Alert>
      </Snackbar>
      <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
        open={showInvalidPasswordUpdatePasswordNotification}
        onClose={handleCloseInvalidPasswordUpdatePasswordNotification}
        autoHideDuration={5000}
        // message={showInvalidPasswordUpdatePasswordNotification ?
        //   "Invalid password to update password" : ""
        // }
        // key={"top" + "center"} 
        action={showInvalidPasswordUpdatePasswordNotificationAction}
      >
        <Alert onClose={handleCloseInvalidPasswordUpdatePasswordNotification} severity="error" sx={{ width: '100%' }}>
          Invalid password to update password!
        </Alert>
      </Snackbar>
      <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
        open={showUnmatchedNewPasswordAndPasswordConfirmNotification}
        onClose={handleCloseUnmatchedNewPasswordAndPasswordConfirmNotification}
        autoHideDuration={5000}
        // message={showUnmatchedNewPasswordAndPasswordConfirmNotification ?
        //   "Your new password and new password confirm is unmatched" : ""
        // }
        // key={"top" + "center"} 
        action={showUnmatchedNewPasswordAndPasswordConfirmNotificationAction}
      >
        <Alert onClose={handleCloseUnmatchedNewPasswordAndPasswordConfirmNotification} severity="error" sx={{ width: '100%' }}>
          Your new password and new password confirm is unmatched!
        </Alert>
      </Snackbar>
      {/* <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
        open={showUnmatchedNewPasswordAndPasswordConfirmNotification}
        onClose={handleCloseUnmatchedNewPasswordAndPasswordConfirmNotification}
        autoHideDuration={5000}
        message={showUnmatchedNewPasswordAndPasswordConfirmNotification ?
          "Your new password and new password confirm is unmatched" : ""
        }
        // key={"top" + "center"} 
        action={showUnmatchedNewPasswordAndPasswordConfirmNotificationAction}
      /> */}
      
      <Box
        sx={{
          pt: 5,
          px: 12,
        }}
      >
        <Typography sx={{textAlign: "center", fontSize: "2rem", my: 1,}} variant="h5">Account Information</Typography>
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: "center",
            justifyContent: "center",
          }}
          component={"form"}
          onSubmit={updateInfoSubmit}
        >
            <TextField
              sx={{mb: 2}}
              label="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            <TextField
              sx={{mb: 2}}
              label="Display Name"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
            />
            <TextField
              sx={{mb: 2}}
              label="Balance"
              value={`$ ${balance}`}
              variant={"filled"}
              InputProps={{
                readOnly: true,
              }}
            />
            <TextField
              sx={{mb: 2}}
              label="Role"
              value={role}
              variant={"filled"}
              InputProps={{
                readOnly: true,
              }}
            />
            <TextField
              sx={{mb: 2}}
              label="Password to update information"
              value={passwordUpdateInformation}
              type="password"
              onChange={(e) => setPasswordUpdateInformation(e.target.value)}
            />
          <Button sx={{mb: 2}} color="secondary" variant="contained" type="submit">
            Update
          </Button>
        </Box>

        <Divider variant="middle"/>

        <Typography sx={{textAlign: "center", fontSize: "2rem", my: 2,}} variant="h5">Update password</Typography>
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: "center",
            justifyContent: "center",
          }}
          component={"form"}
          onSubmit={updatePasswordSubmit}
          >
            <TextField
              sx={{mb: 2}}
              label="New Password"
              type="password"
              onChange={(e) => setNewPassword(e.target.value)}
            />
            <TextField
              sx={{mb: 2}}
              label="New Password Confirm"
              type="password"
              onChange={(e) => setNewPasswordConfirm(e.target.value)}
            />
            <TextField
              sx={{mb: 2}}
              label="Password to update password"
              value={passwordUpdatePassword}
              type="password"
              onChange={(e) => setPasswordUpdatePassword(e.target.value)}
            />
          <Button color="secondary" variant="contained" type="submit">
            Update
          </Button>
        </Box>
      </Box>
        <Backdrop
          sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
          open={openBackdrop}
        >
          <CircularProgress color="inherit" />
        </Backdrop>
    </Layout>
  );
}

export default Profile;
