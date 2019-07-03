import * as React from "react";
import * as ReactDOM from "react-dom";
import { Theme } from "@material-ui/core/styles/createMuiTheme";
import { Container, Typography, Avatar } from "@material-ui/core";
import { createStyles } from "@material-ui/styles";

import { Provider } from "react-redux";
import Form from "./components/form";
import { setupStore } from "./store";
import withStyles, {WithStyles, StyleRules} from "@material-ui/core/styles/withStyles";

const store = setupStore();

const styles = (theme: Theme): StyleRules => createStyles({
  root: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
  },
});

interface IProps extends WithStyles<typeof styles> {
}

const App: React.FC<IProps> = ({ classes }: IProps) => {
  return (
    <Container component="main" maxWidth="sm">
      <div className={classes.root}>
        <Typography component="h1" variant="h2">WebSocket Chat</Typography>
        <Form />
      </div>
    </Container>
  );
};
const AppWithStyle = withStyles(styles)(App);

ReactDOM.render(
  <Provider store={store}><AppWithStyle /></Provider>,
  document.getElementById("root"),
);
