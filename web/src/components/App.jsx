import React, { useEffect } from "react";
import { BrowserRouter, Route } from "react-router-dom";
import { Provider } from "react-redux";
import { createStore, applyMiddleware } from "redux";
import thunk from "redux-thunk";
import { ToastContainer } from "react-toastify";
import CssBaseline from "@material-ui/core/CssBaseline";
import Container from "@material-ui/core/Container";
import { loadReCaptcha } from "react-recaptcha-v3";

import config from "../config";
import Header from "../common/Header";
import Footer from "../common/Footer";
import Login from "./auth/Login";
import Signup from "./auth/Signup";
import BlogCreate from "./blog/BlogCreate";
import BlogEdit from "./blog/BlogEdit";
import BlogList from "./blog/BlogList";
import BlogView from "./blog/BlogView";
import reducers from "../reducers";
import "./App.css";

const store = createStore(reducers, applyMiddleware(thunk));

const App = () => {
  useEffect(() => {
    document.title = config.pageTitle;
    loadReCaptcha(config.recaptchaKey);
  });

  return (
    <BrowserRouter>
      <CssBaseline />
      <Provider store={store}>
        <Container maxWidth="lg">
          <ToastContainer />
          <Header />
          <Route path="/" exact component={BlogList} />
          <Route path="/auth/login" exact component={Login} />
          <Route path="/auth/signup" exact component={Signup} />
          <Route path="/blog/new" exact component={BlogCreate} />
          <Route path="/blog/edit/:id" component={BlogEdit} />
          <Route path="/blog/view/:id" component={BlogView} />
        </Container>
        <Footer />
      </Provider>
    </BrowserRouter>
  );
};

export default App;
