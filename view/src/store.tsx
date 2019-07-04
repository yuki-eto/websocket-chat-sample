import { combineReducers, configureStore, getDefaultMiddleware } from "redux-starter-kit";
import userModule from "./modules/userModule";
import roomModule from "./modules/roomModule";

const rootReducer = combineReducers({
  user: userModule.reducer,
  room: roomModule.reducer,
});

export const setupStore = () => {
  const middleware = getDefaultMiddleware();
  return configureStore({
    reducer: rootReducer,
    middleware,
  });
};
