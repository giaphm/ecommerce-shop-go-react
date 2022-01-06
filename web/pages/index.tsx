import * as React from "react";
import type { NextPage } from "next";
import Router from "next/router";

import {
  List,
  ListItem,
  IconButton,
  ListItemAvatar,
  Avatar,
  ListItemText,
  ListSubheader,
} from "@mui/material";

import { styled } from '@mui/material/styles';
import DeleteIcon from '@mui/icons-material/Delete';
import FolderIcon from '@mui/icons-material/Folder';
import Button from "@mui/material/Button";
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


import * as ProductsAPI from '../src/repositories/products';
import * as UsersAPI from '../src/repositories/users';
import { Auth, setApiClientsAuth } from '../src/repositories/auth';

import Layout from "../components/layout";
import CurrentUserAppCtx from "../store/current-user-context";

function Copyright() {
  return (
    <Typography variant="body2" color="text.secondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="https://mui.com/">
        Your Website
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

function generate(element: React.ReactElement) {
  return [0, 1, 2, 3, 4].map((value) =>
    React.cloneElement(element, {
      key: value,
    }),
  );
}

const Paper = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.secondary.main,
}));

interface OrderItem {
  productUuid: string;
  productTitle: string;
  productImage: string;
  productPrice: string;
  quantity: number;
}

interface Order {
  productUuids: string[];
  totalPrice: number;
}

