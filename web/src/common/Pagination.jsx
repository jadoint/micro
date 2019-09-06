import React from "react";
import { useSelector, useDispatch } from "react-redux";
import Button from "@material-ui/core/Button";

import useStyles from "../hooks/styles";
import { fetchNext, fetchPrev } from "../actions/listingsAction";

const Pagination = () => {
  const classes = useStyles();

  const listings = useSelector(state => state.listingsReducer.listings);
  const pageNum = useSelector(state => state.listingsReducer.pageNum);
  const prevPageNum = useSelector(state => state.listingsReducer.prevPageNum);
  const nextPageNum = useSelector(state => state.listingsReducer.nextPageNum);

  const dispatch = useDispatch();

  return (
    <div style={{ marginBottom: "1rem" }}>
      {prevPageNum > 0 && (
        <Button
          variant="outlined"
          size="small"
          color="secondary"
          className={classes.rightSpacer}
          onClick={() => dispatch(fetchPrev())}
        >
          {`< Page ${prevPageNum}`}
        </Button>
      )}
      {nextPageNum > pageNum && listings && listings.length > 0 && (
        <Button
          variant="outlined"
          size="small"
          color="secondary"
          onClick={() => dispatch(fetchNext())}
        >
          {`Page ${nextPageNum} >`}
        </Button>
      )}
    </div>
  );
};

export default Pagination;
