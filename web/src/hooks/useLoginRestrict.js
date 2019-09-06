import { useEffect } from "react";
import { useSelector } from "react-redux";

const useLoginRestrict = props => {
  let username = useSelector(state => state.authReducer.username);

  useEffect(() => {
    let isLoggedIn = false;
    if (username) isLoggedIn = true;
    if (!username && "localStorage" in window) {
      const storedUsername = localStorage.getItem("username");
      if (storedUsername) isLoggedIn = true;
    }
    if (!isLoggedIn) {
      props.history.replace("/auth/login?error");
    }
  }, [username, props]);
};

export default useLoginRestrict;
