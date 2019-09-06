import config from "../config";
import http from "../services/httpService";

export const fetchListings = endpoint => async (dispatch, getState) => {
  try {
    dispatch({
      type: "IS_LOADING",
      payload: { isLoading: true }
    });

    const res = await http.get(endpoint);

    dispatch({
      type: "FETCH_LISTINGS",
      payload: { ...res.data, isLoading: false }
    });
  } catch (error) {
    dispatch({
      type: "IS_LOADING",
      payload: { isLoading: false }
    });
  }
};

export const fetchNext = () => (dispatch, getState) => {
  try {
    const pageNum = getState().listingsReducer.pageNum;
    const prevPageNum = getState().listingsReducer.prevPageNum;
    const nextPageNum = getState().listingsReducer.nextPageNum;

    dispatch({
      type: "FETCH_NEXT",
      payload: {
        pageNum: pageNum + 1,
        prevPageNum: prevPageNum + 1,
        nextPageNum: nextPageNum + 1
      }
    });
  } catch (error) {}
};

export const fetchPrev = () => (dispatch, getState) => {
  try {
    const pageNum = getState().listingsReducer.pageNum;
    const prevPageNum = getState().listingsReducer.prevPageNum;
    const nextPageNum = getState().listingsReducer.nextPageNum;

    dispatch({
      type: "FETCH_NEXT",
      payload: {
        pageNum: pageNum - 1,
        prevPageNum: prevPageNum - 1,
        nextPageNum: nextPageNum - 1
      }
    });
  } catch (error) {}
};

export const fetchListingsByTag = tagFilter => (dispatch, getState) => {
  try {
    dispatch({
      type: "FETCH_LISTINGS_BY_TAG",
      payload: { tagFilter, pageNum: 1, prevPageNum: 0, nextPageNum: 2 }
    });
  } catch (error) {}
};

export const fetchAuthors = ids => async (dispatch, getState) => {
  try {
    const res = await http.post(`${config.userApiUrl}/names`, { ids });

    const names = res.data.usernames;

    const authors = {};
    names.forEach(item => {
      authors[item.id] = item.username;
    });

    dispatch({
      type: "FETCH_AUTHORS",
      payload: { authors }
    });
  } catch (error) {}
};

export const resetListingsState = () => (dispatch, getState) => {
  dispatch({
    type: "RESET_LISTINGS_STATE",
    payload: {}
  });
};

export const fullResetListingsState = () => (dispatch, getState) => {
  dispatch({
    type: "FULL_RESET_LISTINGS_STATE",
    payload: {}
  });
};
