import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import TextField from "@material-ui/core/TextField";
import FormGroup from "@material-ui/core/FormGroup";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import Checkbox from "@material-ui/core/Checkbox";
import Button from "@material-ui/core/Button";
import Grid from "@material-ui/core/Grid";
import Divider from "@material-ui/core/Divider";
import Container from "@material-ui/core/Container";
import CKEditor from "@ckeditor/ckeditor5-react";
import ClassicEditor from "@ckeditor/ckeditor5-build-classic";
import LinearProgress from "@material-ui/core/LinearProgress";

import useStyles from "../../hooks/styles";
import Sidebar from "../../common/Sidebar";
import useLoginRestrict from "../../hooks/useLoginRestrict";
import {
  fetchBlogWithAuth,
  updateTitle,
  updatePost,
  updateIsDraft,
  updateIsUnlisted,
  updateBlog
} from "../../actions/blogAction";

const BlogEdit = props => {
  useLoginRestrict(props);

  const classes = useStyles();

  const idPost = useSelector(state => state.blogReducer.idPost);
  const title = useSelector(state => state.blogReducer.title);
  const post = useSelector(state => state.blogReducer.post);
  const isDraft = useSelector(state => state.blogReducer.isDraft);
  const isUnlisted = useSelector(state => state.blogReducer.isUnlisted);
  const isLoading = useSelector(state => state.blogReducer.isLoading);

  const dispatch = useDispatch();

  useEffect(() => {
    const { id } = props.match.params;
    dispatch(fetchBlogWithAuth(id, props));
  }, [dispatch, props]);

  const doSubmit = e => {
    e.preventDefault();
    updateBlog({
      idPost,
      title,
      post,
      isDraft,
      isUnlisted
    });
  };

  return (
    <React.Fragment>
      <Container maxWidth="lg">
        <main>
          <Grid container spacing={5} className={classes.mainGrid}>
            {isLoading ? (
              <div className={classes.root}>
                <LinearProgress />
              </div>
            ) : (
              <Grid item xs={12} md={8}>
                {idPost > 0 && (
                  <Button
                    component={RouterLink}
                    to={`/blog/view/${idPost}`}
                    variant="outlined"
                    size="small"
                    color="primary"
                  >
                    View Post
                  </Button>
                )}
                <form
                  onSubmit={e => doSubmit(e)}
                  autoComplete="off"
                  method="POST"
                >
                  <TextField
                    id="outlined-name"
                    name="title"
                    label="Title"
                    className={classes.textField}
                    onChange={e => dispatch(updateTitle(e.currentTarget.value))}
                    margin="normal"
                    variant="outlined"
                    value={title}
                  />
                  <Divider />
                  <CKEditor
                    editor={ClassicEditor}
                    onChange={(e, editor) => {
                      const data = editor.getData();
                      dispatch(updatePost(data));
                    }}
                    data={post}
                  />
                  <FormGroup row>
                    <FormControlLabel
                      control={
                        <Checkbox
                          name="isDraft"
                          onChange={e =>
                            dispatch(updateIsDraft(e.currentTarget.checked))
                          }
                          color="primary"
                          checked={isDraft}
                          value={isDraft}
                        />
                      }
                      label="Draft"
                    />
                    <FormControlLabel
                      control={
                        <Checkbox
                          name="isUnlisted"
                          onChange={e =>
                            dispatch(updateIsUnlisted(e.currentTarget.checked))
                          }
                          color="primary"
                          checked={isUnlisted}
                          value={isUnlisted}
                        />
                      }
                      label="Unlisted"
                    />
                  </FormGroup>
                  <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    className={classes.submit}
                  >
                    Save
                  </Button>
                </form>
              </Grid>
            )}
            <Sidebar />
          </Grid>
        </main>
      </Container>
    </React.Fragment>
  );
};

export default BlogEdit;
