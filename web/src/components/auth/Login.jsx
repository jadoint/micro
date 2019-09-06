import React, { useState, useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import Container from "@material-ui/core/Container";
import Avatar from "@material-ui/core/Avatar";
import LockOutlinedIcon from "@material-ui/icons/LockOutlined";
import Typography from "@material-ui/core/Typography";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import Link from "@material-ui/core/Link";
import Grid from "@material-ui/core/Grid";

import useStyles from "../../hooks/styles";
import { handleTextChange } from "../../utils/input";
import { login, setUsername } from "../../actions/authAction";

const Login = props => {
  const classes = useStyles();

  // Redux
  const username = useSelector(state => state.authReducer.username);

  // Local state
  const [inputState, setInputState] = useState("");

  const dispatch = useDispatch();

  const error = props.location.search;

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

  const doSubmit = e => {
    e.preventDefault();
    dispatch(
      login({ username: inputState.username, password: inputState.password })
    );
  };

  return (
    <Container component="main" maxWidth="xs">
      <div className={classes.paper}>
        {error && (
          <Typography component="h6" variant="subtitle1" color="error">
            Please log in
          </Typography>
        )}
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign in
        </Typography>
        <form onSubmit={e => doSubmit(e)} className={classes.form}>
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
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            className={classes.submit}
          >
            Sign In
          </Button>
          <Grid container>
            <Grid item xs>
              <Link
                component={RouterLink}
                to="/auth/forgot_password"
                variant="body2"
              >
                Forgot password?
              </Link>
            </Grid>
            <Grid item>
              <Link component={RouterLink} to="/auth/signup" variant="body2">
                Don't have an account? Sign Up
              </Link>
            </Grid>
          </Grid>
        </form>
      </div>
    </Container>
  );
};

export default Login;
