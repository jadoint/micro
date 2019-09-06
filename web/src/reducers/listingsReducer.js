const initialState = {
  listings: [],
  authors: null,
  listingType: "",
  tagFilter: "",
  sort: "",
  prevPageNum: 0,
  pageNum: 1,
  nextPageNum: 2,
  pageCount: null,
  scrollListener: null,
  url: "",
  endpoint: "",
  isLoading: false
};

export default (state = initialState, action) => {
  if (!action.payload) return state;

  switch (action.type) {
    case "FETCH_LISTINGS": {
      return { ...state, ...action.payload };
    }
    case "FETCH_NEXT": {
      return { ...state, ...action.payload };
    }
    case "FETCH_PREV": {
      return { ...state, ...action.payload };
    }
    case "FETCH_LISTINGS_BY_TAG": {
      return { ...state, ...action.payload };
    }
    case "FETCH_AUTHORS": {
      return { ...state, ...action.payload };
    }
    case "IS_LOADING": {
      return { ...state, ...action.payload };
    }
    case "RESET_LISTINGS_STATE": {
      return {
        ...initialState,
        pageNum: state.pageNum,
        nextPageNum: state.nextPageNum,
        prevPageNum: state.prevPageNum,
        tagFilter: state.tagFilter
      };
    }
    case "FULL_RESET_LISTINGS_STATE": {
      return {
        ...initialState
      };
    }
    default:
      return state;
  }
};
