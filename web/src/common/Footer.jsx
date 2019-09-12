import React, { Fragment } from "react";
import Container from "@material-ui/core/Container";
import Typography from "@material-ui/core/Typography";

import config from "../config";
import useStyles from "../hooks/styles";
import BackToTop from "./BackToTop";

const Footer = () => {
  const classes = useStyles();

  return (
    <Fragment>
      <footer className={classes.footer}>
        <Container maxWidth="lg">
          <Typography variant="h6" align="center" gutterBottom>
            {config.footer.title}
          </Typography>
          <Typography
            variant="subtitle1"
            align="center"
            color="textSecondary"
            component="p"
          >
            {config.footer.subtitle}
          </Typography>
          <Copyright />
        </Container>
      </footer>
      <BackToTop />
    </Fragment>
  );
};

export default Footer;

const Copyright = () => {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {"Copyright © "}
      {config.footer.copyright} {new Date().getFullYear()}
    </Typography>
  );
};
