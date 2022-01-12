import React, { useEffect, SyntheticEvent, useState } from "react";
import { useRouter } from "next/router";
import axios from "axios";
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
} from "@mui/material";
import { SelectChangeEvent } from '@mui/material/Select';

import { isEmpty } from "lodash";

import CurrentUserAppCtx from "../../store/current-user-context";

import Layout from "../../components/layout";
import { Auth, setApiClientsAuth } from '../../src/repositories/auth';
import * as ProductsAPI from "../../src/repositories/products";
import * as UsersAPI from "../../src/repositories/users";

const availableCategories = ["", "tshirt"];

export default function ProductForm(props: any) {
  const [category, setCategory] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [image, setImage] = useState("");
  const [price, setPrice] = useState("");
  const [quantity, setQuantity] = useState("");
  const [redirect, setRedirect] = useState(false);

  const router = useRouter();
  const { query } = router;

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
          // fetchCurrentUser is async => new useEffect to chase the updates
          currentUserAppCtx!.fetchCurrentUser({
            uuid: currentUser["uuid"],
            email: currentUser["email"],
            displayName: currentUser["name"],
            role: currentUser["role"],
            balance: currentUser["balance"],
          });
        })
    }
    else if (!currentUserAppCtx!["uuid"]) {
      router.push("/login");
    }
    
  }, [])

  React.useEffect(() => {
    console.log("currentUserAppCtx", currentUserAppCtx)
  }, [currentUserAppCtx])


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

    ProductsAPI.addProduct(category, title, image, description, parseFloat(price), parseInt(quantity))

    setRedirect(true);
  };

  React.useEffect(() => {
    if (redirect) {
      router.push("/products");
    }
  }, [redirect]);


  return (
    <Layout>
    <Box
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
          Create a new product
        </Typography>
      </Container>
    </Box>
      <form onSubmit={submit}>
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
            Post this product
          </Button>
        </div>
      </form>
    </Layout>
  );
}
