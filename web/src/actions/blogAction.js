import http from "../services/httpService";
import config from "../config";
import { toast } from "react-toastify";

export const fetchBlog = idPost => async (dispatch, getState) => {
  try {
    dispatch({
      type: "IS_LOADING",
      payload: { isLoading: true }
    });

    const res = await http.get(`${config.blogApiUrl}/${idPost}`);

    const { isDraft, isUnlisted } = res.data;
    let status = "";
    if (isDraft) status += "Draft ";
    if (isUnlisted) status += "Unlisted ";

    dispatch({
      type: "FETCH_BLOG",
      payload: { ...res.data, status, isLoading: false }
    });
  } catch (error) {
    dispatch({
      type: "FETCH_BLOG",
      payload: { isLoading: false }
    });
  }
};

export const fetchBlogWithAuth = (endpoint, props) => async (
  dispatch,
  getState
) => {
  try {
    dispatch({
      type: "IS_LOADING",
      payload: { isLoading: true }
    });

    const res = await http.get(endpoint);

    if (res.data.idAuthor !== res.data.idVisitor) {
      toast("You are unauthorized for this action.", {
        type: toast.TYPE.ERROR
      });
      dispatch({
        type: "FETCH_BLOG",
        payload: { isLoading: false }
      });
      props.history.replace("/");
    } else {
      dispatch({
        type: "FETCH_BLOG",
        payload: { ...res.data, isLoading: false }
      });
    }
  } catch (error) {
    dispatch({
      type: "FETCH_BLOG",
      payload: { isLoading: false }
    });
    props.history.replace("/");
  }
};

export const fetchAuthorName = idAuthor => async (dispatch, getState) => {
  try {
    const res = await http.get(`${config.userApiUrl}/name/${idAuthor}`);

    const author = res.data.username;

    dispatch({
      type: "FETCH_AUTHOR_NAME",
      payload: { author }
    });
  } catch (error) {}
};

export const updateTitle = title => (dispatch, getState) => {
  dispatch({
    type: "UPDATE_TITLE",
    payload: { title }
  });
};

export const updatePost = post => (dispatch, getState) => {
  dispatch({
    type: "UPDATE_POST",
    payload: { post }
  });
};

export const updateIsDraft = isDraft => (dispatch, getState) => {
  // Drafts are automatically unlisted
  let payload = isDraft ? { isDraft, isUnlisted: true } : { isDraft };

  dispatch({
    type: "UPDATE_IS_DRAFT",
    payload
  });
};

export const updateIsUnlisted = isUnlisted => (dispatch, getState) => {
  dispatch({
    type: "UPDATE_IS_UNLISTED",
    payload: { isUnlisted }
  });
};

export const postBlog = async reqPayload => {
  try {
    const res = await http.post(config.blogApiUrl, reqPayload);
    toast("Saved", { type: toast.TYPE.SUCCESS });
    return res.data.idPost;
  } catch (error) {}
};

export const updateBlog = async reqPayload => {
  try {
    await http.put(`${config.blogApiUrl}/${reqPayload.idPost}`, reqPayload);
    toast("Saved", { type: toast.TYPE.SUCCESS });
  } catch (error) {}
};

export const resetBlogState = () => (dispatch, getState) => {
  dispatch({
    type: "RESET_BLOG_STATE",
    payload: {}
  });
};
