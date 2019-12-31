const hostname = window && window.location && window.location.hostname;

let authApiUrl = "https://www.sitename.com/api/v1/auth";
let userApiUrl = "https://www.sitename.com/api/v1/user";
let blogApiUrl = "https://www.sitename.com/api/v1/blog";
if (hostname === "local.sitename.com") {
  authApiUrl = "https://local.sitename.com/api/v1/auth";
  userApiUrl = "https://local.sitename.com/api/v1/user";
  blogApiUrl = "https://local.sitename.com/api/v1/blog";
}

const pageTitle = "Micro Go Blog";
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
