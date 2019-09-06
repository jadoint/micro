import http from "../services/httpService";
import config from "../config";

export const login = ({ username, password }) => async (dispatch, getState) => {
  try {
    const res = await http.post(`${config.authApiUrl}/login`, {
      username,
      password
    });

    const { id: idVisitor, username: dbUsername } = res.data;

    if (idVisitor > 0 && "localStorage" in window) {
      localStorage.setItem("username", dbUsername);
    }

    dispatch({
      type: "LOGIN",
      payload: { username: dbUsername }
    });
  } catch (error) {}
};

export const logout = () => async (dispatch, getState) => {
  try {
    await http.post(`${config.authApiUrl}/logout`);

    if ("localStorage" in window) {
      localStorage.removeItem("username");
    }

    dispatch({
      type: "LOGOUT",
      payload: {}
    });
  } catch (error) {}
};

export const signup = ({ username, email, password }) => async (
  dispatch,
  getState
) => {
  try {
    const res = await http.post(`${config.authApiUrl}/signup`, {
      username,
      email,
      password
    });

    const { id: idVisitor, username: dbUsername } = res.data;

    if (idVisitor > 0 && "localStorage" in window) {
      localStorage.setItem("username", dbUsername);
    }

    dispatch({
      type: "LOGIN",
      payload: { username: dbUsername }
    });
  } catch (error) {}
};

export const setUsername = username => (dispatch, getState) => {
  dispatch({
    type: "SET_USERNAME",
    payload: { username }
  });
};
