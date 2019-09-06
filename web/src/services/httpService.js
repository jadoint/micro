import axios from "axios";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.min.css";

axios.defaults.withCredentials = true;
axios.defaults.timeout = 5000;

axios.interceptors.response.use(null, error => {
  const expectedError =
    error.response &&
    error.response.status >= 400 &&
    error.response.status < 500;

  let { error: errMsg } = error.response.data;

  if (!expectedError) {
    toast("An unexpected error occurred.", { type: toast.TYPE.ERROR });
  } else if (errMsg) {
    toast(errMsg, { type: toast.TYPE.ERROR });
  } else if (error.response.status === 400) {
    toast("Invalid request to server.", { type: toast.TYPE.ERROR });
  } else if (error.response.status === 401) {
    toast(
      "You are unauthorized for this action. Please check if you are logged in.",
      { type: toast.TYPE.ERROR }
    );
  } else if (error.response.status === 403) {
    toast("You are unauthorized for this action.", { type: toast.TYPE.ERROR });
  } else if (error.response.status === 405) {
    toast("You are unauthorized for this action. (405 Method Not Allowed)", {
      type: toast.TYPE.ERROR
    });
  } else if (error.response.status === 429) {
    toast(
      "You are making too many requests too fast. (429 Too Many Requests)",
      {
        type: toast.TYPE.ERROR
      }
    );
  } else if (error.response.status === 500) {
    toast("Something went wrong on the server (500 Internal Server Error)", {
      type: toast.TYPE.ERROR
    });
  } else if (error.response.status > 405) {
    toast(`Unexpected server error [Code: ${error.response.status}]`, {
      type: toast.TYPE.ERROR
    });
  }

  return Promise.reject(error);
});

export default {
  get: axios.get,
  post: axios.post,
  put: axios.put,
  delete: axios.delete
};
