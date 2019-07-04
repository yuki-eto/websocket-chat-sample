import { IUserState } from "./userModule";
import { IRoomState } from "./roomModule";

export interface IRootState {
  user: IUserState;
  room: IRoomState;
}
