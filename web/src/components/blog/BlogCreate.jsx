import React, { useState } from "react";
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

import useStyles from "../../hooks/styles";
import Sidebar from "../../common/Sidebar";
import useLoginRestrict from "../../hooks/useLoginRestrict";
import { handleTextChange, handleCheckboxChange } from "../../utils/input";
import { postBlog, updateBlog } from "../../actions/blogAction";

const BlogCreate = props => {
  useLoginRestrict(props);

  const classes = useStyles();

  const [idPost, setIdPost] = useState(0);
  const [inputState, setInputState] = useState({
    title: "",
    post: "",
    isDraft: false,
    isUnlisted: false
  });
  const [editorState, setEditorState] = useState("");

  const doSubmit = async e => {
    e.preventDefault();
    if (!idPost) {
      // New blog post
      const newIdPost = await postBlog({
        title: inputState.title,
        post: editorState,
        isDraft: inputState.isDraft,
        isUnlisted: inputState.isUnlisted
      });
      setIdPost(newIdPost);
    } else {
      // Updating just posted blog in /blog/new
      updateBlog({
        idPost,
        title: inputState.title,
        post: editorState,
        isDraft: inputState.isDraft,
        isUnlisted: inputState.isUnlisted
      });
    }
  };

  return (
    <React.Fragment>
      <Container maxWidth="lg">
        <main>
          <Grid container spacing={5} className={classes.mainGrid}>
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
                  onChange={event => {
                    handleTextChange({ event, inputState, setInputState });
                  }}
                  margin="normal"
                  variant="outlined"
                  value={inputState.title}
                />
                <Divider />
                <CKEditor
                  editor={ClassicEditor}
                  onChange={(e, editor) => {
                    const data = editor.getData();
                    setEditorState(data);
                  }}
                  data={editorState}
                />
                <FormGroup row>
                  <FormControlLabel
                    control={
                      <Checkbox
                        name="isDraft"
                        onChange={event =>
                          handleCheckboxChange({
                            event,
                            inputState,
                            setInputState
                          })
                        }
                        color="primary"
                        value={inputState.isDraft}
                      />
                    }
                    label="Draft"
                  />
                  <FormControlLabel
                    control={
                      <Checkbox
                        name="isUnlisted"
                        onChange={event =>
                          handleCheckboxChange({
                            event,
                            inputState,
                            setInputState
                          })
                        }
                        color="primary"
                        value={inputState.isUnlisted}
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
                  Post
                </Button>
              </form>
            </Grid>
            <Sidebar />
          </Grid>
        </main>
      </Container>
    </React.Fragment>
  );
};

export default BlogCreate;
