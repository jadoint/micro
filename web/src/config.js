const hostname = window && window.location && window.location.hostname;

let authApiUrl = "https://www.davidado.com/api/v1/auth";
let userApiUrl = "https://www.davidado.com/api/v1/user";
let blogApiUrl = "https://www.davidado.com/api/v1/blog";
if (hostname === "local.davidado.com") {
  authApiUrl = "http://local.davidado.com/api/v1/auth";
  userApiUrl = "http://local.davidado.com/api/v1/user";
  blogApiUrl = "http://local.davidado.com/api/v1/blog";
}

const pageTitle = "David Ado";
const headerTitle = "David Ado";

const header = [
  { label: "Javascript", tag: "javascript", link: "/?tag=javascript" },
  { label: "Go", tag: "go", link: "/?tag=go" },
  { label: "Databases", tag: "databases", link: "/?tag=databases" },
  { label: "Cache", tag: "cache", link: "/?tag=cache" }
];

const footer = {
  title: "Shifting bits since dial-up",
  subtitle: "Powered by micro",
  copyright: "David Ado"
};

// Get reCaptcha key from Google reCaptcha
const recaptchaKey = "6LdSJrcUAAAAAFvw7xIl1pXg4cX2zu4WeNDkv7k0";

export default {
  authApiUrl,
  userApiUrl,
  blogApiUrl,
  photoUrl: "https://photo.sitename.com",
  thumbUrl: "https://photo.sitename.com/thumb_",
  pageTitle,
  headerTitle,
  header,
  footer,
  recaptchaKey
};
