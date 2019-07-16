import * as React from "react";
import {useDispatch, useSelector} from "react-redux";
import {TextField, Button, Theme, WithStyles} from "@material-ui/core";
import {StyleRules, createStyles} from "@material-ui/core/styles";

import {roomActions, sendMessage} from "../modules/roomModule";
import {IRootState} from "../modules";
import withStyles from "@material-ui/core/styles/withStyles";
import {FormEvent} from "react";

const styles = (theme: Theme): StyleRules => createStyles({
  form: {
    width: "100%",
  },
});
interface IProps extends WithStyles<typeof styles> {
}

const Chat: React.FC<IProps> = ({ classes }: IProps) => {
  const dispatch = useDispatch();
  const roomState = useSelector((state: IRootState) => state.room);

  if (!roomState.room || !roomState.room.id) {
    return null;
  }

  const textFieldProps = {
    label: "Message",
    value: roomState.message,
    onChange: (e) => {
      const { value } = e.target;
      dispatch(roomActions.setMessage(value));
    },
  };
  const handleSendBtn = (e: FormEvent) => {
    if (e) {
      e.preventDefault();
    }
    dispatch(sendMessage());
  };

  return (
    <form className={classes.form} onSubmit={(e) => handleSendBtn(e)}>
      <TextField
        {...textFieldProps}
        variant="outlined"
        margin="normal"
        fullWidth
        required
        disabled={roomState.isLoading}
      />
      <Button
        color="primary"
        variant="contained"
        onClick={() => handleSendBtn(null)}
        fullWidth
        disabled={!roomState.message}
      >
        Send
      </Button>
    </form>
  );
};

export default withStyles(styles)(Chat);
