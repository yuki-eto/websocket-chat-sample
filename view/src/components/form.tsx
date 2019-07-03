import * as React from "react";
import withStyles, {StyleRules, WithStyles} from "@material-ui/styles/withStyles/withStyles";
import {createStyles, Theme} from "@material-ui/core/styles";

import UserName from "./user_name";

const styles = (theme: Theme): StyleRules => createStyles({
  form: {
    width: "100%",
  },
});

interface IProps extends WithStyles<typeof styles> {
}

const Form: React.FC<IProps> = ({ classes }: IProps) => {
  return (
    <form className={classes.form}>
      <UserName />
    </form>
  );
};

export default withStyles(styles)(Form);
