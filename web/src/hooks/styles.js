import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles(theme => ({
  toolbar: {
    borderBottom: `1px solid ${theme.palette.divider}`
  },
  toolbarTitle: {
    flex: 1
  },
  toolbarSecondary: {
    justifyContent: "space-between",
    overflowX: "auto",
    marginBottom: theme.spacing(2)
  },
  toolbarLink: {
    padding: theme.spacing(1),
    flexShrink: 0
  },
  mainFeaturedPost: {
    position: "relative",
    backgroundColor: theme.palette.grey[800],
    color: theme.palette.common.white,
    marginBottom: theme.spacing(4),
    backgroundImage: "url(https://source.unsplash.com/user/erondu)",
    backgroundSize: "cover",
    backgroundRepeat: "no-repeat",
    backgroundPosition: "center"
  },
  mainPostBody: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2)
  },
  overlay: {
    position: "absolute",
    top: 0,
    bottom: 0,
    right: 0,
    left: 0,
    backgroundColor: "rgba(0,0,0,.3)"
  },
  mainFeaturedPostContent: {
    position: "relative",
    padding: theme.spacing(3),
    [theme.breakpoints.up("md")]: {
      padding: theme.spacing(6),
      paddingRight: 0
    }
  },
  mainGrid: {
    marginTop: theme.spacing(3)
  },
  card: {
    display: "flex",
    marginBottom: theme.spacing(2)
  },
  cardDetails: {
    flex: 1
  },
  cardMedia: {
    width: 160
  },
  sidebarAboutBox: {
    padding: theme.spacing(2),
    backgroundColor: theme.palette.grey[200]
  },
  sidebarSection: {
    marginTop: theme.spacing(3)
  },
  footer: {
    backgroundColor: theme.palette.background.paper,
    marginTop: theme.spacing(8),
    padding: theme.spacing(6, 0)
  },
  root: {
    flexGrow: 1,
    marginTop: theme.spacing(2)
  },
  textField: {
    width: "100%"
  },
  smallTextField: {
    width: "200px"
  },
  paper: {
    marginTop: theme.spacing(8),
    display: "flex",
    flexDirection: "column",
    alignItems: "center"
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.secondary.main
  },
  form: {
    width: "100%", // Fix IE 11 issue.
    marginTop: theme.spacing(1)
  },
  submit: {
    margin: theme.spacing(3, 0, 2)
  },
  button: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2)
  },
  leftSpaceButton: {
    margin: theme.spacing(2, 0, 2, 2)
  },
  leftSpacer: {
    marginLeft: theme.spacing(1)
  },
  rightSpacer: {
    marginRight: theme.spacing(1)
  },
  topSpacer: {
    marginTop: theme.spacing(1)
  },
  botSpacer: {
    marginBottom: theme.spacing(1)
  },
  fixedBottomRight: {
    position: "fixed",
    bottom: "20px",
    right: "30px",
    zIndex: "99",
    backgroundColor: "white",
    borderRadius: "5px",
    padding: "5px 10px"
  }
}));

export default useStyles;
