import React, { Fragment, useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";

import useStyles from "../hooks/styles";
import { handleTextChange } from "../utils/input";
import {
  fetchAboutUser,
  updateAboutUser,
  deleteAboutUser,
  resetAboutState
} from "../actions/aboutAction";

const About = () => {
  const classes = useStyles();

  // Redux
  const username = useSelector(state => state.authReducer.username);
  const aboutUser = useSelector(state => state.aboutReducer.username);
  const aboutIdUser = useSelector(state => state.aboutReducer.idUser);
  const title = useSelector(state => state.aboutReducer.title);
  const about = useSelector(state => state.aboutReducer.about);

  const isOwner = username && aboutUser === username;

  const dispatch = useDispatch();

  // Local state
  const [showOwnerView, setShowOwnerView] = useState(false);

  // Fetch user's "about" details
  useEffect(() => {
    if (aboutIdUser && (!title || !about))
      dispatch(fetchAboutUser(aboutIdUser));
  }, [dispatch, aboutIdUser, title, about]);

  // Reset "about" state on unmount
  useEffect(() => {
    return () => {
      dispatch(resetAboutState());
    };
  }, [dispatch]);

  // Show owner view to owner by default
  // only if not title or about found.
  useEffect(() => {
    if (isOwner && !title && !about) {
      setShowOwnerView(true);
    } else if (title || about) {
      setShowOwnerView(false);
    }
  }, [isOwner, title, about]);

  let fullView = null;
  if (about || showOwnerView) {
    fullView = (
      <Paper elevation={0} className={classes.sidebarAboutBox}>
        {about && <AboutView />}
        {isOwner && (
          <Button
            variant="outlined"
            color="default"
            className={classes.submit}
            onClick={() => setShowOwnerView(!showOwnerView)}
          >
            Update About
          </Button>
        )}
        {showOwnerView && <OwnerView />}
      </Paper>
    );
  }

  return fullView;
};

export default About;

const AboutView = () => {
  const title = useSelector(state => state.aboutReducer.title);
  const about = useSelector(state => state.aboutReducer.about);

  return (
    <Fragment>
      <Typography variant="h6" gutterBottom>
        {title}
      </Typography>
      <Typography>{about}</Typography>
    </Fragment>
  );
};

const OwnerView = () => {
  const classes = useStyles();

  // Redux
  const username = useSelector(state => state.authReducer.username);
  const aboutUser = useSelector(state => state.aboutReducer.username);
  const aboutIdUser = useSelector(state => state.aboutReducer.idUser);
  const title = useSelector(state => state.aboutReducer.title);
  const about = useSelector(state => state.aboutReducer.about);

  const isOwner = username && aboutUser === username;

  // Local state
  const [inputState, setInputState] = useState({ title: "", about: "" });

  const dispatch = useDispatch();

  // Initializes input state to an
  // existing title and about.
  useEffect(() => {
    if (isOwner && (title || about)) {
      setInputState({ title, about });
    }
  }, [isOwner, title, about]);

  const doSubmit = e => {
    e.preventDefault();
    dispatch(
      updateAboutUser(aboutIdUser, {
        title: inputState.title,
        about: inputState.about
      })
    );
  };

  return (
    <Fragment>
      <form onSubmit={e => doSubmit(e)} autoComplete="off" method="POST">
        <TextField
          id="about-title"
          name="title"
          label="Name or pseudonym"
          className={classes.textField}
          onChange={event => {
            handleTextChange({ event, inputState, setInputState });
          }}
          margin="normal"
          value={!inputState.title ? title : inputState.title}
        />
        <TextField
          id="about-author"
          name="about"
          label="A little about yourself"
          className={classes.textField}
          onChange={event => {
            handleTextChange({ event, inputState, setInputState });
          }}
          margin="normal"
          style={{ marginBottom: "1.5rem" }}
          value={!inputState.about ? about : inputState.about}
        />
        <Button type="submit" variant="outlined" color="primary">
          Save
        </Button>
        {(title || about) && (
          <Button
            variant="outlined"
            color="secondary"
            style={{ marginLeft: ".5rem" }}
            onClick={() => dispatch(deleteAboutUser(aboutIdUser))}
          >
            Delete
          </Button>
        )}
      </form>
    </Fragment>
  );
};
