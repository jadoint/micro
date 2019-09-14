import React, { useState, useEffect, Fragment } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import Typography from "@material-ui/core/Typography";
import Grid from "@material-ui/core/Grid";
import Divider from "@material-ui/core/Divider";
import Container from "@material-ui/core/Container";
import LinearProgress from "@material-ui/core/LinearProgress";
import Button from "@material-ui/core/Button";
import TextField from "@material-ui/core/TextField";
import Link from "@material-ui/core/Link";
import DeleteForeverIcon from "@material-ui/icons/DeleteForever";

import createMarkup from "../../utils/createMarkup";
import useStyles from "../../hooks/styles";
import Sidebar from "../../common/Sidebar";
import Nothing from "../../common/Nothing";
import {
  fetchBlog,
  fetchAuthorName,
  resetBlogState
} from "../../actions/blogAction";
import {
  fetchTags,
  updateTag,
  addTag,
  deleteTag,
  resetBlogTagsState
} from "../../actions/tagAction";
import { setAboutCredentials } from "../../actions/aboutAction";
import { fetchListingsByTag } from "../../actions/listingsAction";

const BlogView = props => {
  const classes = useStyles();

  // Redux
  const idPost = useSelector(state => state.blogReducer.idPost);
  const idAuthor = useSelector(state => state.blogReducer.idAuthor);
  const author = useSelector(state => state.blogReducer.author);
  const idVisitor = useSelector(state => state.blogReducer.idVisitor);
  const title = useSelector(state => state.blogReducer.title);
  const post = useSelector(state => state.blogReducer.post);
  const wordCount = useSelector(state => state.blogReducer.wordCount);
  const created = useSelector(state => state.blogReducer.created);
  const modified = useSelector(state => state.blogReducer.modified);
  const modifiedDatetime = useSelector(
    state => state.blogReducer.modifiedDatetime
  );
  const status = useSelector(state => state.blogReducer.status);
  const isLoading = useSelector(state => state.blogReducer.isLoading);
  const tag = useSelector(state => state.tagReducer.tag);
  const tags = useSelector(state => state.tagReducer.tags);

  const dispatch = useDispatch();

  // Local state
  const [showTagForm, setShowTagForm] = useState(false);
  const [showTagDelete, setShowTagDelete] = useState(false);

  const { id } = props.match.params;

  useEffect(() => {
    dispatch(fetchBlog(id, props));
  }, [dispatch, props, id]);

  useEffect(() => {
    if (modifiedDatetime) dispatch(fetchTags(id, modifiedDatetime));
  }, [dispatch, props, id, modifiedDatetime]);

  useEffect(() => {
    return () => {
      dispatch(resetBlogState());
      dispatch(resetBlogTagsState());
    };
  }, [dispatch]);

  useEffect(() => {
    if (idAuthor > 0) {
      dispatch(fetchAuthorName(idAuthor));
    }
  }, [dispatch, idAuthor]);

  // About
  useEffect(() => {
    if (idAuthor > 0 && author) {
      dispatch(setAboutCredentials(idAuthor, author));
    }
  }, [dispatch, idAuthor, author]);

  const doSubmitTag = e => {
    e.preventDefault();
    dispatch(addTag(idPost, tag));
  };

  const doDeleteTag = (e, deletedTag) => {
    e.preventDefault();
    dispatch(deleteTag(idPost, deletedTag));
  };

  return (
    <Fragment>
      <Container maxWidth="lg">
        <main>
          <Grid container spacing={5} className={classes.mainGrid}>
            <Grid item xs={12} md={8}>
              {isLoading ? (
                <div className={classes.root}>
                  <LinearProgress />
                </div>
              ) : (
                <Fragment>
                  {idVisitor > 0 && idVisitor === idAuthor && (
                    <Fragment>
                      <Button
                        component={RouterLink}
                        to={`/blog/edit/${idPost}`}
                        variant="outlined"
                        size="small"
                        color="default"
                        className={classes.button}
                      >
                        Edit Post
                      </Button>

                      <Button
                        variant="outlined"
                        size="small"
                        color="default"
                        className={classes.leftSpaceButton}
                        onClick={() => setShowTagForm(!showTagForm)}
                      >
                        Add Tag
                      </Button>

                      {showTagForm && (
                        <form
                          onSubmit={e => doSubmitTag(e)}
                          autoComplete="off"
                          method="POST"
                          style={{ margin: "1rem 0 1rem 0" }}
                        >
                          <TextField
                            id="tag"
                            name="tag"
                            label="Tag"
                            className={classes.smallTextField}
                            onChange={e =>
                              dispatch(updateTag(e.currentTarget.value))
                            }
                            value={tag}
                          />
                          <Button
                            type="submit"
                            variant="contained"
                            color="primary"
                            style={{ margin: ".75rem 0 0 .5rem" }}
                          >
                            Add Tag
                          </Button>
                        </form>
                      )}
                    </Fragment>
                  )}
                  {!idAuthor ? (
                    <Nothing />
                  ) : (
                    <Fragment>
                      <Typography
                        component="h1"
                        variant="h4"
                        color="inherit"
                        gutterBottom
                      >
                        {title}
                      </Typography>
                      {idAuthor > 0 && (
                        <Fragment>
                          <Divider />
                          <Typography
                            variant="body2"
                            color="textSecondary"
                            className={classes.topSpacer}
                          >
                            {created} by {author}
                            {modified !== created && (
                              <Typography
                                variant="caption"
                                color="textSecondary"
                              >
                                {" "}
                                (Updated {modified})
                              </Typography>
                            )}
                          </Typography>
                        </Fragment>
                      )}
                      {status && (
                        <Typography
                          variant="body2"
                          color="textSecondary"
                          className={classes.topSpacer}
                        >
                          Status: {status}
                        </Typography>
                      )}
                      <Typography variant="body2" color="textSecondary">
                        {wordCount} words
                      </Typography>
                      <Typography
                        className={classes.mainPostBody}
                        dangerouslySetInnerHTML={createMarkup(post)}
                      />
                      {tags.length > 0 &&
                        idVisitor > 0 &&
                        idVisitor === idAuthor && (
                          <Button
                            href="#"
                            variant="outlined"
                            size="small"
                            color="default"
                            className={classes.button}
                            onClick={e => {
                              e.preventDefault();
                              setShowTagDelete(!showTagDelete);
                            }}
                          >
                            Manage Tags
                          </Button>
                        )}
                      {tags.length > 0 && (
                        <Typography
                          variant="body2"
                          color="textSecondary"
                          className={classes.topSpacer}
                        >
                          Tags:{" "}
                          {tags.map(t => (
                            <Fragment key={t}>
                              {showTagDelete && (
                                <Link
                                  href="#"
                                  color="secondary"
                                  className={classes.leftSpacer}
                                  onClick={e => doDeleteTag(e, t)}
                                >
                                  <DeleteForeverIcon />
                                </Link>
                              )}{" "}
                              <Link
                                component={RouterLink}
                                to={`/?pageNum=1&tag=${t}`}
                                color="inherit"
                                variant="body2"
                                className={classes.rightSpacer}
                                onClick={() => dispatch(fetchListingsByTag(t))}
                              >
                                {t}
                              </Link>{" "}
                            </Fragment>
                          ))}
                        </Typography>
                      )}
                    </Fragment>
                  )}
                </Fragment>
              )}
            </Grid>
            <Sidebar />
          </Grid>
        </main>
      </Container>
    </Fragment>
  );
};

export default BlogView;
