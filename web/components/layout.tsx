import React from "react";
import { Fragment } from "react";
import Head from "next/head";
import Router from "next/router";

import {
  Backdrop,
  CircularProgress,
  Divider,
} from "@mui/material";

import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import Menu from '@mui/material/Menu';
import MenuIcon from '@mui/icons-material/Menu';
import Container from '@mui/material/Container';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import Tooltip from '@mui/material/Tooltip';
import MenuItem from '@mui/material/MenuItem';
import Link from "@mui/material/Link";

import AccountCircleOutlinedIcon from '@mui/icons-material/AccountCircleOutlined';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import ProductionQuantityLimitsIcon from '@mui/icons-material/ProductionQuantityLimits';

import { Auth } from "../src/repositories/auth";
import CurrentUserAppCtx from "../store/current-user-context";


function Copyright() {
  return (
    <Typography variant="body2" color="text.secondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="https://mui.com/">
        ecommerce-shop-go-react
      </Link>{' '}
      {new Date().getFullYear() - 1} - {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}


const pages = ['Products', 'Orders', 'Checkouts'];
const settings = ['Profile', 'Logout'];

const Layout = (props: any) => {
  const [anchorElNav, setAnchorElNav] = React.useState<(EventTarget & Element) | null>(null);
  const [anchorElUser, setAnchorElUser] = React.useState<(EventTarget & Element) | null>(null);
  const [openBackdrop, setOpenBackdrop] = React.useState(false);
  
  const currentUserAppCtx = React.useContext(CurrentUserAppCtx);

  console.log("openBackdrop", openBackdrop)

  console.log("currentUserAppCtx", currentUserAppCtx)

  console.log("props", props)

  const handleOpenBackdrop = () => {
    setOpenBackdrop(true);
  }
  
  const handleCloseBackdrop = () => {
    setOpenBackdrop(false);
  }

  const handleOpenNavMenu = (event: React.SyntheticEvent) => {
    setAnchorElNav(event.currentTarget);
  };
  const handleOpenUserMenu = (event: React.SyntheticEvent) => {
    setAnchorElUser(event.currentTarget);
  };

  const handleCloseNavMenu = () => {
    setAnchorElNav(null);
  };

  const handleCloseUserMenu = () => {
    setAnchorElUser(null);
  };

  const logout = () => {
    currentUserAppCtx!.removeCurrentUser()
    
    Auth.logout().then(() => {
      console.log("Logout successfully")
      Router.push({
        pathname: "/login",
        query: {
          loggedOut: true,
        }
      })
    })
  }

  return (
    <Fragment>
      <Head>
        <link
          href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css"
          rel="stylesheet"
          integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC"
          crossOrigin="anonymous"
        ></link>
        <script src="https://js.stripe.com/v3/"></script>
      </Head>
      <div className="container-fluid">
        
        <AppBar color="secondary" position="static">
          <Container maxWidth="xl">
            <Toolbar disableGutters>
              {/* <Typography
                variant="h6"
                noWrap
                component="div"
                sx={{ mr: 2, display: { xs: 'none', md: 'flex' } }}
              > */}
                {/* LOGO */}
                {/* <ProductionQuantityLimitsIcon />
              </Typography> */}
    
              <Box sx={{ flexGrow: 1, display: { xs: 'flex', md: 'flex' } }}>
                <IconButton
                  size="large"
                  aria-label="account of current user"
                  aria-controls="menu-appbar"
                  aria-haspopup="true"
                  onClick={handleOpenNavMenu}
                  color="inherit"
                >
                  <MenuIcon />
                </IconButton>
                <Menu
                  id="menu-appbar"
                  anchorEl={anchorElNav}
                  anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                  }}
                  keepMounted
                  transformOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                  }}
                  open={Boolean(anchorElNav)}
                  onClose={handleCloseNavMenu}
                  sx={{
                    display: { xs: 'block', md: 'block' },
                  }}
                >
                  { props.role === "shopkeeper" ? (
                    <div>
                      <MenuItem key={"Your products"} onClick={() => {
                        handleCloseNavMenu();
                        console.log("handleOpenBackdrop();")
                        handleOpenBackdrop();
                        setTimeout(() => Router.push(`/product/your-products`), 500);
                      }}>
                        <Typography textAlign="center">{"Your Products"}</Typography>
                      </MenuItem>
                      <Divider />
                    </div>)
                  : ""}
                  {pages.map((page) => (
                    <MenuItem key={page} onClick={() => {
                      handleCloseNavMenu();
                      console.log("handleOpenBackdrop();")
                      handleOpenBackdrop();
                      setTimeout(() => Router.push(page === "Products" ? `/` : `/${page.toLowerCase()}`), 500)
                    }}>
                      <Typography textAlign="center">{page}</Typography>
                    </MenuItem>
                  ))}
                </Menu>
              </Box>
              {/* <Typography
                variant="h6"
                noWrap
                component="div"
                sx={{ flexGrow: 1, display: { xs: 'flex', md: 'none' } }}
              >
                <ProductionQuantityLimitsIcon />
              </Typography> */}
              {/* <Box sx={{ flexGrow: 1, display: { xs: 'none', md: 'flex' } }}>
                {pages.map((page) => (
                  <Button
                    key={page}
                    onClick={() => {
                      console.log("handleOpenBackdrop();")
                      handleOpenBackdrop();
                      setTimeout(() => Router.push(page === "Products" ? `/` : `/${page.toLowerCase()}`), 500)
                      setTimeout(() => handleCloseBackdrop(), 3000);
                    }}
                    sx={{ my: 2, color: 'white', display: 'block' }}
                  >
                    {page}
                  </Button>
                ))}
              </Box> */}
    
              {/* { props.role === "shopkeeper" ? 
              <Box sx={{ flexGrow: 0 }}>
                <Button
                key={"Your products"}
                onClick={() => Router.push(`/product/your-products`)}
                sx={{ my: 2, color: 'white', display: { xs: 'none', md: 'flex' } }}
                >
                  {"Your products"}
                </Button>
              </Box> : "" } */}
              <Box sx={{ flexGrow: 0 }}>
                <Typography sx={{mr: "8px"}}>
                  {currentUserAppCtx!.displayName ? `Hello, ${currentUserAppCtx!.displayName} !` : ""}
                </Typography>
                {/* <Button
                key={"Your products"}
                onClick={() => Router.push(`/product/your-products`)}
                sx={{ my: 2, color: 'white', display: { xs: 'none', md: 'flex' } }}
                >
                  {"Your products"}
                </Button> */}
              </Box>
              <Box sx={{ flexGrow: 0 }}>
                <Tooltip title="Open settings">
                  <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                    <AccountCircleOutlinedIcon style={{fill: "white"}} fontSize={"large"} />
                  </IconButton>
                </Tooltip>
                <Menu
                  sx={{ mt: '45px' }}
                  id="menu-appbar"
                  anchorEl={anchorElUser}
                  anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                  }}
                  keepMounted
                  transformOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                  }}
                  open={Boolean(anchorElUser)}
                  onClose={handleCloseUserMenu}
                >
                {/* <MenuItem style={{cursor: "context-menu"}}>
                  <Typography>{currentUserAppCtx!.email}</Typography>
                </MenuItem> */}
                  <MenuItem style={{cursor: "context-menu"}}>
                    <Typography>Balance: ${currentUserAppCtx!.balance}</Typography>
                  </MenuItem>
                  {settings.map((setting) => (
                    <MenuItem key={setting} onClick={() => {
                      handleCloseUserMenu()
                      console.log("handleOpenBackdrop();")
                      handleOpenBackdrop();
                      setTimeout(() => setting === "Logout" ? logout() : Router.push(`/${setting.toLowerCase()}`), 500)
                      }
                    }>
                      <Typography textAlign="center">{setting}</Typography>
                    </MenuItem>
                  ))}
                </Menu>
              </Box>
            </Toolbar>
          </Container>
        </AppBar>
        {props.children}
      </div>
        {/* Footer */}
        <Box sx={{ bgcolor: 'background.paper', py: 6 }} component="footer">
          <Divider />
          <Typography sx={{mt: 1}} variant="h6" align="center" gutterBottom>
            Footer
          </Typography>
          <Typography
            variant="subtitle1"
            align="center"
            color="text.secondary"
            component="p"
          >
            Something here to give the footer a purpose!
          </Typography>
          <Copyright />
        </Box>
        {/* End footer */}
        <Backdrop
          sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
          open={openBackdrop}
        >
          <CircularProgress color="inherit" />
        </Backdrop>
    </Fragment>
  );
};

export default Layout;
