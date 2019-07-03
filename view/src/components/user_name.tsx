import * as React from "react";
import {useDispatch, useSelector} from "react-redux";
import { TextField, Button } from "@material-ui/core";

import {IRootState} from "../modules";
import {userActions, createUser} from "../modules/userModule";

const UserName: React.FC<{}> = () => {
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
  const handleLoginBtn = () => dispatch(createUser());

  return (
    <>
      <TextField
        {...textFieldProps}
        variant="outlined"
        margin="normal"
        fullWidth
        disabled={userState.isLogin || userState.isLoading}
      />
      <Button
        color="primary"
        variant="contained"
        fullWidth
        onClick={handleLoginBtn}
        disabled={userState.isLogin || userState.isLoading}
      >
        Login
      </Button>
    </>
  );
};

export default UserName;
