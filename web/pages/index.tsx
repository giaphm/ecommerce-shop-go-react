import * as React from "react";
import type { NextPage } from "next";
import Router from "next/router";

import { Scrollbar } from "smooth-scrollbar-react";
import {
  Button,
  List,
  ListItem,
  IconButton,
  ListItemAvatar,
  Avatar,
  ListItemText,
  ListSubheader,
} from "@mui/material";

import { styled } from "@mui/material/styles";
import { makeStyles } from "@mui/styles";

import DeleteIcon from "@mui/icons-material/Delete";
import FolderIcon from "@mui/icons-material/Folder";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import CardMedia from "@mui/material/CardMedia";
import CssBaseline from "@mui/material/CssBaseline";
import Grid from "@mui/material/Grid";
import Stack from "@mui/material/Stack";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";
import Link from "@mui/material/Link";

import { isEmpty } from "lodash";

import * as ProductsAPI from "../src/repositories/products";
import * as UsersAPI from "../src/repositories/users";
import * as OrdersAPI from "../src/repositories/orders";
import { Auth, setApiClientsAuth } from "../src/repositories/auth";

import Layout from "../components/layout";
import CurrentUserAppCtx from "../store/current-user-context";

const Paper = styled("div")(({ theme }) => ({
  backgroundColor: theme.palette.primary.main,
}));

interface OrderItemToRequest {
  productUuid: string;
  quantity: number;
}

interface OrderToRequest {
  userUuid: string | null;
  orderItems: OrderItemToRequest[];
  totalPrice: number;
}

interface OrderItem {
  productUuid: string;
  productTitle: string;
  productImage: string;
  productPrice: number;
  quantity: number;
}

interface Order {
  orderItems: OrderItem[];
  totalPrice: number;
}

const useStyles = makeStyles({
  customScrollBar: {
    '&::-webkit-scrollbar': {
      width: '0.1em',
    },
    '&::-webkit-scrollbar-track': {
      boxShadow: 'inset 0 0 6px rgba(0,0,0,0.00)',
      webkitBoxShadow: 'inset 0 0 6px rgba(0,0,0,0.00)',
    },
    '&::-webkit-scrollbar-thumb': {
      backgroundColor: 'rgba(0,0,0,.1)',
      // outline: '1px solid slategrey',
    },
  }
})

