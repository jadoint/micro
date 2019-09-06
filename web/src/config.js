const hostname = window && window.location && window.location.hostname;

let authApiUrl = "https://www.sitename.com/api/v1/auth";
let userApiUrl = "https://www.sitename.com/api/v1/user";
let blogApiUrl = "https://www.sitename.com/api/v1/blog";
if (hostname === "localhost") {
  authApiUrl = "http://localhost:8000/api/v1/auth";
  userApiUrl = "http://localhost:8000/api/v1/user";
  blogApiUrl = "http://localhost:8001/api/v1/blog";
}

const headerTitle = "Micro Go Blog";

const header = [
  { label: "Javascript", tag: "javascript", link: "/?tag=javascript" },
  { label: "Go", tag: "go", link: "/?tag=go" },
  { label: "Databases", tag: "databases", link: "/?tag=databases" },
  { label: "Cache", tag: "cache", link: "/?tag=cache" }
];

export default {
  authApiUrl,
  userApiUrl,
  blogApiUrl,
  photoUrl: "https://photo.sitename.com",
  thumbUrl: "https://photo.sitename.com/thumb_",
  headerTitle,
  header
};
