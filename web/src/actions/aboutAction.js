import { toast } from "react-toastify";
import config from "../config";
import http from "../services/httpService";

export const fetchAboutUser = idUser => async (dispatch, getState) => {
  try {
    const res = await http.get(`${config.userApiUrl}/about/${idUser}`);

    dispatch({
      type: "FETCH_ABOUT_USER",
      payload: { ...res.data }
    });
  } catch (error) {}
};

export const updateAboutUser = (idUser, reqPayload) => async (
  dispatch,
  getState
) => {
  try {
    await http.put(`${config.userApiUrl}/about/${idUser}`, reqPayload);
    toast("Saved", { type: toast.TYPE.SUCCESS });

    dispatch({
      type: "UPDATE_ABOUT_USER",
      payload: { ...reqPayload }
    });
  } catch (error) {}
};

export const deleteAboutUser = idUser => async (dispatch, getState) => {
  try {
    await http.delete(`${config.userApiUrl}/about/${idUser}`);
    toast("Deleted", { type: toast.TYPE.INFO });

    dispatch(fetchAboutUser(idUser));
  } catch (error) {}
};

export const setAboutCredentials = (idUser, username) => (
  dispatch,
  getState
) => {
  dispatch({ type: "SET_ABOUT_CREDENTIALS", payload: { idUser, username } });
};

export const resetAboutState = () => (dispatch, getState) => {
  dispatch({
    type: "RESET_ABOUT_STATE",
    payload: {}
  });
};
