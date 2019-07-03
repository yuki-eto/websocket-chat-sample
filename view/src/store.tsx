import { combineReducers, configureStore, getDefaultMiddleware } from "redux-starter-kit";
import userModule from "./modules/userModule";

const rootReducer = combineReducers({
  user: userModule.reducer,
});

export const setupStore = () => {
  const middleware = getDefaultMiddleware();
  return configureStore({
    reducer: rootReducer,
    middleware,
  });
};
