const initialState = {
  idUser: 0,
  username: "",
  title: "",
  about: ""
};

export default (state = initialState, action) => {
  if (!action.payload) return state;

  switch (action.type) {
    case "FETCH_ABOUT_USER": {
      return {
        ...state,
        ...action.payload
      };
    }
    case "UPDATE_ABOUT_USER": {
      return {
        ...state,
        ...action.payload
      };
    }
    case "SET_ABOUT_CREDENTIALS": {
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
