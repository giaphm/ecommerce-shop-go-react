import * as React from 'react';
import Router from "next/router";

import {
  Paper,
  InputLabel,
  MenuItem,
  Select,
  FormControl,
  FormHelperText,
  Backdrop,
  CircularProgress,
  Snackbar,
  IconButton,
} from '@mui/material';

import CloseIcon from '@mui/icons-material/Close';

import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';

import * as UsersAPI from "../src/repositories/users";
import { Auth } from "../src/repositories/auth";

function Copyright(props: any) {
  return (
    <Typography variant="body2" color="text.secondary" align="center" {...props}>
      {'Copyright © '}
      <Link color="inherit" href="https://mui.com/">
        ecommerce-shop-go-react
      </Link>{' '}
      {new Date().getFullYear() - 1} - {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}


export default function SignUp() {
  const [role, setRole] = React.useState("");
  const [checkbox, setCheckbox] = React.useState(false);
  const [openBackdrop, setOpenBackdrop] = React.useState(false);
  const [showInvalidSignupNotification, setShowInvalidSignupNotification] = React.useState(false);

  React.useEffect(() => {
    if (Auth.isLoggedIn()) {
      Router.push("/");
    }
  }, []);

  const handleCloseInvalidSignupNotification = () => {
    setShowInvalidSignupNotification(false);
  }

  const showInvalidSignupNotificationAction = (
      <IconButton
        size="small"
        aria-label="close"
        color="inherit"
        onClick={handleCloseInvalidSignupNotification}
      >
        <CloseIcon fontSize="small" />
      </IconButton>
  )

  const handleOpenBackdrop = () => {
    setOpenBackdrop(true);
  }
  
  const handleCloseBackdrop = () => {
    setOpenBackdrop(false);
  }

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    handleOpenBackdrop()
    event.preventDefault();
    console.log(event.currentTarget)
    const data = new FormData(event.currentTarget);
    // eslint-disable-next-line no-console
    console.log({
      displayName: data.get("displayName"),
      role: data.get("role"),
      email: data.get('email'),
      password: data.get('password'),
    });

    const dataToRequest = {
      displayName: data.get("displayName"),
      role: data.get("role"),
      email: data.get('email'),
      password: data.get('password'),
    }

    const { displayName, email, password, role } = dataToRequest;

    console.log("displayName", displayName)
    console.log("email", email)
    console.log("password", password)
    console.log("role", role)

    UsersAPI.signupUser(displayName, email, password, role, (response: any) => {
      console.log("response", response)
      if(response.statusCode === 201) {
        Router.push({
          pathname: "/login",
          query: {
            signedUp: true,
          }
        })
      }
    })
  };

  return (
    <React.Fragment>
      <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
        open={showInvalidSignupNotification}
        onClose={handleCloseInvalidSignupNotification}
        autoHideDuration={5000}
        message={showInvalidSignupNotification ? "Sign up failure" : ""}
        key={"top" + "right"}
        action={showInvalidSignupNotificationAction}
      />

      <Paper style={{backgroundImage: "url(https://source.unsplash.com/random?shopping)"}}>
      <Container style={{backgroundColor: "white"}} sx={{height: "100vh"}} component="main" maxWidth="xs">
        <CssBaseline />
        <Box
          sx={{
            // marginTop: 8,
            pt: 8,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
          }}
        >
          <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component="h1" variant="h5">
            Sign up
          </Typography>
          <Box component="form" noValidate onSubmit={handleSubmit} sx={{ mt: 3 }}>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={8}>
                <TextField
                  autoComplete="given-display-name"
                  name="displayName"
                  required
                  fullWidth
                  id="displayName"
                  label="Display Name"
                  autoFocus
                  color="primary"
                />
              </Grid>
              <Grid item xs={12} sm={4}>
                {/* <TextField
                  required
                  fullWidth
                  id="role"
                  label="Role"
                  name="role"
                  autoComplete="role"
                /> */}
                <FormControl sx={{ minWidth: 120 }} error={!role}>
                  <InputLabel id="role">Role *</InputLabel>
                  <Select
                    labelId="role"
                    value={role}
                    required
                    label="Role"
                    name="role"
                    onChange={(e) => setRole(e.target.value)}
                    renderValue={(value) => `${value}`}
                  >
                    {/* <MenuItem value="">
                      <em>None</em>
                    </MenuItem> */}
                    <MenuItem value={"shopkeeper"}>Shopkeeper</MenuItem>
                    <MenuItem value={"user"}>User</MenuItem>
                  </Select>
                  <FormHelperText>Required</FormHelperText>
                </FormControl>
              </Grid>
              <Grid item xs={12}>
                <TextField
                  required
                  fullWidth
                  id="email"
                  label="Email Address"
                  name="email"
                  autoComplete="email"
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  required
                  fullWidth
                  name="password"
                  label="Password"
                  type="password"
                  id="password"
                  autoComplete="new-password"
                />
              </Grid>
              <Grid item xs={12}>
                <FormControlLabel
                  control={<Checkbox value="allowExtraEmails" color="secondary" />}
                  label="I agree the policy of this website for signing up an account."
                  onChange={(e: React.SyntheticEvent, checked: boolean) => {
                    console.log(checked)
                    // console.log(e.target.checked)
                    setCheckbox(checked)
                  }}
                />
              </Grid>
              <Grid item xs={12}>
                <FormControlLabel
                  control={<Checkbox value="allowExtraEmails" color="secondary" />}
                  label="I want to receive inspiration, marketing promotions and updates via email."
                />
              </Grid>
            </Grid>
            <Button
              type="submit"
              fullWidth
              variant="contained"
              color="secondary"
              sx={{ mt: 3, mb: 2 }}
              disabled={!checkbox}
            >
              Sign Up
            </Button>
            <Grid container justifyContent="flex-end">
              <Grid item>
                <Link href="/login" color="secondary.main" variant="body2">
                  Already have an account? Sign in
                </Link>
              </Grid>
            </Grid>
          </Box>
        </Box>
        <Copyright sx={{ mt: 5 }} />
      </Container>
      </Paper>
        <Backdrop
          sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
          open={openBackdrop}
        >
          <CircularProgress color="inherit" />
        </Backdrop>
    </React.Fragment>
  );
}