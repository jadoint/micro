import { combineReducers } from "redux";
import authReducer from "./authReducer";
import aboutReducer from "./aboutReducer";
import listingsReducer from "./listingsReducer";
import blogReducer from "./blogReducer";
import tagReducer from "./tagReducer";

export default combineReducers({
  authReducer,
  aboutReducer,
  listingsReducer,
  blogReducer,
  tagReducer
});
