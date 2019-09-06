const initialState = {
  tag: "",
  tags: [],
  frequentTags: []
};

export default (state = initialState, action) => {
  if (!action.payload) return state;

  switch (action.type) {
    case "FETCH_TAGS": {
      return { ...state, ...action.payload };
    }
    case "FETCH_FREQUENT_TAGS": {
      return { ...state, ...action.payload };
    }
    case "UPDATE_TAG": {
      return { ...state, ...action.payload };
    }
    case "ADD_TAG": {
      return { ...state, ...action.payload };
    }
    case "DELETE_TAG": {
      return { ...state, ...action.payload };
    }
    case "RESET_BLOG_TAGS_STATE": {
      return { ...state, tag: "", tags: [] };
    }
    default:
      return state;
  }
};
