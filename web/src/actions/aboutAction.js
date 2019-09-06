import { toast } from "react-toastify";
import config from "../config";
import http from "../services/httpService";

export const fetchAboutAuthor = idAuthor => async (dispatch, getState) => {
  try {
    const res = await http.get(`${config.userApiUrl}/about/${idAuthor}`);

    dispatch({
      type: "FETCH_ABOUT_AUTHOR",
      payload: { ...res.data }
    });
  } catch (error) {}
};

export const updateAboutAuthor = reqPayload => async (dispatch, getState) => {
  try {
    await http.put(`${config.userApiUrl}/about/${reqPayload.id}`, reqPayload);
    toast("Saved", { type: toast.TYPE.SUCCESS });

    dispatch(fetchAboutAuthor(reqPayload.id));
  } catch (error) {}
};

export const deleteAboutAuthor = idAuthor => async (dispatch, getState) => {
  try {
    await http.delete(`${config.userApiUrl}/about/${idAuthor}`);
    toast("Deleted", { type: toast.TYPE.INFO });

    dispatch(fetchAboutAuthor(idAuthor));
  } catch (error) {}
};

export const resetAboutState = () => (dispatch, getState) => {
  dispatch({
    type: "RESET_ABOUT_STATE",
    payload: {}
  });
};
