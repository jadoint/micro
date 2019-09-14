import React, { Fragment, useEffect } from "react";
import { Link as RouterLink } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import Link from "@material-ui/core/Link";

import config from "../config";
import useStyles from "../hooks/styles";
import { setUsername, logout } from "../actions/authAction";
import {
  fullResetListingsState,
  fetchListingsByTag
} from "../actions/listingsAction";

const Header = () => {
  const username = useSelector(state => state.authReducer.username);
  const dispatch = useDispatch();

  const classes = useStyles();

  const headerNav = config.header ? config.header : [];

  useEffect(() => {
    if (!username) {
      // Check local storage to repopulate visitor
      // state if the browser is refreshed. Note that
      // actual authorization with the server is done
      // through the site cookie, not through values
      // from local storage.
      if ("localStorage" in window) {
        const storedUsername = localStorage.getItem("username");
        if (storedUsername) dispatch(setUsername(storedUsername));
      }
    }
  }, [dispatch, username]);

  return (
    <React.Fragment>
      <Toolbar id="top" className={classes.toolbar}>
        <Link
          component={RouterLink}
          to={"/"}
          color="inherit"
          variant="h5"
          align="center"
          noWrap
          className={classes.toolbarTitle}
          onClick={() => dispatch(fullResetListingsState())}
        >
          {config.headerTitle}
        </Link>
        {username ? (
          <Fragment>
            <Typography
              color="inherit"
              noWrap
              variant="body2"
              style={{ marginRight: ".5rem" }}
            >
              <Link
                component={RouterLink}
                to={`/profile/${username}`}
                color="inherit"
                variant="body2"
                noWrap
              >
                {username}
              </Link>
            </Typography>
            <Button
              onClick={() => dispatch(logout())}
              variant="outlined"
              size="small"
            >
              Logout
            </Button>
          </Fragment>
        ) : (
          <Fragment>
            <Button
              component={RouterLink}
              to="/auth/login"
              variant="outlined"
              size="small"
              style={{ marginRight: ".3rem" }}
            >
              Login
            </Button>
            <Button
              component={RouterLink}
              to="/auth/signup"
              variant="outlined"
              size="small"
            >
              Sign up
            </Button>
          </Fragment>
        )}
      </Toolbar>
      <Toolbar
        component="nav"
        variant="dense"
        className={classes.toolbarSecondary}
      >
        {headerNav.map(nav => (
          <Link
            component={RouterLink}
            to={nav.link}
            color="inherit"
            noWrap
            key={nav.label}
            variant="body2"
            className={classes.toolbarLink}
            onClick={() => dispatch(fetchListingsByTag(nav.tag))}
          >
            {nav.label}
          </Link>
        ))}
      </Toolbar>
    </React.Fragment>
  );
};

export default Header;
