import * as React from "react";
import {useDispatch, useSelector} from "react-redux";
import {TextField, Button, Theme, WithStyles} from "@material-ui/core";
import {StyleRules, createStyles} from "@material-ui/core/styles";

import {IRootState} from "../modules";
import {userActions, createUser} from "../modules/userModule";
import withStyles from "@material-ui/core/styles/withStyles";
import {FormEvent} from "react";

const styles = (theme: Theme): StyleRules => createStyles({
  form: {
    width: "100%",
  },
});
interface IProps extends WithStyles<typeof styles> {
}

const UserName: React.FC<IProps> = ({ classes }: IProps) => {
  const dispatch = useDispatch();
  const userState = useSelector((state: IRootState) => state.user);
  const textFieldProps = {
    label: "Name",
    value: userState.user.name,
    onChange: (e) => {
      const { value } = e.target;
      dispatch(userActions.setName(value));
    },
  };
  const handleLoginBtn = (e: FormEvent) => {
    if (e) {
      e.preventDefault();
    }
    dispatch(createUser());
  };

  return (
    <form className={classes.form} onSubmit={(e) => handleLoginBtn(e)}>
      <TextField
        {...textFieldProps}
        variant="outlined"
        margin="normal"
        fullWidth
        required
        disabled={userState.isLogin || userState.isLoading}
      />
      <Button
        color="primary"
        variant="contained"
        fullWidth
        onClick={() => handleLoginBtn(null)}
        disabled={userState.isLogin || userState.isLoading || !userState.user.name}
      >
        Login
      </Button>
    </form>
  );
};

export default withStyles(styles)(UserName);
