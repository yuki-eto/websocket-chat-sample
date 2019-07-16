import * as React from "react";

import UserName from "./user_name";
import Chat from "./chat";

const Form: React.FC<{}> = () => {
  return (
    <>
      <UserName />
      <Chat />
    </>
  );
};

export default Form;
