const initialState = {
  idVisitor: 0,
  username: "",
  recaptchaToken: ""
};

export default (state = initialState, action) => {
  if (!action.payload) return state;

  switch (action.type) {
    case "LOGIN": {
      return { ...state, ...action.payload };
    }
    case "LOGOUT": {
      return { ...initialState };
    }
    case "SIGNUP": {
      return { ...state, ...action.payload };
    }
    case "SET_USERNAME": {
      return { ...state, ...action.payload };
    }
    case "SET_RECAPTCHA_TOKEN": {
      return { ...state, ...action.payload };
    }
    default:
      return state;
  }
};
