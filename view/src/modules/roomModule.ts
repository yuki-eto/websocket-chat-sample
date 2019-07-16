import { createSlice, PayloadAction } from "redux-starter-kit";
import { join, message } from "../requests/room";

export enum StreamType {
  Chat = "chat",
  Join = "join",
  Leave = "leave",
}

interface IRoom {
  id: string;
  name: string;
}

export interface IMessage {
  name: string;
  text: string;
  time: string;
}

export interface IUser {
  id: number;
  name: string;
}

export interface IStream {
  type: StreamType;
  user: IUser;
  message: IMessage;
}

export interface IRoomState {
  isLoading: boolean;
  isLoaded: boolean;
  isError: boolean;
  room: IRoom;
  users: {[key: string]: IUser};
  streams: IStream[];
  message: string;
}

const initialState = {
  isLoading: false,
  isLoaded: false,
  isError: false,
  room: {},
  users: {},
  streams: [],
  message: "",
};

const roomModule = createSlice({
  slice: "stream",
  initialState,
  reducers: {
    initializeState: (state: IRoomState) => {
      state.room = null;
      state.users = {};
      state.streams = [];
      state.isLoading = false;
      state.isLoaded = false;
      state.isError = false;
    },
    initializeLoadingState: (state: IRoomState) => {
      state.isLoading = false;
      state.isLoaded = false;
      state.isError = false;
    },
    loadedState: (state: IRoomState, action: PayloadAction<{ isLoaded: boolean, isError: boolean }>) => {
      state.isLoading = false;
      state.isLoaded = action.payload.isLoaded;
      state.isError = action.payload.isError;
    },
    setRoom: (state: IRoomState, action: PayloadAction<IRoom>) => {
      state.room = action.payload;
    },
    setUsers: (state: IRoomState, action: PayloadAction<IUser[]>) => {
      const users = {};
      Object.values(action.payload).forEach((user) => {
        users[user.id] = user;
      });
      state.users = users;
    },
    setMessages: (state: IRoomState, action: PayloadAction<IMessage[]>) => {
      const streams = action.payload.map((msg) => ({
          type: StreamType.Chat,
          message: msg,
          user: null,
      }));
      state.streams = [...state.streams, ...streams];
    },
    addStreamLog: (state: IRoomState, action: PayloadAction<IStream>) => {
      state.streams.unshift(action.payload);
    },
    joinUser: (state: IRoomState, action: PayloadAction<IUser>) => {
      state.users[action.payload.id] = action.payload;
    },
    leaveUser: (state: IRoomState, action: PayloadAction<IUser>) => {
      delete state.users[action.payload.id];
    },
    setMessage: (state: IRoomState, action: PayloadAction<string>) => {
      state.message = action.payload;
    },
  },
});

export const { actions: roomActions } = roomModule;
export default roomModule;

export const joinRoom = () => {
  return async (dispatch, getState) => {
    const state = getState();
    const { room: roomState } = state;
    if (roomState.isLoading) {
      return;
    }

    const { loginToken, accessToken } = state.user.user;
    const roomId = "test_room_id";
    dispatch(roomActions.initializeState());
    try {
      const resp = await join(roomId, loginToken, accessToken);
      const { status, data } = resp;
      if (status !== 200) {
        dispatch(roomActions.loadedState({ isLoaded: true, isError: true }));
        return;
      }

      dispatch(roomActions.setRoom(data.room));
      dispatch(roomActions.setUsers(data.users));
      dispatch(roomActions.setMessages(data.messages));
      dispatch(roomActions.loadedState({ isLoaded: true, isError: false }));
    } catch (e) {
      console.error(e);
      dispatch(roomActions.loadedState({ isLoaded: true, isError: true }));
    }
  };
};

export const sendMessage = () => {
  return async (dispatch, getState) => {
    const state = getState();
    const { room: roomState } = state;
    if (roomState.isLoading || !roomState.room || !roomState.room.id || !roomState.message) {
      return;
    }

    const { loginToken, accessToken } = state.user.user;
    dispatch(roomActions.initializeLoadingState());
    dispatch(roomActions.setMessage(""));
    try {
      const resp = await message(loginToken, accessToken, roomState.message);
      const { status } = resp;
      if (status !== 200) {
        dispatch(roomActions.loadedState({ isLoaded: true, isError: true }));
        return;
      }
      dispatch(roomActions.loadedState({ isLoaded: true, isError: false }));
    } catch (e) {
      console.error(e);
      dispatch(roomActions.loadedState({ isLoaded: true, isError: true }));
    }
  };
};
