import React, { Fragment, useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";

import useStyles from "../hooks/styles";
import { handleTextChange } from "../utils/input";
import {
  fetchAboutAuthor,
  updateAboutAuthor,
  deleteAboutAuthor,
  resetAboutState
} from "../actions/aboutAction";

const About = () => {
  const classes = useStyles();

  // Redux
  const username = useSelector(state => state.authReducer.username);
  const author = useSelector(state => state.blogReducer.author);
  const idAuthor = useSelector(state => state.blogReducer.idAuthor);
  const title = useSelector(state => state.aboutReducer.title);
  const about = useSelector(state => state.aboutReducer.about);

  const dispatch = useDispatch();

  // Local state
  const [inputState, setInputState] = useState({ title: "", about: "" });
  const [showOwnerView, setShowOwnerView] = useState(false);

  useEffect(() => {
    if (idAuthor && (!title || !about)) dispatch(fetchAboutAuthor(idAuthor));
  }, [dispatch, idAuthor, title, about]);

  useEffect(() => {
    return () => {
      dispatch(resetAboutState());
    };
  }, [dispatch]);

  useEffect(() => {
    if (username && author === username && !title && !about) {
      setShowOwnerView(true);
    } else if (title || about) {
      setShowOwnerView(false);
    }
  }, [author, username, title, about]);

  const doSubmit = e => {
    e.preventDefault();
    dispatch(
      updateAboutAuthor({
        id: idAuthor,
        title: inputState.title,
        about: inputState.about
      })
    );
  };

  const doDelete = () => {
    dispatch(deleteAboutAuthor(idAuthor));
  };

  let aboutView = null;
  if (about) {
    aboutView = (
      <Fragment>
        <Typography variant="h6" gutterBottom>
          {title}
        </Typography>
        <Typography>{about}</Typography>
        {author === username && (
          <Button
            variant="outlined"
            color="default"
            className={classes.submit}
            onClick={() => setShowOwnerView(!showOwnerView)}
          >
            Update About
          </Button>
        )}
      </Fragment>
    );
  }
  let ownerView = null;
  if (showOwnerView) {
    ownerView = (
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
              onClick={() => doDelete()}
            >
              Delete
            </Button>
          )}
        </form>
      </Fragment>
    );
  }
  let fullView = null;
  if (aboutView || ownerView) {
    fullView = (
      <Paper elevation={0} className={classes.sidebarAboutBox}>
        {aboutView}
        {ownerView}
      </Paper>
    );
  }

  return fullView;
};

export default About;
