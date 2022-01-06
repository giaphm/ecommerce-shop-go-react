import * as React from "react";
import { useRouter } from "next/router";

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
import { createTheme, ThemeProvider } from "@mui/material/styles";

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
  const [showLoader, setShowLoader] = React.useState(false);
  const [showLoggedOutNotification, setShowLoggedOutNotification] = React.useState(false);

  const currentUserContext = React.useContext(CurrentUserAppCtx);
  console.log("currentUserContext", currentUserContext)

  const router = useRouter()
  const { query } = router;
  console.log("query", query)

  let loggedOut: boolean = false;

  if (!isEmpty(query)) {
    loggedOut = JSON.parse(query["loggedOut"] as string);
  }

  console.log("loggedOut", loggedOut)

  React.useEffect(() => {
    if (Auth.isLoggedIn()) {
      router.push("/");
    }
    setShowLoggedOutNotification(loggedOut);
  }, []);

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const data = new FormData(event.currentTarget);

    const email = data.get("email")
    const password = data.get("password")
    // eslint-disable-next-line no-console
    // console.log(data)
    console.log({
      email: data.get("email"),
      password: data.get("password"),
    });
    setShowLoader(true);

    UsersAPI.loginUser(email, password)
      .then(function (currentUser: any) {
        // toast.message("Hey buddy!")
        currentUserContext!.fetchCurrentUser({
          uuid: currentUser["uuid"],
          email: currentUser["email"],
          displayName: currentUser["displayName"],
          role: currentUser["role"],
        })
        router.push("/");
      })
      .catch((error: unknown) => {
        // toast.error("Failed to log in")
        console.error(error);
        setShowLoader(false);
      });
  };

  return (
    <React.Fragment>
      <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
        open={showLoggedOutNotification}
        // onClose={handleClose}
        message={showLoggedOutNotification ? "You have been logged out" : ""}
        key={"top" + "right"}
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
                sx={{ mt: 3, mb: 2 }}
              >
                Sign In
              </Button>
              <Grid container>
                <Grid item xs>
                  <Link href="#" variant="body2">
                    Forgot password?
                  </Link>
                </Grid>
                <Grid item>
                  <Link href="#" variant="body2">
                    {"Don't have an account? Sign Up"}
                  </Link>
                </Grid>
              </Grid>
              <Copyright sx={{ mt: 5 }} />
            </Box>
          </Box>
        </Grid>
      </Grid>
    </React.Fragment>
  );
}
