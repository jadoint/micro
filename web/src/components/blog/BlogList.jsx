import React, { useEffect, Fragment } from "react";
import { Link as RouterLink } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import Typography from "@material-ui/core/Typography";
import Grid from "@material-ui/core/Grid";
import Card from "@material-ui/core/Card";
import CardActionArea from "@material-ui/core/CardActionArea";
import CardContent from "@material-ui/core/CardContent";
import Container from "@material-ui/core/Container";
import LinearProgress from "@material-ui/core/LinearProgress";
import Button from "@material-ui/core/Button";

import config from "../../config";
import useStyles from "../../hooks/styles";
import Sidebar from "../../common/Sidebar";
import Pagination from "../../common/Pagination";
import {
  fetchListings,
  fetchAuthors,
  resetListingsState
} from "../../actions/listingsAction";

const BlogList = props => {
  const classes = useStyles();

  const listings = useSelector(state => state.listingsReducer.listings);
  const tagFilter = useSelector(state => state.listingsReducer.tagFilter);
  const pageNum = useSelector(state => state.listingsReducer.pageNum);
  const isLoading = useSelector(state => state.listingsReducer.isLoading);
  const dispatch = useDispatch();

  useEffect(() => {
    let endpoint = `${config.blogApiUrl}/latest?pageNum=${pageNum}`;
    if (tagFilter) {
      endpoint += `&tag=${tagFilter}`;
    }
    dispatch(fetchListings(endpoint));
  }, [dispatch, props, pageNum, tagFilter]);

  // Fetch authors from author IDs results from fetchListings
  useEffect(() => {
    if (listings && listings.length > 0) {
      let authorIds = listings.map(l => {
        return l.idAuthor;
      });
      // Remove duplicate author IDs
      authorIds = authorIds.filter((v, i, self) => {
        return self.indexOf(v) === i;
      });
      if (authorIds.length > 0) {
        dispatch(fetchAuthors(authorIds));
      }
    }
  }, [dispatch, listings]);

  useEffect(() => {
    return () => {
      dispatch(resetListingsState());
    };
  }, [dispatch]);

  return (
    <React.Fragment>
      <Container maxWidth="lg">
        <main>
          <Grid container spacing={4} className={classes.mainGrid}>
            <Grid item xs={12} md={8}>
              {isLoading ? (
                <div className={classes.root}>
                  <LinearProgress />
                </div>
              ) : (
                <Fragment>
                  <NewBlogButton />
                  <Filter tagFilter={tagFilter} classes={classes} />
                  <Pagination />
                  <Listings listings={listings} classes={classes} />
                  <Pagination />
                </Fragment>
              )}
            </Grid>
            <Sidebar />
          </Grid>
        </main>
      </Container>
    </React.Fragment>
  );
};

export default BlogList;

const NewBlogButton = () => {
  const username = useSelector(state => state.authReducer.username);
  if (!username) return null;

  return (
    <Button
      component={RouterLink}
      to="/blog/new"
      variant="outlined"
      size="small"
      color="primary"
      style={{ marginBottom: "1rem" }}
    >
      New Blog
    </Button>
  );
};

const Tags = ({ list }) => {
  if (!list.tags) return null;

  return (
    <Fragment>
      <Typography
        variant="caption"
        color="textSecondary"
        style={{ display: "block" }}
      >
        Tags:{" "}
        {list.tags.map(t => (
          <Fragment key={t}>{t} </Fragment>
        ))}
      </Typography>
    </Fragment>
  );
};

const Filter = ({ tagFilter, classes }) => {
  if (!tagFilter) return null;

  return (
    <Typography
      component="h6"
      variant="h6"
      color="secondary"
      className={classes.botSpacer}
    >
      {tagFilter.charAt(0).toUpperCase() + tagFilter.substring(1)}
    </Typography>
  );
};

const Listings = ({ listings, classes }) => {
  const authors = useSelector(state => state.listingsReducer.authors);
  if (!listings) return null;

  return (
    <Fragment>
      {listings.map(list => (
        <CardActionArea
          key={list.idPost}
          component={RouterLink}
          to={`/blog/view/${list.idPost}`}
        >
          <Card className={classes.card}>
            <div className={classes.cardDetails}>
              <CardContent style={{ overflowWrap: "break-word" }}>
                <Typography component="h2" variant="h5">
                  {list.title}
                </Typography>
                <Typography variant="caption" color="textSecondary">
                  {list.modified}
                  {" by "}
                  {authors ? authors[list.idAuthor] : ""}
                </Typography>
                <Typography
                  variant="caption"
                  color="textSecondary"
                  style={{ display: "block" }}
                >
                  {list.wordCount}
                  {" words"}
                </Typography>
                <Typography
                  variant="body1"
                  paragraph
                  className={classes.topSpacer}
                >
                  {list.post}
                </Typography>
                <Tags list={list} />
              </CardContent>
            </div>
          </Card>
        </CardActionArea>
      ))}
    </Fragment>
  );
};
