import { default as React, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { List } from "@material-ui/core";
import withStyles, {StyleRules, WithStyles} from "@material-ui/styles/withStyles/withStyles";
import {createStyles, Theme} from "@material-ui/core/styles";

import { IRootState } from "../modules";
import {joinRoom, StreamType} from "../modules/roomModule";
import {openConnection} from "../requests/websocket";

import Message from "./message";

const styles = (theme: Theme): StyleRules => createStyles({
  root: {
    width: "100%",
    maxWidth: 640,
    backgroundColor: theme.palette.background.paper,
  },
});

interface IProps extends WithStyles<typeof styles> {
}

const Streams: React.FC<IProps> = ({ classes }: IProps) => {
  const dispatch = useDispatch();
  const userState = useSelector((state: IRootState) => state.user);
  const roomState = useSelector((state: IRootState) => state.room);

  useEffect(() => {
    if (!userState.isLogin) {
      return;
    }
    if (!userState.isSocketConnected) {
      dispatch(openConnection());
      dispatch(joinRoom());
    }
  }, [userState.isLogin, userState.isSocketConnected]);

  let key = 0;
  const streams = roomState.streams.map((stream) => {
    if (!stream) {
      return null;
    }
    switch (stream.type) {
      case StreamType.Chat:
        const { message } = stream;
        return <Message key={key++} text={message.text} name={message.name} />;
      case StreamType.Join:
        const { user } = stream;
        return <Message key={key++} text="Joined" name={user.name} />;
    }
  }).filter((stream) => stream !== null);

  return (
    <List className={classes.root}>
      {streams}
    </List>
  );
};

export default withStyles(styles)(Streams);
