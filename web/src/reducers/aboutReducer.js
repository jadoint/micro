const initialState = {
  title: "",
  about: ""
};

export default (state = initialState, action) => {
  if (!action.payload) return state;

  switch (action.type) {
    case "FETCH_ABOUT_AUTHOR": {
      return {
        ...state,
        ...action.payload
      };
    }
    case "RESET_ABOUT_STATE": {
      return { ...initialState };
    }
    default:
      return state;
  }
};
