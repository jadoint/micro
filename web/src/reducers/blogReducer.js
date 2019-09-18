const initialState = {
  idPost: 0,
  idAuthor: 0,
  author: "",
  title: "",
  post: "",
  wordCount: 0,
  views: 0,
  created: "",
  modified: "",
  modifiedDatetime: "",
  isDraft: false,
  isUnlisted: false,
  status: "",
  isLoading: false,
  idVisitor: 0
};

export default (state = initialState, action) => {
  if (!action.payload) return state;

  switch (action.type) {
    case "FETCH_BLOG": {
      return { ...state, ...action.payload };
    }
    case "FETCH_AUTHOR_NAME": {
      return { ...state, ...action.payload };
    }
    case "SET_ID_AUTHOR": {
      return { ...state, ...action.payload };
    }
    case "UPDATE_TITLE": {
      return { ...state, ...action.payload };
    }
    case "UPDATE_POST": {
      return { ...state, ...action.payload };
    }
    case "UPDATE_IS_DRAFT": {
      return { ...state, ...action.payload };
    }
    case "UPDATE_IS_UNLISTED": {
      return { ...state, ...action.payload };
    }
    case "INCR_VIEWS": {
      return { ...state, ...action.payload };
    }
    case "IS_LOADING": {
      return { ...state, ...action.payload };
    }
    case "RESET_BLOG_STATE": {
      return { ...initialState };
    }
    default:
      return state;
  }
};
