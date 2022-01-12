import React, { useEffect, SyntheticEvent, useState } from "react";
import { useRouter } from "next/router";
import { 
  TextField,
  Button,
  Box,
  Select,
  InputLabel,
  MenuItem,
  FormHelperText,
  FormControl,
  Container,
  Typography,
  Grid,
  Card,
  CardMedia,
  CardContent,
} from "@mui/material";
import { SelectChangeEvent } from '@mui/material/Select';

import { isEmpty } from "lodash";

import CurrentUserAppCtx from "../../../store/current-user-context";

import Layout from "../../../components/layout";
import { Auth, setApiClientsAuth } from '../../../src/repositories/auth';
import * as ProductsAPI from "../../../src/repositories/products";
import * as UsersAPI from "../../../src/repositories/users";

const availableCategories = ["", "tshirt"];

const categoryName = (category: string) => category.charAt(0) + category.slice(1);

export default function ViewProduct(props: any) {
  const [users, setUsers] = React.useState([]);
  const [displayName, setDisplayName] = useState("");
  const [uuid, setUuid] = useState("");
  const [userUuid, setUserUuid] = useState("");
  const [category, setCategory] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [image, setImage] = useState("");
  const [price, setPrice] = useState("");
  const [quantity, setQuantity] = useState("");
  const [redirect, setRedirect] = useState(false);
  const [isLoading, setIsLoading] = React.useState(false);

  console.log("users", users);

  const router = useRouter();
  const { query } = router;

  console.log("query", query);

  const currentUserAppCtx = React.useContext(CurrentUserAppCtx);

  React.useEffect(() => {
    console.log(typeof window)
    console.log("currentUserAppCtx", currentUserAppCtx)
    // const mockUserLoggedIn = JSON.parse(localStorage.getItem("_mock_user") || "{}")
    const isCurrentUserLoggedIn = Auth.isLoggedIn()
    console.log("isCurrentUserLoggedIn", isCurrentUserLoggedIn)
    if (isCurrentUserLoggedIn) {
      const currentUser = Auth.currentUser();
      console.log("currentUser", currentUser);
      // fetchCurrentUser is async => new useEffect to chase the updates
      currentUserAppCtx!.fetchCurrentUser({
        uuid: currentUser["uuid"],
        email: currentUser["email"],
        displayName: currentUser["name"],
        role: currentUser["role"],
        balance: currentUser["balance"],
      });
      // set token in header again
      Auth.waitForAuthReady()
        .then(() => {
          return Auth.getJwtToken(false)
        })
        .then((token: string) => setApiClientsAuth(token))
        .then(() => {
          console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);
    
          console.log("UsersAPI.usersClient", UsersAPI.usersClient);
    
          console.log("LoggedIn and set currentUserAppCtx again")
          setIsLoading(true);
        })
    }
    else if (!currentUserAppCtx!["uuid"]) {
      router.push("/login");
    }
    

  }, [])


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
    
    
    ProductsAPI.getProduct(query.product_uuid, (product: any) => {
      console.log("ProductsAPI.productsClient", ProductsAPI.productsClient);
      console.log("ProductsAPI.productsClient.authentications", ProductsAPI.productsClient.authentications);
      console.log("product", product);
      setUuid(product.uuid);
      users.forEach((user: any) => {
        // console.log("user", user);
        if(user.uuid === product.userUuid) {
          console.log("user", user);
          setDisplayName(user.displayName)
        }
      })
      setUserUuid(product.userUuid);
      setCategory(product.category);
      setTitle(product.title);
      setDescription(product.description);
      setImage(product.image);
      setPrice(product.price);
      setQuantity(product.quantity);
    });
  }, [users])

  const submit = (e: SyntheticEvent) => {
    e.preventDefault();

    const data = {
      category,
      title,
      description,
      image,
      price: parseFloat(price),
      quantity: parseInt(quantity),
    };
    
    ProductsAPI.updateProduct(uuid, userUuid, category, title, description, image, parseFloat(price), parseInt(quantity))

    setRedirect(true);
  };

  React.useEffect(() => {
    if (redirect) {
      router.push("/product/your-products");
    }
  }, [redirect]);


  return (
    <Layout>
      <Box
        sx={{
          mt: 10,
          width: "100%",
          display: 'flex',
          flexDirection: 'column',
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Card
          sx={{ height: '50%' }}
        >
          <CardMedia
            component="img"
            sx={{
              // 16:9
              pt: '56.25%',
            }}
            image={image}
            alt={category}
          />
          <CardContent sx={{ flexGrow: 1 }}>
            <Typography gutterBottom variant="h5" component="h2">
              {title}
            </Typography>
            <Typography>
              Category: {categoryName(category)}
            </Typography>
            <Typography>
              Product owner: {displayName}
            </Typography>
            <Typography>
              Description: {description}
            </Typography>
            <Typography>
              Price: ${price}
            </Typography>
            <Typography>
              Quantity: {quantity}
            </Typography>
          </CardContent>
        </Card>
    {/* <Box
      sx={{
        bgcolor: 'background.paper',
        pt: 5,
        pb: 1,
      }}
    >
      <Container maxWidth="xl">
        <Typography
          component="h1"
          variant="h4"
          align="left"
          color="text.primary"
          gutterBottom
        >
          Update your product
        </Typography>
      </Container>
    </Box> */}
      {/* <form onSubmit={submit}>
        <div className="mt-5 mb-3">
          <FormControl sx={{ minWidth: 120 }} error={!availableCategories.includes(category)}>
            <InputLabel id="category">Category</InputLabel>
            <Select
              labelId="category"
              value={category}
              label="Category"
              onChange={(e) => setCategory(e.target.value)}
              renderValue={(value) => value === "tshirt" ? `${value}` : `⚠️ ${value}`}
            >
              <MenuItem value="">
                <em>None</em>
              </MenuItem>
              <MenuItem value={"tshirt"}>T-Shirt</MenuItem>
              <MenuItem value={"accessories"}>Accessories</MenuItem>
              <MenuItem value={"jeans"}>Jeans</MenuItem>
            </Select>
          <FormHelperText>{availableCategories.includes(category) ? "" : "Error"}</FormHelperText>
        </FormControl>
        </div>
        <div className="mb-3">
          <TextField
            label="Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
        </div>
        <div className="mb-3">
          <TextField
            label="Description"
            value={description}
            rows={4}
            multiline
            onChange={(e) => setDescription(e.target.value)}
          />
        </div>
        <div className="mb-3">
          <TextField
            label="Image url"
            value={image}
            onChange={(e) => setImage(e.target.value)}
          />
        </div>
        <div className="mb-3">
          <TextField
            label="Price"
            type="number"
            value={price}
            onChange={(e) => {
              setPrice(e.target.value);
              console.log(typeof e.target.value);
            }}
          />
        </div>
        <div className="mb-3">
          <TextField
            label="Quantity"
            type="number"
            value={quantity}
            onChange={(e) => {
              setQuantity(e.target.value);
              console.log(typeof e.target.value);
            }}
          />
        </div>
        <div className="mb-3">
          <Button variant="contained" color="primary" type="submit">
            Update this product
          </Button>
        </div>
      </form> */}
      </Box> 
    </Layout>
  );
}
