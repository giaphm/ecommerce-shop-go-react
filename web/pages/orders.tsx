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

interface OrderItem {
  uuid: string;
  productUuid: string;
  quantity: number;
}

interface Order {
  uuid: string;
  userUuid: string;
  orderItems: any[]; // to add the productTitle property
  totalPrice: number;
  status: string;
  proposedTime: Date;
  expiresAt: Date;
}

const Orders = () => {
  const [orders, setOrders] = React.useState<Order[]>([]);
  const [products, setProducts] = React.useState<any[]>([]);
  const [openDialog, setOpenDialog] = React.useState(false);
  const [openBackdrop, setOpenBackdrop] = React.useState(false);
  const [note, setNote] = React.useState("");
  const [orderToRequest, setOrderToRequest] = React.useState<Order>();
  
  const currentUserAppCtx = React.useContext(CurrentUserAppCtx);

  console.log("orders", orders);

  console.log("openDialog", openDialog);

  console.log("note", note);
  
  console.log("orderToRequest", orderToRequest);

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
          console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);
        
          console.log("UsersAPI.usersClient", UsersAPI.usersClient);
          console.log("OrdersAPI.ordersClient", OrdersAPI.ordersClient);
          // toast.message("Hey buddy!")
          UsersAPI.getCurrentUser((currentUser: any) => {
            currentUserAppCtx!.fetchCurrentUser({
              uuid: currentUser["uuid"],
              email: currentUser["email"],
              displayName: currentUser["displayName"],
              role: currentUser["role"],
              balance: currentUser["balance"],
            });
            console.log("currentUserAppCtx", currentUserAppCtx);
            console.log("currentUser", currentUser);
          });
          setOrdersWithTitle(currentUser.uuid);
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

  const setOrdersWithTitle = (currentUserUuid: string) => {
    
    OrdersAPI.getUserOrders(currentUserUuid, (userOrders: any) => {
      console.log("userOrders", userOrders); 
      // setOrders(userOrders);
      ProductsAPI.getProducts((products: any) => {
        console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);
        console.log(
          "ProductsAPI.productsClient.authentications",
          ProductsAPI.productsClient.authentications
        );
        console.log("products", products);
        const ordersWithTitle: any[] = [];
        console.log("userOrders", userOrders);
        userOrders.map((order: any) => {
          // console.log("product", product);
          const orderItemsWithTitle: any[] = [];
          order.orderItems.forEach((orderItem: any) => {
            products.forEach((product: any) => {
              // console.log("user", user);
              if (orderItem.productUuid === product.uuid) {
                console.log("product", product);
                orderItemsWithTitle.push({
                  productTitle: product.title,
                  ...orderItem,
                });
              }
            });
          });
          console.log("orderItemsWithTitle", orderItemsWithTitle)
          ordersWithTitle.push({
            ...order,
            orderItems: [...orderItemsWithTitle],
          })
        });
        console.log("ordersWithTitle", ordersWithTitle);
        console.log("setOrders(ordersWithTitle)");
        setOrders(ordersWithTitle);
      });
    });
  }

  const handleOpenBackdrop = () => {
    setOpenBackdrop(true);
  }
  
  const handleCloseBackdrop = () => {
    setOpenBackdrop(false);
  }

  const handleCloseDialog = () => {
    setOpenDialog(false);
  };

  const openNoteHandler = (order: Order) => {
    setOrderToRequest(order);
    setOpenDialog(true);
    console.log("openDialog", openDialog);
  }

  const cancelOrderHandler = (orderUuid: string) => {
    OrdersAPI.cancelOrder(orderUuid, (response: any) => {
      setOrdersWithTitle(currentUserAppCtx!.uuid as string);
    });
  }

  const payOrderHandler = (tokenId: string) => {
    const newDate = new Date();
    console.log("newDate", newDate)
    handleCloseDialog();
    setTimeout(() => handleOpenBackdrop(), 1000);
    CheckoutsAPI.createCheckout(orderToRequest?.uuid, note, newDate, tokenId, (response: any) => {
      console.log("response", response)
      if(response.statusCode === 201) {
        Router.push("/checkouts");
        handleCloseBackdrop()
      }
      else {
        console.log("failed")
        handleCloseBackdrop();
      }
    });
  }

  return (
    <Layout>
      <Grid sx={{mt: 5, mb: 3}} container spacing={2}>
        {/* order uuid */}
        <Grid item xs={1} md={1}>
          <Item>Uuid</Item>
        </Grid>
        {/* orderItems */}
        <Grid item xs={2} md={2}>
          <Item>Order items</Item>
        </Grid>
        {/* totalPrice */}
        <Grid item xs={1} md={1}>
          <Item>Total price</Item>
        </Grid>
        {/* status */}
        <Grid item xs={1.5} md={1.5}>
          <Item>Status</Item>
        </Grid>
        {/* proposedTime */}
        <Grid item xs={2} md={2}>
          <Item>Proposed time</Item>
        </Grid>
        {/* expiresAt */}
        <Grid item xs={2} md={2}>
          <Item>Expires At</Item>
        </Grid>
        {/* actions */}
        <Grid item xs={2.5} md={2.5}>
          <Item>Actions</Item>
        </Grid>
      </Grid>
      <Divider />
      <Grid sx={{mt: 3,}} container spacing={2}>
        {orders.map((order: Order) => {
          console.log("order", order)
          return (
            <React.Fragment key={order.uuid}>
              <Grid item xs={1} md={1}>
                <Item>{order.uuid}</Item>
              </Grid>
              <Grid item xs={2} md={2}>
                
              <Item>
                <List
                    // className={classes.customScrollBar}
                    dense={false}
                    sx={{
                      pt: 0,
                      pb: 0,
                      width: "100%",
                      bgcolor: "background.paper",
                      position: "relative",
                      // overflow: "auto",
                      // maxHeight: 100,
                      "& ul": { padding: 0 },
                    }}
                  >
                      {order.orderItems.map((orderItem) => (
                        <ListItem
                          key={orderItem.uuid}
                          sx={{
                            py: 0,
                            width: "100%",
                            // px: 0,
                          }}
                          // secondaryAction={
                          //   <IconButton
                          //     edge="end"
                          //     aria-label="delete"
                          //     onClick={() => {
                          //       removeProductInOrderHandler(
                          //         orderItem.productUuid
                          //       );
                          //     }}
                          //   >
                          //     <DeleteIcon />
                          //   </IconButton>
                          // }
                        >
                          {/* <ListItemAvatar
                            sx={{
                              width: "100%",
                              maxWidth: 50,
                              minWidth: 20,
                              // maxHeight: "30px",
                            }}
                          >
                            <Avatar
                              sx={{
                                maxWidth: 20,
                              }}
                              src={orderItem.productImage}
                            />
                          </ListItemAvatar> */}
                          <ListItemText
                            sx={{
                              width: "100%",
                              maxWidth: 150,
                              // px: 0,
                              maxHeight: "60px",
                            }}
                            primary={`${orderItem.productTitle} - x${orderItem.quantity}`}
                            primaryTypographyProps={{
                              fontSize: "0.875rem",
                            }}
                            secondary={null}
                          />
                        </ListItem>
                      ))}
                  </List>
                </Item>
                {/* <ul>{order.orderItems.map((orderItem: OrderItem) => (
                  <div>
                    <li>{orderItem.uuid}</li>
                    <li>{orderItem.productUuid}</li>
                    <li>{orderItem.quantity}</li>
                  </div>
                ))}</ul> */}
              </Grid>
              <Grid item xs={1} md={1}>
                <Item>${order.totalPrice}</Item>
              </Grid>
              <Grid item xs={1.5} md={1.5}>
                <Item>{order.status}</Item>
              </Grid>
              <Grid item xs={2} md={2}>
                <Item>{formatDateTime(order.proposedTime)}</Item>
              </Grid>
              <Grid item xs={2} md={2}>
                <Item>{formatDateTime(order.expiresAt)}</Item>
              </Grid>
              {/* payment action */}
              <Grid item xs={1} md={1}>
                    <Button
                      // sx={{
                      //   minWidth: "60%",
                      //   fontSize: "10px",
                      //   borderRadius: 0,
                      // }}
                      variant="contained"
                      color="success"
                      disabled={order.status !== "created"}
                      onClick={() => openNoteHandler(order)}
                    >
                      Pay
                    </Button>
              </Grid>
              {/* cancel action */}
              <Grid item xs={1.5} md={1.5}>
                    <Button
                      color="warning"
                      // sx={{
                      //   minWidth: "30%",
                      //   fontSize: "10px",
                      //   borderRadius: 0,
                      // }}
                      variant="contained"
                      disabled={order.status !== "created"}
                      onClick={() => cancelOrderHandler(order.uuid)}
                    >
                      Cancel
                    </Button>
              </Grid>
            </React.Fragment>
          )
        })}
      </Grid>
      <Dialog open={openDialog} onClose={handleCloseDialog}>
        <DialogTitle>Order's Note</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Please enter the note for your order, it's optional.
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="note"
            label="Note"
            type="text"
            fullWidth
            variant="standard"
            onChange={(e) => setNote(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button color="warning" onClick={handleCloseDialog}>Cancel</Button>
          {/* <Button onClick={() => payOrderHandler()}>Finish</Button> */}
          {orderToRequest ? ( 
            <StripeCheckout
            token={(token: any) => {
              console.log("token", token)
              payOrderHandler(token.id)
            }}
            stripeKey={
              "pk_test_51JyZ0sDfLgIVzReB6QjXOl8oKXIZK9enFlA206neSNKBzfy09xDPcNzqzXQqzQxch7anuT8XweOc4fQUHYqaP2MZ008P9i8qt1"
            }
            amount={orderToRequest!.totalPrice * 100}
            email={currentUserAppCtx!.email as string}
            >
              <Button color="success">Finish</Button>
            </StripeCheckout>
          ) : ""}
          
        </DialogActions>
      </Dialog>
      <Backdrop
        sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
        open={openBackdrop}
      >
        <CircularProgress color="inherit" />
      </Backdrop>
    </Layout>
  );
};
export default Orders;