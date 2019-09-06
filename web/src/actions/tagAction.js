import http from "../services/httpService";
import config from "../config";
import { toast } from "react-toastify";

export const fetchTags = (idPost, modifiedDatetime) => async (
  dispatch,
  getState
) => {
  try {
    const res = await http.get(
      `${config.blogApiUrl}/tag/${idPost}/blog_tag_${modifiedDatetime}.json`
    );

    dispatch({
      type: "FETCH_TAGS",
      payload: { ...res.data }
    });
  } catch (error) {}
};

export const fetchFrequentTags = () => async (dispatch, getState) => {
  try {
    const res = await http.get(`${config.blogApiUrl}/tag/frequent`);

    dispatch({
      type: "FETCH_FREQUENT_TAGS",
      payload: { ...res.data }
    });
  } catch (error) {}
};

export const updateTag = tag => (dispatch, getState) => {
  const re = RegExp(/^[a-z0-9-]{0,25}$/, "ig");
  if (!re.test(tag)) {
    toast(
      "Tag must contain only dashes or alphanumeric characters and be between 3 and 25 characters in length",
      {
        type: toast.TYPE.INFO
      }
    );
    return;
  }

  dispatch({
    type: "UPDATE_TAG",
    payload: { tag: tag.toLowerCase() }
  });
};

export const addTag = (idPost, tag) => async (dispatch, getState) => {
  try {
    const tags = getState().tagReducer.tags;
    if (tags.length >= 20) {
      toast("Limit of 20 tags reached", { type: toast.TYPE.INFO });
      return;
    }

    await http.post(`${config.blogApiUrl}/tag/${idPost}`, { tag });
    toast("Saved", { type: toast.TYPE.SUCCESS });

    const newTags = [...getState().tagReducer.tags, tag];

    dispatch({
      type: "ADD_TAG",
      payload: { tags: newTags }
    });
  } catch (error) {}
};

export const deleteTag = (idPost, tag) => async (dispatch, getState) => {
  try {
    await http.delete(`${config.blogApiUrl}/tag/${idPost}/${tag}`);
    toast("Deleted", { type: toast.TYPE.INFO });

    const newTags = getState().tagReducer.tags.filter(t => t !== tag);

    dispatch({
      type: "DELETE_TAG",
      payload: { tags: newTags }
    });
  } catch (error) {}
};

export const resetBlogTagsState = () => (dispatch, getState) => {
  dispatch({
    type: "RESET_BLOG_TAGS_STATE",
    payload: {}
  });
};