const Home: NextPage = () => {
  const [users, setUsers] = React.useState([]);
  const [products, setProducts] = React.useState<any[]>([]);
  // const [orderItems, setOrderItems] = React.useState<OrderItem[]>([]);
  const [order, setOrder] = React.useState<Order>({
    orderItems: [],
    totalPrice: 0.0,
  });
  const [isLoading, setIsLoading] = React.useState(false);

  const classes = useStyles();

  console.log("users", users);
  console.log("products", products);
  // console.log("orderItems", orderItems)
  console.log("order", order);

  const currentUserAppCtx = React.useContext(CurrentUserAppCtx);
  console.log("alo");
  console.log(typeof window);

  console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);

  console.log("UsersAPI.usersClient", UsersAPI.usersClient);

  React.useEffect(() => {
    console.log(typeof window);
    console.log("currentUserAppCtx", currentUserAppCtx);
    // const mockUserLoggedIn = JSON.parse(localStorage.getItem("_mock_user") || "{}")
    const isCurrentUserLoggedIn = Auth.isLoggedIn();
    console.log("isCurrentUserLoggedIn", isCurrentUserLoggedIn);
    if (isCurrentUserLoggedIn) {
      const currentUser = Auth.currentUser();
      console.log("currentUser", currentUser);
      UsersAPI.loginUser(currentUser["email"], currentUser["password"])
        .then(function (currentUser: any) {
          // toast.message("Hey buddy!")
          // "go inside hooks update currentUserAppCtx"
          currentUserAppCtx!.fetchCurrentUser({
            uuid: currentUser["uuid"],
            email: currentUser["email"],
            displayName: currentUser["displayName"],
            role: currentUser["role"],
          });
          console.log("currentUserAppCtx", currentUserAppCtx);
        })
        .catch((error: unknown) => {
          // toast.error("Failed to log in")
          console.error(error);
        });

      // // set token in header again
      // Auth.getJwtToken(false).then((token: any) => {
      //   console.log("token", token)
      //   setApiClientsAuth(token)
      // })

      // console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);

      // console.log("UsersAPI.usersClient", UsersAPI.usersClient);

      // console.log("LoggedIn and set currentUserAppCtx again")
      // // fetchCurrentUser is async => new useEffect to chase the updates
      // const currentUser = Auth.currentUser();
      // console.log("currentUser", currentUser);
      // currentUserAppCtx!.fetchCurrentUser({
      //   uuid: currentUser["uuid"],
      //   email: currentUser["email"],
      //   displayName: currentUser["name"],
      //   role: currentUser["role"],
      // });
      console.log("setIsLoading(true);");
      setIsLoading(true);
    } else if (!currentUserAppCtx!["uuid"]) {
      Router.push("/login");
    }
  }, []);

  React.useEffect(() => {
    console.log("currentUserAppCtx", currentUserAppCtx);

    // "go inside hooks update users"
    UsersAPI.getUsers((users: any) => {
      console.log("UsersAPI.usersClient", UsersAPI.usersClient);
      console.log(
        "UsersAPI.usersClient.authentications",
        UsersAPI.usersClient.authentications
      );
      console.log("users", users);
      console.log("setUsers(users);");
      setUsers(users);
    });
  }, [currentUserAppCtx]);

  React.useEffect(() => {
    console.log("currentUserAppCtx", currentUserAppCtx);

    ProductsAPI.getProducts((products: any) => {
      console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);
      console.log(
        "ProductsAPI.productsClient.authentications",
        ProductsAPI.productsClient.authentications
      );
      console.log("products", products);
      const productsWithDisplayName: any[] = [];
      products.map((product: any) => {
        // console.log("product", product);
        users.forEach((user: any) => {
          // console.log("user", user);
          if (user.uuid === product.userUuid) {
            productsWithDisplayName.push({
              displayName: user.displayName,
              ...product,
            });
          }
        });
      });
      console.log("setProducts(productsWithDisplayName)");
      setProducts(productsWithDisplayName);
    });

    console.log("setIsLoading(true);");
    setIsLoading(true);
  }, [users]);

  const viewProductHandler = (productUuid: string) => {
    Router.push({
      pathname: `/product/view/${productUuid}`,
    });
  };

  const addProductToOrderHandler = (
    productUuid: string,
    productTitle: string,
    productImage: string,
    productPrice: number
  ) => {
    setOrder((order: Order) => {
      const newOrderItems = [...order.orderItems];
      if (newOrderItems.length === 0) {
        const orderItem: OrderItem = {
          productUuid,
          productTitle,
          productImage,
          productPrice,
          quantity: 1,
        };
        return {
          ...order,
          orderItems: [orderItem],
          totalPrice: orderItem.productPrice,
        };
      }

      let foundOrderItem = false;
      newOrderItems.map((orderItem) => {
        console.log("orderItem", orderItem);
        if (orderItem.productUuid === productUuid) {
          orderItem.quantity += 1;
          foundOrderItem = true;
        }
        return orderItem;
      });

      if (foundOrderItem) {
        return {
          ...order,
          orderItems: newOrderItems,
          totalPrice: order.totalPrice + productPrice,
        };
      }

      const orderItem: OrderItem = {
        productUuid,
        productTitle,
        productImage,
        productPrice,
        quantity: 1,
      };
      return {
        ...order,
        orderItems: [...order.orderItems, orderItem],
        totalPrice: order.totalPrice + productPrice,
      };
    });
  };

  const removeProductInOrderHandler = (productUuid: string) => {
    setOrder((order: Order) => {
      let foundOrderItem = false;
      const newOrderItems = [...order.orderItems];
      let newTotalPrice = order.totalPrice;
      newOrderItems.map((orderItem) => {
        console.log("orderItem", orderItem);
        if (orderItem.productUuid === productUuid) {
          if (orderItem.quantity > 0) {
            orderItem.quantity -= 1;
            foundOrderItem = true;
          }
        }

        return orderItem;
      });

      if (foundOrderItem) {
        // substract totalPrice
        newOrderItems.forEach((orderItem) => {
          if (orderItem.productUuid === productUuid) {
            newTotalPrice -= orderItem.productPrice;
          }
        });
      }
      return {
        orderItems: newOrderItems.filter(
          (newOrderItem) => newOrderItem.quantity > 0
        ),
        totalPrice: newTotalPrice,
      };
    });
  };

  const addOrderHandler = () => {
    console.log("order", order);
    console.log("currentUserAppCtx", currentUserAppCtx);
    const orderToRequest: OrderToRequest = {
      userUuid: currentUserAppCtx!.uuid,
      orderItems: [],
      totalPrice: order.totalPrice,
    }
    console.log("orderToRequest", orderToRequest);
    order.orderItems.forEach((orderItem) => {
      console.log("orderItem", orderItem);
      orderToRequest.orderItems.push({
        productUuid: orderItem.productUuid,
        quantity: orderItem.quantity,
      })
    })
    console.log("orderToRequest", orderToRequest);
    OrdersAPI.createOrder(
      orderToRequest.userUuid,
      orderToRequest.orderItems,
      orderToRequest.totalPrice
    );
  }

  return isLoading ? (
    <Layout role={currentUserAppCtx!.role}>
      <main>
        {/* Hero unit */}
        <Box
          sx={{
            bgcolor: "background.paper",
            pt: 8,
            pb: 6,
          }}
        >
          <Container maxWidth="sm">
            <Typography
              component="h1"
              variant="h2"
              align="center"
              color="text.primary"
              gutterBottom
            >
              E-commerce shop
            </Typography>
            {/* <Typography variant="h5" align="center" color="text.secondary" paragraph>
                Something short and leading about the collection below—its contents,
                the creator, etc. Make it short and sweet, but not too short so folks
                don&apos;t simply skip over it entirely.
              </Typography> */}
            {/* <Stack
                sx={{ pt: 4 }}
                direction="row"
                spacing={2}
                justifyContent="center"
              >
                <Button
                  variant="contained"
                  onClick={postNewProduct}
                >
                  Post a new product
                </Button>
                <Button variant="outlined">Secondary action</Button>
              </Stack> */}
          </Container>
        </Box>
        <Container sx={{ py: 8 }} maxWidth="md">
          {/* End hero unit */}
          <Grid container spacing={4}>
            {products.map((product: any) => (
              <Grid item key={product.uuid} xs={12} sm={6} md={4}>
                <Card
                  sx={{
                    height: "100%",
                    display: "flex",
                    flexDirection: "column",
                  }}
                >
                  <CardMedia
                    component="img"
                    sx={{
                      // 16:9
                      pt: "56.25%",
                    }}
                    image={product.image}
                    alt="random"
                  />
                  <CardContent sx={{ flexGrow: 1 }}>
                    <Typography gutterBottom variant="h5" component="h2">
                      {product.title}
                    </Typography>
                    <Typography>Owner: {product.displayName}</Typography>
                    <Typography>Description: {product.description}</Typography>
                    <Typography>Price: ${product.price}</Typography>
                    <Typography>Quantity: {product.quantity}</Typography>
                  </CardContent>
                  <CardActions>
                    <Button
                      size="small"
                      onClick={() => viewProductHandler(product.uuid)}
                    >
                      View
                    </Button>

                    {currentUserAppCtx!.uuid === product.userUuid ? (
                      ""
                    ) : (
                      <Button
                        size="small"
                        onClick={() => {
                          addProductToOrderHandler(
                            product.uuid,
                            product.title,
                            product.image,
                            product.price
                          );
                        }}
                      >
                        Add to your order
                      </Button>
                    )}
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
        {order.orderItems.length > 0 ? (
          <Box
            sx={{
              position: "fixed",
              top: "65%",
              right: "20px",
              maxWidth: 150,
              minWidth: "100px",
            }}
          >
            <Grid item xs={12} md={12}>
              <Typography
                align="center"
                sx={{
                  bgcolor: "primary.main",
                  border: "1px dashed black",
                  mt: 1,
                  pb: 1,
                }}
                variant="h6"
                component="div"
              >
                Your order
              </Typography>
              <Paper>
                <List
                  className={classes.customScrollBar}
                  dense={false}
                  sx={{
                    width: "100%",
                    maxWidth: 150,
                    bgcolor: "background.paper",
                    position: "relative",
                    overflow: "auto",
                    maxHeight: 100,
                    "& ul": { padding: 0 },
                  }}
                >
                  {/* <Scrollbar
                    plugins={{
                      overscroll: {
                        effect: "bounce", 
                      } as const,
                    }}
                  > */}
                    {order.orderItems.map((orderItem) => (
                      <ListItem
                        key={orderItem.productUuid}
                        sx={{
                          width: "100%",
                          maxWidth: 150,
                          // px: 0,
                          maxHeight: "40px",
                        }}
                        secondaryAction={
                          <IconButton
                            edge="end"
                            aria-label="delete"
                            onClick={() => {
                              removeProductInOrderHandler(
                                orderItem.productUuid
                              );
                            }}
                          >
                            <DeleteIcon />
                          </IconButton>
                        }
                      >
                        <ListItemAvatar
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
                        </ListItemAvatar>
                        <ListItemText
                          sx={{
                            width: "100%",
                            maxWidth: 150,
                            // px: 0,
                            maxHeight: "60px",
                          }}
                          primary={`${orderItem.productTitle} x${orderItem.quantity}`}
                          primaryTypographyProps={{
                            fontSize: "0.6rem",
                          }}
                          secondary={null}
                        />
                      </ListItem>
                    ))}
                  {/* </Scrollbar> */}
                </List>
                {/* button add order */}
                {/* button cancel order in ui */}
                <Grid
                  sx={{ width: "100%", maxHeight: "30px" }}
                  item
                  xs={12}
                  md={12}
                >
                  <Stack
                    sx={{ width: "100%", maxHeight: "30px" }}
                    spacing={0}
                    direction="row"
                  >
                    <Button
                      sx={{
                        minWidth: "60%",
                        fontSize: "10px",
                        borderRadius: 0,
                      }}
                      variant="contained"
                      onClick={addOrderHandler}
                    >
                      Add
                    </Button>
                    <Button
                      color="warning"
                      sx={{
                        minWidth: "30%",
                        fontSize: "10px",
                        borderRadius: 0,
                      }}
                      variant="contained"
                    >
                      Cancel
                    </Button>
                  </Stack>
                </Grid>
              </Paper>
            </Grid>
          </Box>
        ) : (
          ""
        )}
      </main>
    </Layout>
  ) : (
    <div>Loading...</div>
  );
};

export default Home;
