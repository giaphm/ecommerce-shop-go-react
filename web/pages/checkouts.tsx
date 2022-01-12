import * as React from 'react';
import Router from "next/router";

import StripeCheckout from "react-stripe-checkout";

import {
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

import { styled } from '@mui/material/styles';
import {
  Grid,
  Paper,
  List,
  ListItem,
  ListItemText,
  Backdrop,
  CircularProgress,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
} from "@mui/material"

import Layout from "../components/layout";
import CurrentUserAppCtx from "../store/current-user-context";

import { Auth, setApiClientsAuth } from "../src/repositories/auth";
import * as CheckoutsAPI from "../src/repositories/checkouts";
import * as OrdersAPI from "../src/repositories/orders";
import * as ProductsAPI from "../src/repositories/products";
import * as UsersAPI from "../src/repositories/users";

import { formatDateTime } from "../src/date";

const Item = styled(Paper)(({ theme }) => ({
  ...theme.typography.body2,
  padding: theme.spacing(1),
  textAlign: 'center',
  color: theme.palette.text.secondary,
  wordWrap: "break-word",
  height: "100%",
}));

interface Checkout {
  uuid: string;
  userUuid: string;
  orderUuid: string;
  totalPrice: number;
  notes: string;
  proposedTime: Date;
}

const Checkouts = () => {
  const [checkouts, setCheckouts] = React.useState<Checkout[]>([]);
  const [products, setProducts] = React.useState<any[]>([]);
  
  const currentUserAppCtx = React.useContext(CurrentUserAppCtx);

  console.log("checkouts", checkouts);

  console.log("currentUserAppCtx", currentUserAppCtx);

  console.log("OrdersAPI.ordersClient", OrdersAPI.ordersClient);

  React.useEffect(() => {
    console.log(typeof window);
    console.log("currentUserAppCtx", currentUserAppCtx);
    // const mockUserLoggedIn = JSON.parse(localStorage.getItem("_mock_user") || "{}")
    const isCurrentUserLoggedIn = Auth.isLoggedIn();
    console.log("isCurrentUserLoggedIn", isCurrentUserLoggedIn);
    if (isCurrentUserLoggedIn) {
      const currentUser = Auth.currentUser();
      console.log("currentUser", currentUser);
      Auth.waitForAuthReady()
        .then(() => {
          return Auth.getJwtToken(false)
        })
        .then((token: string) => setApiClientsAuth(token))
        .then(() => {
          // toast.message("Hey buddy!")
          // "go inside hooks update currentUserAppCtx"
          UsersAPI.getCurrentUser((currentUser: any) =>{
            console.log("currentUser", currentUser)
            currentUserAppCtx!.fetchCurrentUser({
              uuid: currentUser["uuid"],
              email: currentUser["email"],
              displayName: currentUser["displayName"],
              role: currentUser["role"],
              balance: currentUser["balance"],
            });
            console.log("currentUserAppCtx", currentUserAppCtx)
            Auth.login(currentUser) // set to local storage
            .then(() => Auth.waitForAuthReady())
            .then(() => {
                return Auth.getJwtToken(false)
            })
            .then((token: any) => setApiClientsAuth(token))
          })
          console.log("currentUserAppCtx", currentUserAppCtx);

          CheckoutsAPI.getUserCheckouts(currentUser.uuid, (userCheckouts: any) => {
            console.log("userCheckouts", userCheckouts); 

            OrdersAPI.getUserOrders(currentUser.uuid, (userOrders: any) => {
                const checkoutsWithTitle: any[] = [];
                console.log("userOrders", userOrders);
                userCheckouts.forEach((userCheckout: any) => {
                  userOrders.forEach((userOrder: any) => {
                      if (userOrder.uuid === userCheckout.orderUuid) {
                        checkoutsWithTitle.push({
                          ...userCheckout,
                          totalPrice: userOrder.totalPrice,
                        });
                      }
                  });
                });
                console.log("setOrders(ordersWithTitle)");
                setCheckouts(checkoutsWithTitle);
  
            });
            
          })
        })
    } else if (!currentUserAppCtx!["uuid"]) {
      Router.push("/login");
    }
  }, []);

  React.useEffect(() => {
    console.log("currentUserAppCtx", currentUserAppCtx);
    
    // ProductsAPI.getProducts((products: any) => {
    //   console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);
    //   console.log(
    //     "ProductsAPI.productsClient.authentications",
    //     ProductsAPI.productsClient.authentications
    //   );
    //   console.log("products", products);
    //   const ordersWithTitle: any[] = [];
    //   console.log("orders", orders);
    //   orders.map((order: any) => {
    //     // console.log("product", product);
    //     const orderItemsWithTitle: any[] = [];
    //     order.orderItems.forEach((orderItem: any) => {
    //       products.forEach((product: any) => {
    //         // console.log("user", user);
    //         if (orderItem.productUuid === product.uuid) {

    //           orderItemsWithTitle.push({
    //             productTitle: product.title,
    //             ...orderItem,
    //           });
    //         }
    //       });
    //     });
    //     ordersWithTitle.push({
    //       orderItems: orderItemsWithTitle,
    //       ...order,
    //     })
    //   });
    //   console.log("ordersWithTitle", ordersWithTitle);
    //   console.log("setOrders(ordersWithTitle)");
    //   setOrders(ordersWithTitle);
    // });
    
  }, [currentUserAppCtx]);

  return (
    <Layout>
      <Grid container sx={{mt: 5, mb: 3}} spacing={2}>
        {/* order uuid */}
        <Grid item xs={2} md={2}>
          <Item>Uuid</Item>
        </Grid>
        {/* orderItems */}
        <Grid item xs={2} md={2}>
          <Item>Order uuid</Item>
        </Grid>
        {/* totalPrice */}
        <Grid item xs={2} md={2}>
          <Item>Total price</Item>
        </Grid>
        {/* status */}
        <Grid item xs={4} md={4}>
          <Item>Note</Item>
        </Grid>
        {/* proposedTime */}
        <Grid item xs={2} md={2}>
          <Item>Proposed time</Item>
        </Grid>
      </Grid>
      <Divider />
      <Grid sx={{mt: 3,}} container spacing={2}>
        {checkouts.map((checkout: Checkout) => {
          console.log("checkout", checkout)
          return (
            <React.Fragment key={checkout.uuid}>
              <Grid item xs={2} md={2}>
                <Item>{checkout.uuid}</Item>
              </Grid>
              <Grid item xs={2} md={2}>
                <Item>{checkout.orderUuid}</Item>
              </Grid>
              <Grid item xs={2} md={2}>
                <Item>{checkout.totalPrice}</Item>
              </Grid>
              <Grid item xs={4} md={4}>
                <Item>{checkout.notes}</Item>
              </Grid>
              <Grid item xs={2} md={2}>
                <Item>{formatDateTime(checkout.proposedTime)}</Item>
              </Grid>
            </React.Fragment>
          )
        })}
      </Grid>
    </Layout>
  );
};
export default Checkouts;