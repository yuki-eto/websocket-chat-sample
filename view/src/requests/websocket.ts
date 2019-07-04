import { userActions } from "../modules/userModule";
import { roomActions, StreamType } from "../modules/roomModule";

let ws: WebSocket;

export const openConnection = () => {
  return async (dispatch, getState) => {
    if (!!ws) {
      ws.close();
    }

    const state = getState();
    const { loginToken, accessToken } = state.user.user;

    const host = "localhost";
    const port = 19999;
    const tokens = {login_token: loginToken, access_token: accessToken};
    const tokenForms = Object.keys(tokens).map((key) => {
      return `${key}=${encodeURIComponent(tokens[key])}`;
    });
    const url = `ws://${host}:${port}/websocket?${tokenForms.join("&")}`;
    ws = new WebSocket(url);

    ws.addEventListener("open", () => {
      dispatch(userActions.socketConnected());
    });
    ws.addEventListener("close", () => {
      dispatch(userActions.socketDisconnected());

      ws = null;
    });
    ws.addEventListener("error", (e) => {
      console.error(e);
      ws.close();
      ws = null;
      dispatch(userActions.socketDisconnected());
    });
    ws.addEventListener("message", (event) => {
      if (event.data === "ping") {
        ws.send("pong");
        return;
      }

      try {
        const data = JSON.parse(event.data);
        dispatch(roomActions.addStreamLog(data));

        if (data.type === StreamType.Join) {
          dispatch(roomActions.joinUser(data.user));
        } else if (data.type === StreamType.Leave) {
          dispatch(roomActions.leaveUser(data.user));
        }
      } catch (err) {
        console.error(err);
      }
    });
  };
};
