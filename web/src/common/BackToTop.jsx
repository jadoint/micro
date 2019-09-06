import React from "react";
import { useSelector } from "react-redux";
import Link from "@material-ui/core/Link";

import useStyles from "../hooks/styles";

const BackToTop = () => {
  const classes = useStyles();

  const pageNum = useSelector(state => state.listingsReducer.pageNum);
  let pages = pageNum ? `Page ${pageNum}` : "";
  return (
    <span className={classes.fixedBottomRight}>
      <span className={classes.rightSpacer}>{pages}</span>
      <Link href="#top" color="inherit" variant="body2">
        ↑ Back to Top
      </Link>
    </span>
  );
};

export default BackToTop;