const Home: NextPage = () => {
  const [users, setUsers] = React.useState([]);
  const [products, setProducts] = React.useState<any[]>([]);
  const [orderItems, setOrderItems] = React.useState<OrderItem[]>([]);
  const [order, setOrder] = React.useState<Order>({
    productUuids: [],
    totalPrice: 0.0,
  });
  const [isLoading, setIsLoading] = React.useState(false);

  console.log("users", users)
  console.log("products", products)
  console.log("orderItems", orderItems)
  console.log("order", order)
  
  const currentUserAppCtx = React.useContext(CurrentUserAppCtx);
  console.log("alo")
  console.log(typeof window)
  
  console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);

  console.log("UsersAPI.usersClient", UsersAPI.usersClient);

  React.useEffect(() => {
    console.log(typeof window)
    console.log("currentUserAppCtx", currentUserAppCtx)
    // const mockUserLoggedIn = JSON.parse(localStorage.getItem("_mock_user") || "{}")
    const isCurrentUserLoggedIn = Auth.isLoggedIn()
    console.log("isCurrentUserLoggedIn", isCurrentUserLoggedIn)
    if (isCurrentUserLoggedIn) {
      const currentUser = Auth.currentUser();
      console.log("currentUser", currentUser);
      UsersAPI.loginUser(currentUser["email"], currentUser["password"])
        .then(function (currentUser: any) {
          // toast.message("Hey buddy!")
          currentUserAppCtx!.fetchCurrentUser({
            uuid: currentUser["uuid"],
            email: currentUser["email"],
            displayName: currentUser["displayName"],
            role: currentUser["role"],
          })
          Router.push("/");
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
      setIsLoading(true);
    }
    else if (!currentUserAppCtx!["uuid"]) {
      Router.push("/login");
    }
    
    
  }, []);


  React.useEffect(() => {
    console.log("currentUserAppCtx", currentUserAppCtx)

    
    UsersAPI.getUsers((users: any) => {
      console.log("UsersAPI.usersClient", UsersAPI.usersClient);
      console.log("UsersAPI.usersClient.authentications", UsersAPI.usersClient.authentications);
      console.log("users", users);
      setUsers(users);
    })

  }, [currentUserAppCtx])

  React.useEffect(() => {
    console.log("currentUserAppCtx", currentUserAppCtx)
    
    ProductsAPI.getProducts((products: any) => {
      console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);
      console.log("ProductsAPI.productsClient.authentications", ProductsAPI.productsClient.authentications);
      console.log("products", products);
      const productsWithDisplayName: any[] = [];
      products.map((product: any) => {
        // console.log("product", product);
        users.forEach((user: any) => {
          // console.log("user", user);
          if(user.uuid === product.userUuid) {
            productsWithDisplayName.push({
              displayName: user.displayName,
              ...product,
            })
          }
        })
      })
      setProducts(productsWithDisplayName);
    })
  }, [users])

  const viewProductHandler = (productUuid: string) => {
    Router.push({
      pathname: `/product/view/${productUuid}`,
    });
  }

  const addProductToOrderItemsHandler = (
    productUuid: string,
    productTitle: string,
    productImage: string,
    productPrice: string,
  ) => {

    if(orderItems.length === 0){
      const orderItem = {
        productUuid,
        productTitle,
        productImage,
        productPrice,
        quantity: 1,
      }
      setOrderItems([orderItem])
    }
    else {
      const orderItem = {
        productUuid,
        productTitle,
        productImage,
        productPrice,
        quantity: 1,
      }
      const prevOrderItems = [...orderItems];
      prevOrderItems.forEach((prevOrderItem) => {
        console.log("prevOrderItem", prevOrderItem)
        if(prevOrderItem.productUuid === productUuid){
          prevOrderItem.quantity += 1;
          setOrderItems([...prevOrderItems]);
          return
        }
      })
      setOrderItems([...prevOrderItems]);
    }
  }

  const addProductToOrderHandler = (productUuid: string, totalPrice: number) => {
    setOrder((order: Order) => {
      return {
        productUuids: [...order.productUuids, productUuid],
        totalPrice: order.totalPrice + totalPrice,
      }
    })
  }

  return (
    isLoading ?
    <Layout role={currentUserAppCtx!.role}>

        <main>
          {/* Hero unit */}
          <Box
            sx={{
              bgcolor: 'background.paper',
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
                    sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}
                  >
                    <CardMedia
                      component="img"
                      sx={{
                        // 16:9
                        pt: '56.25%',
                      }}
                      image={product.image}
                      alt="random"
                    />
                    <CardContent sx={{ flexGrow: 1 }}>
                      <Typography gutterBottom variant="h5" component="h2">
                        {product.title}
                      </Typography>
                      <Typography>
                        Owner: {product.displayName}
                      </Typography>
                      <Typography>
                        Description: {product.description}
                      </Typography>
                      <Typography>
                        Price: ${product.price}
                      </Typography>
                      <Typography>
                        Quantity: {product.quantity}
                      </Typography>
                    </CardContent>
                    <CardActions>
                      <Button
                        size="small"
                        onClick={() => viewProductHandler(product.uuid)}
                      >
                        View
                      </Button>

                      {currentUserAppCtx!.uuid === product.userUuid ? "" :
                      <Button
                        size="small"
                        onClick={() => {
                          addProductToOrderItemsHandler(product.uuid, product.title, product.image, product.price)
                          addProductToOrderHandler(product.uuid, product.price)
                        }}
                      >
                        Add to your order
                      </Button>
                      }
                    </CardActions>
                  </Card>
                </Grid>
              ))}
            </Grid>
          </Container>
            <Box sx={{
              position: "fixed",
              top: "65%",
              left: "70%",
              }}>
              <Grid item xs={12} md={12}>
                <Typography
                  align="center"
                  sx={{ bgcolor: "primary.main", border: "1px dashed black", mt: 1, pb: 1 }}
                  variant="h6"
                  component="div"
                >
                  Your order
                </Typography>
                <Paper>
                  <List
                    dense={false}
                      sx={{
                        width: '100%',
                        maxWidth: 360,
                        bgcolor: 'background.paper',
                        position: 'relative',
                        overflow: 'auto',
                        maxHeight: 180,
                        '& ul': { padding: 0 },
                      }}>
                        {orderItems.map(orderItem => (
                          <ListItem
                            key={orderItem.productUuid}
                            secondaryAction={
                              <IconButton edge="end" aria-label="delete">
                                <DeleteIcon />
                              </IconButton>
                            }
                          >
                            <ListItemAvatar>
                              <Avatar src={orderItem.productImage}/>
                            </ListItemAvatar>
                            <ListItemText
                              primary={`${orderItem.productTitle} x${orderItem.quantity}`}
                              secondary={null}
                            />
                          </ListItem>
                        ))}
                  </List>
                </Paper>
              </Grid>
            </Box>
        </main>
    </Layout>
    :
    <div>Loading...</div>
  );
};

export default Home;
