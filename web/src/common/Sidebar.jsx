import React, { useEffect, Fragment } from "react";
import { Link as RouterLink } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import Typography from "@material-ui/core/Typography";
import Grid from "@material-ui/core/Grid";
import Link from "@material-ui/core/Link";

import useStyles from "../hooks/styles";
import About from "./About";
import { fetchFrequentTags } from "../actions/tagAction";
import { fetchListingsByTag } from "../actions/listingsAction";

const Sidebar = () => {
  return (
    <Grid item xs={12} md={4}>
      <Fragment>
        <About />
        <FrequentTags />
      </Fragment>
    </Grid>
  );
};

export default Sidebar;

const FrequentTags = () => {
  const classes = useStyles();

  const frequentTags = useSelector(state => state.tagReducer.frequentTags);
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(fetchFrequentTags());
  }, [dispatch]);

  let frequentTagsView = null;
  if (frequentTags && frequentTags.length > 0) {
    frequentTagsView = (
      <Fragment>
        <Typography
          variant="h6"
          gutterBottom
          className={classes.sidebarSection}
        >
          Frequent Tags
        </Typography>
        {frequentTags.map(tag => (
          <Link
            component={RouterLink}
            to={`/?pageNum=1&tag=${tag}`}
            display="block"
            variant="body1"
            key={tag}
            onClick={() => dispatch(fetchListingsByTag(tag))}
          >
            {tag}
          </Link>
        ))}
      </Fragment>
    );
  }

  return frequentTagsView;
};
