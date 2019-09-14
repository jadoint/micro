import React, { useEffect, useState, Fragment } from "react";
import { useSelector, useDispatch } from "react-redux";
import Container from "@material-ui/core/Container";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import MenuList from "@material-ui/core/MenuList";
import MenuItem from "@material-ui/core/MenuItem";
import Avatar from "@material-ui/core/Avatar";
import LockOutlinedIcon from "@material-ui/icons/LockOutlined";
import Typography from "@material-ui/core/Typography";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";

import useStyles from "../../hooks/styles";
import About from "../../common/About";
import Nothing from "../../common/Nothing";
import { handleTextChange } from "../../utils/input";
import { getCredentials, saveNewPassword } from "../../actions/authAction";
import { setAboutCredentials } from "../../actions/aboutAction";

const Profile = props => {
  const classes = useStyles();

  const { username: paramUsername } = props.match.params;

  // Redux
  const idVisitor = useSelector(state => state.authReducer.idVisitor);
  const username = useSelector(state => state.authReducer.username);

  // Local
  const [section, setSection] = useState("about");

  const dispatch = useDispatch();

  // Credentials
  useEffect(() => {
    dispatch(getCredentials());
  }, [dispatch]);

  // About
  useEffect(() => {
    if (idVisitor > 0 && username) {
      dispatch(setAboutCredentials(idVisitor, username));
    }
  }, [dispatch, idVisitor, username]);

  const isOwner = username === paramUsername;

  // Using css to hide components instead of using
  // React conditional statements to avoid making
  // http requests on every menu item selection.
  const displayClasses = {
    about: section === "about" ? "block" : "none",
    password: section === "password" ? "block" : "none"
  };

  return (
    <Container maxWidth="lg">
      <main>
        <Grid container spacing={5} className={classes.mainGrid}>
          {!isOwner ? (
            <Nothing />
          ) : (
            <Fragment>
              <Grid item xs={12} md={4}>
                <Paper className={classes.paper}>
                  <MenuList>
                    <MenuItem onClick={() => setSection("about")}>
                      About
                    </MenuItem>
                    <MenuItem onClick={() => setSection("password")}>
                      Password
                    </MenuItem>
                  </MenuList>
                </Paper>
              </Grid>
              <Grid item xs={12} md={8}>
                <div style={{ display: displayClasses.about }}>
                  <About />
                </div>
                <div style={{ display: displayClasses.password }}>
                  <Password />
                </div>
              </Grid>
            </Fragment>
          )}
        </Grid>
      </main>
    </Container>
  );
};

export default Profile;

const Password = () => {
  const classes = useStyles();

  // Local state
  const [inputState, setInputState] = useState("");

  const dispatch = useDispatch();

  const doSubmit = e => {
    e.preventDefault();
    dispatch(
      saveNewPassword({
        oldPassword: inputState.oldPassword,
        newPassword: inputState.newPassword
      })
    );
  };

  return (
    <div className={classes.paper}>
      <Avatar className={classes.avatar}>
        <LockOutlinedIcon />
      </Avatar>
      <Typography component="h1" variant="h5">
        Change Password
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
          name="oldPassword"
          label="Old Password"
          type="password"
          id="old-password"
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
          name="newPassword"
          label="New Password"
          type="password"
          id="new-password"
        />
        <Button
          type="submit"
          fullWidth
          variant="contained"
          color="primary"
          className={classes.submit}
        >
          Change Password
        </Button>
      </form>
    </div>
  );
};
