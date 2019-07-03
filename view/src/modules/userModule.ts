import { createSlice, PayloadAction } from "redux-starter-kit";
import { create, login } from "../requests/user";

interface IUser {
  name: string;
  loginToken: string;
  accessToken: string;
}

export interface IUserState {
  user: IUser | null;
  isLoading: boolean;
  isLoaded: boolean;
  isError: boolean;
  isLogin: boolean;
}

const initialState: IUserState = {
  isError: false,
  isLoaded: false,
  isLoading: false,
  isLogin: false,
  user: {
    name: "",
    loginToken: "",
    accessToken: "",
  },
};

const userModule = createSlice({
  slice: "user",
  initialState,
  reducers: {
    initializeState: (state: IUserState) => {
      state.isLoading = false;
      state.isLoaded = false;
      state.isError = false;
    },
    loadedState: (state: IUserState, action: PayloadAction<{ isLoaded: boolean, isError: boolean }>) => {
      state.isLoading = false;
      state.isLoaded = action.payload.isLoaded;
      state.isError = action.payload.isError;
    },
    setLoginToken: (state: IUserState, action: PayloadAction<string>) => {
      state.user.loginToken = action.payload;
    },
    setAccessToken: (state: IUserState, action: PayloadAction<string>) => {
      state.user.accessToken = action.payload;
      state.isLogin = true;
    },
    setName: (state: IUserState, action: PayloadAction<string>) => {
      state.user.name = action.payload;
    },
  },
});

export const { actions: userActions } = userModule;

export const createUser = () => {
  return async (dispatch, getState) => {
    const { user: userState } = getState();
    if (userState.isLoading) {
      return;
    }

    try {
      dispatch(userActions.initializeState());
      const createResponse = await create(userState.user.name);
      if (createResponse.status !== 200) {
        dispatch(userActions.loadedState({ isLoaded: true, isError: true }));
        return;
      }

      dispatch(userActions.loadedState({ isLoaded: true, isError: false }));
      dispatch(userActions.setLoginToken(createResponse.data.login_token));

      dispatch(userActions.initializeState());
      const loginResponse = await login(createResponse.data.login_token);
      if (loginResponse.status !== 200) {
        dispatch(userActions.loadedState({ isLoaded: true, isError: true }));
        return;
      }

      dispatch(userActions.loadedState({ isLoaded: true, isError: false }));
      dispatch(userActions.setAccessToken(loginResponse.data.access_token));
    } catch (err) {
      console.error(err);
      dispatch(userActions.loadedState({ isLoaded: true, isError: true }));
    }
  };
};

export const loginUser = () => {
  return async (dispatch, getState) => {
    const { user: userState } = getState();
    if (userState.isLoading) {
      return;
    }
    if (!userState.user.loginToken) {
      return;
    }

    try {
      dispatch(userActions.initializeState());
      const response = await login(userState.user.loginToken);
      if (response.status !== 200) {
        dispatch(userActions.loadedState({ isLoaded: true, isError: true }));
        return;
      }

      dispatch(userActions.loadedState({ isLoaded: true, isError: false }));
      dispatch(userActions.setAccessToken(response.data.access_token));
    } catch (err) {
      console.error(err);
      dispatch(userActions.loadedState({ isLoaded: true, isError: true }));
    }
  };
};

export default userModule;
