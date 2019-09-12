const hostname = window && window.location && window.location.hostname;

let authApiUrl = "https://www.sitename.com/api/v1/auth";
let userApiUrl = "https://www.sitename.com/api/v1/user";
let blogApiUrl = "https://www.sitename.com/api/v1/blog";
if (hostname === "localhost") {
  authApiUrl = "http://localhost:8000/api/v1/auth";
  userApiUrl = "http://localhost:8000/api/v1/user";
  blogApiUrl = "http://localhost:8001/api/v1/blog";
}

const pageTitle = "David Ado";
const headerTitle = "Micro Go Blog";

const header = [
  { label: "Javascript", tag: "javascript", link: "/?tag=javascript" },
  { label: "Go", tag: "go", link: "/?tag=go" },
  { label: "Databases", tag: "databases", link: "/?tag=databases" },
  { label: "Cache", tag: "cache", link: "/?tag=cache" }
];

const footer = {
  title: "Micro Blog",
  subtitle: "Powered by Go microservices",
  copyright: "Micro"
};

// Get reCaptcha key from Google reCaptcha
const recaptchaKey = "abcdefghijklmnopqrstuvwxyz123456";

export default {
  authApiUrl,
  userApiUrl,
  blogApiUrl,
  photoUrl: "https://photo.sitename.com",
  thumbUrl: "https://photo.sitename.com/thumb_",
  pageTitle,
  headerTitle,
  header,
  recaptchaKey
};
