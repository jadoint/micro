import React, { useState, useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import Container from "@material-ui/core/Container";
import Avatar from "@material-ui/core/Avatar";
import LockOutlinedIcon from "@material-ui/icons/LockOutlined";
import Typography from "@material-ui/core/Typography";
import Grid from "@material-ui/core/Grid";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import Link from "@material-ui/core/Link";
import { ReCaptcha } from "react-recaptcha-v3";

import config from "../../config";
import useStyles from "../../hooks/styles";
import { handleTextChange } from "../../utils/input";
import {
  signup,
  setUsername,
  setRecaptchaToken
} from "../../actions/authAction";

const Signup = props => {
  const classes = useStyles();

  // Redux
  const username = useSelector(state => state.authReducer.username);
  const recaptchaToken = useSelector(state => state.authReducer.recaptchaToken);

  // Local state
  const [inputState, setInputState] = useState("");

  const dispatch = useDispatch();

  // Checks if visitor is already logged in
  useEffect(() => {
    if (username !== "") {
      props.history.replace("/");
    } else {
      if ("localStorage" in window) {
        const storedUsername = localStorage.getItem("username");
        if (storedUsername) {
          dispatch(setUsername(storedUsername));
        }
      }
    }
  }, [username, props, dispatch]);

  // Set reCaptcha token in state
  const verifyRecaptcha = token => {
    dispatch(setRecaptchaToken(token));
  };

  const doSubmit = e => {
    e.preventDefault();
    dispatch(
      signup({
        username: inputState.username,
        email: inputState.email,
        password: inputState.password,
        recaptchaToken
      })
    );
  };

  return (
    <Container component="main" maxWidth="xs">
      <ReCaptcha
        sitekey={config.recaptchaKey}
        action="signup"
        verifyCallback={verifyRecaptcha}
      />
      <div className={classes.paper}>
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign up
        </Typography>
        <form onSubmit={event => doSubmit(event)} className={classes.form}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextField
                onChange={event =>
                  handleTextChange({
                    event,
                    inputState,
                    setInputState
                  })
                }
                variant="outlined"
                margin="normal"
                required
                fullWidth
                id="username"
                label="Username"
                name="username"
                autoComplete="username"
                autoFocus
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                onChange={event =>
                  handleTextChange({
                    event,
                    inputState,
                    setInputState
                  })
                }
                variant="outlined"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                onChange={event =>
                  handleTextChange({
                    event,
                    inputState,
                    setInputState
                  })
                }
                variant="outlined"
                margin="normal"
                required
                fullWidth
                name="password"
                label="Password"
                type="password"
                id="password"
                autoComplete="current-password"
              />
            </Grid>
          </Grid>
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            className={classes.submit}
          >
            Sign Up
          </Button>
          <Grid container justify="flex-end">
            <Grid item>
              <Typography gutterBottom>
                <Link component={RouterLink} to="/auth/login" variant="body2">
                  Already have an account? Sign in
                </Link>
              </Typography>
            </Grid>
          </Grid>
          <Grid container justify="flex-end">
            <Grid item>
              <Typography
                component="body2"
                variant="caption"
                color="textSecondary"
              >
                This site is protected by reCAPTCHA and the Google{" "}
                <Link href="https://policies.google.com/privacy">
                  Privacy Policy
                </Link>{" "}
                and{" "}
                <Link href="https://policies.google.com/terms">
                  Terms of Service
                </Link>{" "}
                apply.
              </Typography>
            </Grid>
          </Grid>
        </form>
      </div>
    </Container>
  );
};

export default Signup;
