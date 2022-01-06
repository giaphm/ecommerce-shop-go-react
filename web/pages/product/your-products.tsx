import * as React from "react";
import type { NextPage } from "next";
import Router from "next/router";

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


import * as ProductsAPI from '../../src/repositories/products';
import * as UsersAPI from '../../src/repositories/users';
import { Auth, setApiClientsAuth } from '../../src/repositories/auth';

import Layout from "../../components/layout";
import CurrentUserAppCtx from "../../store/current-user-context";


const YourProducts: NextPage = () => {
  const [ownProducts, setOwnProducts] = React.useState([]);
  const [isLoading, setIsLoading] = React.useState(false);

  console.log("ownProducts", ownProducts)
  
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
      // set token in header again
      Auth.getJwtToken(false).then((token: any) => {
        console.log("token", token)
        setApiClientsAuth(token)
      })
            
      console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);

      console.log("UsersAPI.usersClient", UsersAPI.usersClient);

      console.log("LoggedIn and set currentUserAppCtx again")
      // fetchCurrentUser is async => new useEffect to chase the updates
      const currentUser = Auth.currentUser();
      console.log("currentUser", currentUser);
      currentUserAppCtx!.fetchCurrentUser({
        uuid: currentUser["uuid"],
        email: currentUser["email"],
        displayName: currentUser["name"],
        role: currentUser["role"],
      });
      setIsLoading(true);
    }
    else if (!currentUserAppCtx!["uuid"]) {
      Router.push("/login");
    }
    
    
  }, []);


  React.useEffect(() => {
    console.log("currentUserAppCtx", currentUserAppCtx)

    ProductsAPI.getShopkeeperProducts((ownProducts: any) => {
      console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);
      console.log("ProductsAPI.productsClient.authentications", ProductsAPI.productsClient.authentications);
      console.log("ownProducts", ownProducts);
      setOwnProducts(ownProducts);
    })
    
    UsersAPI.getUsers((users: any) => {
      console.log("UsersAPI.usersClient", UsersAPI.usersClient);
      console.log("UsersAPI.usersClient.authentications", UsersAPI.usersClient.authentications);
      console.log("users", users);
    })
  }, [currentUserAppCtx])

  const postNewProduct = () => {
    Router.push("/product/new-product");
  }

  const editProductHandler = (productUuid: string) => {
    Router.push({
      pathname: `/product/edit/${productUuid}`,
    });
  }

  const deleteProductHandler = (productUuid: string) => {
    ProductsAPI.deleteProduct(productUuid, () => {
      ProductsAPI.getShopkeeperProducts((ownProducts: any) => {
        console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);
        console.log("ProductsAPI.productsClient.authentications", ProductsAPI.productsClient.authentications);
        console.log("ownProducts", ownProducts);
        setOwnProducts(ownProducts);
      })
    });

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
                Your products
              </Typography>
              {/* <Typography variant="h5" align="center" color="text.secondary" paragraph>
                Something short and leading about the collection below—its contents,
                the creator, etc. Make it short and sweet, but not too short so folks
                don&apos;t simply skip over it entirely.
              </Typography> */}
              <Stack
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
                {/* <Button variant="outlined">Secondary action</Button> */}
              </Stack>
            </Container>
          </Box>
          <Container sx={{ py: 8 }} maxWidth="md">
            {/* End hero unit */}
            <Grid container spacing={4}>
              {ownProducts.map((product: any) => (
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
                        onClick={() => editProductHandler(product.uuid)}
                      >
                        Edit
                      </Button>
                      
                      <Button
                        size="small"
                        onClick={() => deleteProductHandler(product.uuid)}
                      >
                        Delete
                      </Button>
                    </CardActions>
                  </Card>
                </Grid>
              ))}
            </Grid>
          </Container>
        </main>
    </Layout>
    :
    <div>Loading...</div>
  );
};

export default YourProducts;
