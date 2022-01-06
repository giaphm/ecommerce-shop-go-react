import React, { useEffect, SyntheticEvent, useState } from "react";
import { TextField, Button } from "@mui/material";
import { useRouter } from "next/router";
import axios from "axios";
import Layout from "../../components/layout";

export default function ProductForm(props: any) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [image, setImage] = useState("");
  const [price, setPrice] = useState("");
  const [redirect, setRedirect] = useState(false);

  const router = useRouter();
  const { query } = router;

  useEffect(() => {
    if (props.match.params.id) {
      console.log(props.match.params);
      (async () => {
        const { data } = await axios.get(`products/${props.match.params.id}`);

        console.log(data);

        setTitle(data.title);
        setDescription(data.description);
        setImage(data.image);
        setPrice(data.price);
      })();
    }
  }, []);

  const submit = async (e: SyntheticEvent) => {
    e.preventDefault();

    const data = {
      title,
      description,
      image,
      price: parseFloat(price),
    };
    if (props.match.params.id) {
      await axios.put(`products/${props.match.params.id}`, data);
    } else {
      await axios.post("products", data);
    }

    setRedirect(true);
  };

  if (redirect) {
    return router.push("/products");
  }

  return (
    <Layout>
      <form onSubmit={submit}>
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
            label="Image"
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
        <Button variant="contained" color="primary" type="submit">
          Submit
        </Button>
      </form>
    </Layout>
  );
}
