import * as React from "react";
import { ListItem, ListItemText, ListItemAvatar, Avatar } from "@material-ui/core";
import withStyles, {StyleRules, WithStyles} from "@material-ui/styles/withStyles/withStyles";
import {createStyles, Theme} from "@material-ui/core/styles";

const styles = (theme: Theme): StyleRules => createStyles({
});

interface IProps extends WithStyles<typeof styles> {
  text: string;
  name: string;
}

const Message: React.FC<IProps> = ({ classes, text, name }: IProps) => {
  return (
    <ListItem alignItems="flex-start">
      <ListItemAvatar>
        <Avatar />
      </ListItemAvatar>
      <ListItemText
        primary={name}
        secondary={text}
      />
    </ListItem>
  );
};

export default withStyles(styles)(Message);
