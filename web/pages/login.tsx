import * as React from "react";
import { useRouter } from "next/router";

import {
  Backdrop,
  CircularProgress,
  IconButton,
} from "@mui/material";

import CloseIcon from '@mui/icons-material/Close';

import Avatar from "@mui/material/Avatar";
import Button from "@mui/material/Button";
import CssBaseline from "@mui/material/CssBaseline";
import TextField from "@mui/material/TextField";
import FormControlLabel from "@mui/material/FormControlLabel";
import Checkbox from "@mui/material/Checkbox";
import Link from "@mui/material/Link";
import Paper from "@mui/material/Paper";
import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import Typography from "@mui/material/Typography";
import Snackbar from '@mui/material/Snackbar';

import * as UsersAPI from "../src/repositories/users";
import { Auth } from "../src/repositories/auth";
import CurrentUserAppCtx from '../store/current-user-context'

import { isEmpty } from "lodash";

function Copyright(props: any) {
  return (
    <Typography
      variant="body2"
      color="text.secondary"
      align="center"
      {...props}
    >
      {"Copyright © "}
      <Link color="inherit" href="https://mui.com/">
        ecommerce-shop-go-react.app
      </Link>{" "}
      {new Date().getFullYear() - 1} - {new Date().getFullYear()}
      {"."}
    </Typography>
  );
}



