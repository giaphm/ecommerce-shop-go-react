// display profile
import React, { SyntheticEvent, useState, useEffect, Dispatch } from "react";
import { TextField, Button } from "@mui/material";
import axios from "axios";

import Layout from "../components/layout";

function Profile(props: any) {
  const [first_name, setFirstName] = useState("");
  const [last_name, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [password_confirm, setPasswordConfirm] = useState("");

  useEffect(() => {
    console.log(props.user);
    setFirstName(props.user.first_name);
    setLastName(props.user.last_name);
    setEmail(props.user.email);
  }, [props.user]);

  const infoSubmit = async (e: SyntheticEvent) => {
    e.preventDefault();

    const { data } = await axios.put("/users/info", {
      first_name,
      last_name,
      email,
    });

    props.setUser(data);
  };
  const passwordSubmit = async (e: SyntheticEvent) => {
    e.preventDefault();

    await axios.put("/users/password", {
      password,
      password_confirm,
    });
  };

  return (
    <Layout>
      <h3>Account Information</h3>
      <form onSubmit={infoSubmit}>
        <div className="mb-3">
          <TextField
            label="First Name"
            value={first_name}
            onChange={(e) => setFirstName(e.target.value)}
          />
        </div>
        <div className="mb-3">
          <TextField
            label="Last Name"
            value={last_name}
            onChange={(e) => setLastName(e.target.value)}
          />
        </div>
        <div className="mb-3">
          <TextField
            label="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>
        <Button color="primary" variant="contained" type="submit">
          Submit
        </Button>
      </form>

      <h3 className="mt-4">Change password</h3>
      <form onSubmit={passwordSubmit}>
        <div className="mb-3">
          <TextField
            label="Password"
            type="password"
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <div className="mb-3">
          <TextField
            label="Password Confirm"
            type="password"
            onChange={(e) => setPasswordConfirm(e.target.value)}
          />
        </div>
        <Button color="primary" variant="contained" type="submit">
          Submit
        </Button>
      </form>
    </Layout>
  );
}

export default Profile;