export default function SignIn() {
  const [showInvalidLoginNotification, setShowInvalidLoginNotification] = React.useState(false);
  const [showLoggedOutNotification, setShowLoggedOutNotification] = React.useState(false);
  const [showSignedUpNotification, setShowSignedUpNotification] = React.useState(false);
  const [openBackdrop, setOpenBackdrop] = React.useState(false);
  

  const currentUserContext = React.useContext(CurrentUserAppCtx);
  console.log("currentUserContext", currentUserContext)

  console.log("showInvalidLoginNotification", showInvalidLoginNotification)
  
  console.log("showLoggedOutNotification", showLoggedOutNotification);
  
  console.log("showSignedUpNotification", showSignedUpNotification)

  const router = useRouter()
  const { query } = router;
  console.log("query", query)

  const { loggedOut, signedUp } = router.query;

  console.log("loggedOut", loggedOut)
  
  console.log("signedUp", signedUp)

  React.useEffect(() => {
    if (Auth.isLoggedIn()) {
      router.push("/");
    }
    setShowLoggedOutNotification(loggedOut ? true : false);
    setShowSignedUpNotification(signedUp ? true : false);
  }, []);

  const handleOpenBackdrop = () => {
    setOpenBackdrop(true);
  }
  
  const handleCloseBackdrop = () => {
    setOpenBackdrop(false);
  }
  
  
  const handleCloseSignedUpNotification = () => {
    setShowSignedUpNotification(false);
  }

  const showSignedUpNotificationAction = (
    <IconButton
      size="small"
      aria-label="close"
      color="inherit"
      onClick={handleCloseSignedUpNotification}
    >
      <CloseIcon fontSize="small" />
    </IconButton>
  )

  const handleCloseLoggedOutNotification = () => {
    setShowLoggedOutNotification(false);
  }

  const showLoggedOutNotificationAction = (
      <IconButton
        size="small"
        aria-label="close"
        color="inherit"
        onClick={handleCloseLoggedOutNotification}
      >
        <CloseIcon fontSize="small" />
      </IconButton>
  )

  const handleCloseInvalidLoginNotification = () => {
    setShowInvalidLoginNotification(false);
  }

  const showInvalidLoginNotificationAction = (
      <IconButton
        size="small"
        aria-label="close"
        color="inherit"
        onClick={handleCloseInvalidLoginNotification}
      >
        <CloseIcon fontSize="small" />
      </IconButton>
  )

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    console.log("handleOpenBackdrop();")
    handleOpenBackdrop();

    const data = new FormData(event.currentTarget);

    const email = data.get("email")
    const password = data.get("password")
    // eslint-disable-next-line no-console
    // console.log(data)
    console.log({
      email: data.get("email"),
      password: data.get("password"),
    });

    // UsersAPI.loginUser(email, password)
    //   .then(function (currentUser: any) {
    //     // toast.message("Hey buddy!")
    //     currentUserContext!.fetchCurrentUser({
    //       uuid: currentUser["uuid"],
    //       email: currentUser["email"],
    //       displayName: currentUser["displayName"],
    //       role: currentUser["role"],
    //       balance: currentUser["balance"],
    //     })
    //     // router.push("/");
    //   })
    //   .catch((error: unknown) => {
    //     // toast.error("Failed to log in")
    //     setShowInvalidLoginNotification(true);
    //     handleCloseBackdrop()
    //     console.error(error);
    //   });

      UsersAPI.loginUser(email, password, (response: any) => {
        console.log(response);
        if(response.statusCode === 200){
          router.push("/");
        }
        else {
          // toast.error("Failed to log in")
          setShowInvalidLoginNotification(true);
          handleCloseBackdrop()
        }
      })
  };

  return (
    <React.Fragment>
      <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
        open={showLoggedOutNotification}
        onClose={handleCloseLoggedOutNotification}
        autoHideDuration={5000}
        message={showLoggedOutNotification ? "You have been logged out" : ""}
        // key={"top" + "center"} 
        action={showLoggedOutNotificationAction}
      />
      <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
        open={showSignedUpNotification}
        onClose={handleCloseSignedUpNotification}
        autoHideDuration={5000}
        message={showSignedUpNotification ? "You have been signed up successfully!" : ""}
        // key={"top" + "center"} 
        action={showSignedUpNotificationAction}
      />

      <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
        open={showInvalidLoginNotification}
        onClose={handleCloseInvalidLoginNotification}
        autoHideDuration={5000}
        message={showInvalidLoginNotification ? "Invalid email or password" : ""}
        key={"top" + "right"}
        action={showInvalidLoginNotificationAction}
      />

      <Grid container component="main" sx={{ height: "100vh" }}>
        <CssBaseline />
        <Grid
          item
          xs={false}
          sm={4}
          md={7}
          sx={{
            backgroundImage: "url(https://source.unsplash.com/random?shopping)",
            backgroundRepeat: "no-repeat",
            backgroundColor: (t) =>
              t.palette.mode === "light"
                ? t.palette.grey[50]
                : t.palette.grey[900],
            backgroundSize: "cover",
            backgroundPosition: "center",
          }}
        />
        <Grid item xs={12} sm={8} md={5} component={Paper} elevation={6} square>
          <Box
            sx={{
              my: 8,
              mx: 4,
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
            }}
          >
            <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
              <LockOutlinedIcon />
            </Avatar>
            <Typography component="h1" variant="h5">
              Sign in
            </Typography>
            <Box
              component="form"
              noValidate
              onSubmit={handleSubmit}
              sx={{ mt: 1 }}
            >
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                autoFocus
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="Password"
                type="password"
                id="password"
                autoComplete="current-password"
              />
              <FormControlLabel
                control={<Checkbox value="remember" color="primary" />}
                label="Remember me"
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                color="secondary"
                sx={{ mt: 3, mb: 2 }}
              >
                Sign In
              </Button>
              <Grid container>
                {/* <Grid item xs>
                  <Link href="#" variant="body2">
                    Forgot password?
                  </Link>
                </Grid> */}
                <Grid item>
                  <Link href="/signup">
                    <Typography variant="body2" color="text.secondary">
                      Don't have an account? Sign up!
                    </Typography>
                  </Link>
                </Grid>
              </Grid>
              <Copyright sx={{ mt: 5 }} />
            </Box>
          </Box>
        </Grid>
      </Grid>
        <Backdrop
          sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
          open={openBackdrop}
        >
          <CircularProgress color="inherit" />
        </Backdrop>
    </React.Fragment>
  );
}
